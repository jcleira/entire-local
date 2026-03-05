package checkpoint

import (
	"bufio"
	"encoding/json"
	"strings"
	"time"
)

type flexTimestamp struct {
	Time time.Time
}

func (ft *flexTimestamp) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		ft.Time = time.Unix(int64(f), 0)
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		t, err := time.Parse(time.RFC3339Nano, s)
		if err == nil {
			ft.Time = t
			return nil
		}
		t, err = time.Parse(time.RFC3339, s)
		if err == nil {
			ft.Time = t
			return nil
		}
	}

	return nil
}

type contentBlock struct {
	Type  string          `json:"type"`
	Text  string          `json:"text,omitempty"`
	Name  string          `json:"name,omitempty"`
	ID    string          `json:"id,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`
}

type messageContent struct {
	Role    string         `json:"role"`
	Content []contentBlock `json:"content"`
}

type jsonlEntry struct {
	Type          string          `json:"type"`
	Role          string          `json:"role,omitempty"`
	Message       json.RawMessage `json:"message,omitempty"`
	Content       json.RawMessage `json:"content,omitempty"`
	Timestamp     flexTimestamp   `json:"timestamp"`
	ContentBlocks []contentBlock  `json:"contentBlocks,omitempty"`
}

func ParseTranscript(content string) ([]TranscriptEntry, error) {
	lineCount := strings.Count(content, "\n") + 1
	entries := make([]TranscriptEntry, 0, lineCount)

	scanner := bufio.NewScanner(strings.NewReader(content))
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var entry jsonlEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue
		}

		text, toolName := extractContent(entry)

		role := entry.Role
		if role == "" {
			role = entry.Type
		}

		entries = append(entries, TranscriptEntry{
			Type:      entry.Type,
			Role:      role,
			Content:   text,
			ToolName:  toolName,
			Timestamp: entry.Timestamp.Time,
		})
	}

	return entries, scanner.Err()
}

const (
	contentTypeText    = "text"
	contentTypeToolUse = "tool_use"
)

func extractContent(entry jsonlEntry) (text, toolName string) {
	if len(entry.ContentBlocks) > 0 {
		for _, b := range entry.ContentBlocks {
			switch b.Type {
			case contentTypeText:
				text += b.Text
			case contentTypeToolUse:
				toolName = b.Name
			}
		}
		return text, toolName
	}

	if entry.Message != nil {
		var mc messageContent
		if json.Unmarshal(entry.Message, &mc) == nil && len(mc.Content) > 0 {
			for _, b := range mc.Content {
				switch b.Type {
				case contentTypeText:
					text += b.Text
				case contentTypeToolUse:
					toolName = b.Name
				}
			}
			return text, toolName
		}

		var ms struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}
		if json.Unmarshal(entry.Message, &ms) == nil && ms.Content != "" {
			return ms.Content, ""
		}

		var s string
		if json.Unmarshal(entry.Message, &s) == nil {
			return s, ""
		}
	}

	if entry.Content != nil {
		var blocks []contentBlock
		if json.Unmarshal(entry.Content, &blocks) == nil && len(blocks) > 0 {
			for _, b := range blocks {
				switch b.Type {
				case contentTypeText:
					text += b.Text
				case contentTypeToolUse:
					toolName = b.Name
				}
			}
			return text, toolName
		}

		var s string
		if json.Unmarshal(entry.Content, &s) == nil {
			return s, ""
		}
	}

	return "", ""
}
