// Package checkpoint loads and parses entire.io checkpoint data from git.
package checkpoint

import "time"

type TokenUsage struct {
	InputTokens         int `json:"input_tokens"`
	CacheCreationTokens int `json:"cache_creation_tokens"`
	CacheReadTokens     int `json:"cache_read_tokens"`
	OutputTokens        int `json:"output_tokens"`
	APICallCount        int `json:"api_call_count"`
}

type SessionRef struct {
	Metadata    string `json:"metadata"`
	Transcript  string `json:"transcript"`
	Context     string `json:"context"`
	ContentHash string `json:"content_hash"`
	Prompt      string `json:"prompt"`
}

type Attribution struct {
	AgentLines      int     `json:"agent_lines"`
	HumanAdded      int     `json:"human_added"`
	HumanModified   int     `json:"human_modified"`
	TotalCommitted  int     `json:"total_committed"`
	AgentPercentage float64 `json:"agent_percentage"`
}

type Metadata struct {
	CheckpointID     string       `json:"checkpoint_id"`
	CLIVersion       string       `json:"cli_version"`
	Strategy         string       `json:"strategy"`
	Branch           string       `json:"branch"`
	CheckpointsCount int          `json:"checkpoints_count"`
	FilesTouched     []string     `json:"files_touched"`
	Sessions         []SessionRef `json:"sessions"`
	TokenUsage       TokenUsage   `json:"token_usage"`
}

type SessionMetadata struct {
	CheckpointID string      `json:"checkpoint_id"`
	SessionID    string      `json:"session_id"`
	Strategy     string      `json:"strategy"`
	CreatedAt    string      `json:"created_at"`
	Branch       string      `json:"branch"`
	FilesTouched []string    `json:"files_touched"`
	Agent        string      `json:"agent"`
	TurnID       string      `json:"turn_id"`
	TokenUsage   TokenUsage  `json:"token_usage"`
	Attribution  Attribution `json:"initial_attribution"`
}

type Checkpoint struct {
	Metadata
	Session SessionMetadata
	Prompt  string
	Context string
	Diff    string
	Plan    string
}

type CheckpointSummary struct {
	ID            string
	CreatedAt     time.Time
	Branch        string
	Agent         string
	AgentPercent  int
	CommitHash    string
	CommitMessage string
	Author        string
	Context       string
	Prompt        string
	TotalTokens   int
	Duration      string
	FileCount     int
	Additions     int
	Deletions     int
	SessionCount  int
}

type TranscriptEntry struct {
	Type      string
	Role      string
	Content   string
	ToolName  string
	Timestamp time.Time
	Tokens    int
}
