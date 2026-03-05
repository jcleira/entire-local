package dashboard

import "github.com/jcleira/entire-local/pkg/checkpoint"

type screen int

const (
	screenOverview screen = iota
	screenDetail
)

type detailTab int

const (
	tabTranscript detailTab = iota
	tabFiles
	tabPlan
)

type checkpointsLoadedMsg struct {
	summaries []checkpoint.CheckpointSummary
}

type checkpointDetailLoadedMsg struct {
	checkpoint *checkpoint.Checkpoint
}

type transcriptLoadedMsg struct {
	entries []checkpoint.TranscriptEntry
}

type errMsg struct {
	err error
}
