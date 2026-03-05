package checkpoint

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jcleira/entire-local/pkg/git"
)

type Loader struct {
	reader *git.Reader
}

func NewLoader(reader *git.Reader) *Loader {
	return &Loader{reader: reader}
}

func (l *Loader) ListCheckpoints() ([]CheckpointSummary, error) {
	shards, err := l.reader.ListShards()
	if err != nil {
		return nil, err
	}

	var summaries []CheckpointSummary

	for _, shard := range shards {
		rests, err := l.reader.ListCheckpointRests(shard)
		if err != nil {
			continue
		}

		for _, rest := range rests {
			id := shard + rest
			summary, err := l.loadSummary(id, shard, rest)
			if err != nil {
				continue
			}
			summaries = append(summaries, *summary)
		}
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].CreatedAt.After(summaries[j].CreatedAt)
	})

	return summaries, nil
}

func (l *Loader) LoadCheckpoint(id string) (*Checkpoint, error) {
	shard := Shard(id)
	rest := Rest(id)
	prefix := shard + "/" + rest

	metaJSON, err := l.reader.ReadFile(prefix + "/metadata.json")
	if err != nil {
		return nil, err
	}

	var meta Metadata
	err = json.Unmarshal([]byte(metaJSON), &meta)
	if err != nil {
		return nil, err
	}

	var sessionMeta SessionMetadata
	if sessionJSON, readErr := l.reader.ReadFile(prefix + "/0/metadata.json"); readErr == nil {
		if unmarshalErr := json.Unmarshal([]byte(sessionJSON), &sessionMeta); unmarshalErr != nil {
			sessionMeta = SessionMetadata{}
		}
	}

	var prompt, context, diff, plan string
	if v, readErr := l.reader.ReadFile(prefix + "/0/prompt.txt"); readErr == nil {
		prompt = v
	}
	if v, readErr := l.reader.ReadFile(prefix + "/0/context.md"); readErr == nil {
		context = v
	}
	if v, readErr := l.reader.ReadFile(prefix + "/0/diff.patch"); readErr == nil {
		diff = v
	}
	if v, readErr := l.reader.ReadFile(prefix + "/0/plan.md"); readErr == nil {
		plan = v
	}

	if plan == "" {
		plan = extractPlanFromPrompt(prompt)
	}

	if diff == "" && meta.Branch != "" {
		if createdAt, parseErr := time.Parse(time.RFC3339Nano, sessionMeta.CreatedAt); parseErr == nil && !createdAt.IsZero() {
			if info, findErr := l.reader.FindCommitBefore(meta.Branch, createdAt.Add(10*time.Second)); findErr == nil {
				if d, diffErr := l.reader.DiffContent(info.Hash); diffErr == nil {
					diff = d
				}
			}
		}
	}

	return &Checkpoint{
		Metadata: meta,
		Session:  sessionMeta,
		Prompt:   prompt,
		Context:  context,
		Diff:     diff,
		Plan:     plan,
	}, nil
}

func (l *Loader) LoadTranscript(id string) ([]TranscriptEntry, error) {
	shard := Shard(id)
	rest := Rest(id)

	content, err := l.reader.ReadFile(shard + "/" + rest + "/0/full.jsonl")
	if err != nil {
		return nil, err
	}

	return ParseTranscript(content)
}

func (l *Loader) loadSummary(id, shard, rest string) (*CheckpointSummary, error) {
	prefix := shard + "/" + rest

	metaJSON, err := l.reader.ReadFile(prefix + "/metadata.json")
	if err != nil {
		return nil, err
	}

	var meta Metadata
	err = json.Unmarshal([]byte(metaJSON), &meta)
	if err != nil {
		return nil, err
	}

	summary := &CheckpointSummary{
		ID:           id,
		Branch:       meta.Branch,
		FileCount:    len(meta.FilesTouched),
		SessionCount: len(meta.Sessions),
		TotalTokens:  meta.TokenUsage.OutputTokens,
	}

	var sessionMeta SessionMetadata
	if sessionJSON, readErr := l.reader.ReadFile(prefix + "/0/metadata.json"); readErr == nil {
		if json.Unmarshal([]byte(sessionJSON), &sessionMeta) == nil {
			summary.Agent = sessionMeta.Agent
			summary.AgentPercent = int(sessionMeta.Attribution.AgentPercentage)
			createdAt, parseErr := time.Parse(time.RFC3339, sessionMeta.CreatedAt)
			if parseErr != nil {
				if t, nanoErr := time.Parse(time.RFC3339Nano, sessionMeta.CreatedAt); nanoErr == nil {
					createdAt = t
				}
			}
			summary.CreatedAt = createdAt
		}
	}

	if summary.SessionCount == 0 {
		entries, listErr := l.reader.ListEntries(prefix)
		if listErr == nil {
			for _, name := range entries {
				if _, atoiErr := strconv.Atoi(name); atoiErr == nil {
					summary.SessionCount++
				}
			}
		}
	}

	prompt, err := l.reader.ReadFile(prefix + "/0/prompt.txt")
	if err == nil {
		summary.Prompt = prompt
	}

	context, err := l.reader.ReadFile(prefix + "/0/context.md")
	if err == nil {
		summary.Context = context
	}

	diff, err := l.reader.ReadFile(prefix + "/0/diff.patch")
	if err == nil {
		_, adds, dels := parseDiffStats(diff)
		summary.Additions = adds
		summary.Deletions = dels
	}

	if !summary.CreatedAt.IsZero() && meta.Branch != "" {
		info, err := l.reader.FindCommitBefore(meta.Branch, summary.CreatedAt.Add(10*time.Second))
		if err == nil {
			summary.CommitMessage = info.Subject
			summary.CommitHash = info.Hash
			summary.Author = info.Author
			adds, dels, err := l.reader.DiffStats(info.Hash)
			if err == nil && (adds > 0 || dels > 0) {
				summary.Additions = adds
				summary.Deletions = dels
			}
		}
	}

	return summary, nil
}

func extractPlanFromPrompt(prompt string) string {
	markers := []string{
		"Implement the following plan:\n",
		"Implement the following plan:\r\n",
	}
	for _, m := range markers {
		if idx := strings.Index(prompt, m); idx >= 0 {
			return strings.TrimSpace(prompt[idx+len(m):])
		}
	}
	return ""
}

func Shard(id string) string {
	if len(id) < 2 {
		return id
	}
	return id[:2]
}

func Rest(id string) string {
	if len(id) <= 2 {
		return ""
	}
	return id[2:]
}

func parseDiffStats(diff string) (files, adds, dels int) {
	filesSeen := make(map[string]bool)
	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "diff --git"):
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				filesSeen[parts[3]] = true
			}
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			adds++
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			dels++
		}
	}
	files = len(filesSeen)
	return
}
