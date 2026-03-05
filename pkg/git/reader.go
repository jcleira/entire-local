package git

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CommitInfo struct {
	Hash    string
	Subject string
	Author  string
}

type Reader struct {
	repoRoot string
}

func NewReader(repoRoot string) *Reader {
	return &Reader{repoRoot: repoRoot}
}

func (r *Reader) ListShards() ([]string, error) {
	out, err := ExecGit(r.repoRoot, "ls-tree", "--name-only", CheckpointBranch)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

func (r *Reader) ListCheckpointRests(shard string) ([]string, error) {
	out, err := ExecGit(r.repoRoot, "ls-tree", "--name-only", CheckpointBranch+":"+shard)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

func (r *Reader) ReadFile(path string) (string, error) {
	return ExecGit(r.repoRoot, "show", CheckpointBranch+":"+path)
}

func (r *Reader) FindCommitBefore(branch string, before time.Time) (CommitInfo, error) {
	out, err := ExecGit(r.repoRoot, "log", branch, "--format=%H%x09%s%x09%an", "-1", "--before="+before.Format(time.RFC3339))
	if err != nil {
		return CommitInfo{}, err
	}
	if out == "" {
		return CommitInfo{}, fmt.Errorf("no commit found before %s", before)
	}
	parts := strings.SplitN(out, "\t", 3)
	if len(parts) < 3 {
		return CommitInfo{}, fmt.Errorf("unexpected log format: %s", out)
	}
	return CommitInfo{
		Hash:    parts[0],
		Subject: parts[1],
		Author:  parts[2],
	}, nil
}

func (r *Reader) DiffStats(hash string) (adds, dels int, err error) {
	out, err := ExecGit(r.repoRoot, "diff", "--numstat", hash+"~1", hash)
	if err != nil {
		return 0, 0, err
	}
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		a, aErr := strconv.Atoi(fields[0])
		if aErr != nil {
			continue
		}
		d, dErr := strconv.Atoi(fields[1])
		if dErr != nil {
			continue
		}
		adds += a
		dels += d
	}
	return adds, dels, nil
}

func (r *Reader) DiffContent(hash string) (string, error) {
	return ExecGit(r.repoRoot, "diff", hash+"~1", hash)
}

func (r *Reader) ListEntries(path string) ([]string, error) {
	out, err := ExecGit(r.repoRoot, "ls-tree", "--name-only", CheckpointBranch+":"+path)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}
