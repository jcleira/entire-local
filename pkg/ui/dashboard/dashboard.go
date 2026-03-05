// Package dashboard implements the Bubble Tea TUI for browsing checkpoints.
package dashboard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jcleira/entire-local/pkg/checkpoint"
)

type model struct {
	loader     *checkpoint.Loader
	summaries  []checkpoint.CheckpointSummary
	stats      checkpoint.OverviewStats
	detail     *checkpoint.Checkpoint
	transcript []checkpoint.TranscriptEntry

	screen       screen
	cursor       int
	width        int
	height       int
	filter       string
	filtering    bool
	showHelp     bool
	err          error
	detailTab    detailTab
	scrollOffset int
}

func (m model) contentLineEstimate() int {
	switch m.detailTab {
	case tabFiles:
		if m.detail != nil {
			return strings.Count(m.detail.Diff, "\n") + len(parseDiffStats(m.detail.Diff).files) + 10
		}
	case tabPlan:
		if m.detail != nil {
			return strings.Count(m.detail.Plan, "\n") + 10
		}
	case tabTranscript:
		return len(m.transcript) * 5
	}
	return 0
}

func RunDashboard(loader *checkpoint.Loader) error {
	m := model{
		loader: loader,
		screen: screenOverview,
		width:  80,
		height: 24,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m model) Init() tea.Cmd {
	return m.loadCheckpoints()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.filtering {
			return m.handleFilterInput(msg)
		}
		return m.handleKeyPress(msg)

	case checkpointsLoadedMsg:
		m.summaries = msg.summaries
		m.stats = checkpoint.ComputeStats(msg.summaries)
		return m, nil

	case checkpointDetailLoadedMsg:
		m.detail = msg.checkpoint
		return m, nil

	case transcriptLoadedMsg:
		m.transcript = msg.entries
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	if m.showHelp {
		return renderHelp(m.width, m.height)
	}

	if m.err != nil {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(errClr).Render("Error: "+m.err.Error()),
		)
	}

	switch m.screen {
	case screenOverview:
		return m.renderOverviewScreen()
	case screenDetail:
		return m.renderDetailScreen()
	}

	return ""
}

func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.Help):
		m.showHelp = !m.showHelp
		return m, nil

	case key.Matches(msg, keys.Refresh):
		return m, m.loadCheckpoints()

	case key.Matches(msg, keys.Filter):
		if m.screen == screenOverview {
			m.filtering = true
			return m, nil
		}
	}

	switch m.screen {
	case screenOverview:
		return m.handleOverviewKeys(msg)
	case screenDetail:
		return m.handleDetailKeys(msg)
	}

	return m, nil
}

func (m model) handleOverviewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	filtered := filterSummaries(m.summaries, m.filter)

	switch {
	case key.Matches(msg, keys.Down):
		if m.cursor < len(filtered)-1 {
			m.cursor++
		}

	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}

	case key.Matches(msg, keys.Enter):
		if len(filtered) > 0 && m.cursor < len(filtered) {
			selected := filtered[m.cursor]
			m.screen = screenDetail
			m.detail = nil
			m.transcript = nil
			m.detailTab = tabTranscript
			m.scrollOffset = 0
			return m, tea.Batch(
				m.loadDetail(selected.ID),
				m.loadTranscript(selected.ID),
			)
		}
	}

	return m, nil
}

func (m model) handleDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Back):
		m.screen = screenOverview
		m.detail = nil
		m.transcript = nil
		m.scrollOffset = 0
		return m, nil

	case key.Matches(msg, keys.Tab):
		m.detailTab = (m.detailTab + 1) % 3
		m.scrollOffset = 0
		return m, nil

	case key.Matches(msg, keys.ShiftTab):
		m.detailTab = (m.detailTab + 2) % 3
		m.scrollOffset = 0
		return m, nil

	case key.Matches(msg, keys.Num1):
		m.detailTab = tabTranscript
		m.scrollOffset = 0
		return m, nil

	case key.Matches(msg, keys.Num2):
		m.detailTab = tabFiles
		m.scrollOffset = 0
		return m, nil

	case key.Matches(msg, keys.Num3):
		m.detailTab = tabPlan
		m.scrollOffset = 0
		return m, nil

	case key.Matches(msg, keys.Down):
		if m.scrollOffset < m.contentLineEstimate() {
			m.scrollOffset++
		}
		return m, nil

	case key.Matches(msg, keys.Up):
		if m.scrollOffset > 0 {
			m.scrollOffset--
		}
		return m, nil
	}

	return m, nil
}

func (m model) handleFilterInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter, tea.KeyEscape:
		m.filtering = false
		if msg.Type == tea.KeyEscape {
			m.filter = ""
		}
		m.cursor = 0
		return m, nil
	case tea.KeyBackspace:
		if m.filter != "" {
			m.filter = m.filter[:len(m.filter)-1]
		}
		m.cursor = 0
		return m, nil
	default:
		if len(msg.Runes) > 0 {
			m.filter += string(msg.Runes)
			m.cursor = 0
		}
		return m, nil
	}
}

func (m model) renderOverviewScreen() string {
	panelWidth := m.width - 2

	overview := renderOverview(m.stats, panelWidth-2)
	overviewPanel := panelStyle.Width(panelWidth).Render(overview)
	overviewHeight := lipgloss.Height(overviewPanel)

	footerContent := m.renderFooterContent()
	footerPanel := panelStyle.Width(panelWidth).Render(footerContent)
	footerHeight := lipgloss.Height(footerPanel)

	listPanelOuter := m.height - overviewHeight - footerHeight
	if listPanelOuter < 5 {
		listPanelOuter = 5
	}
	listInnerHeight := listPanelOuter - 2
	if listInnerHeight < 3 {
		listInnerHeight = 3
	}

	contentWidth := panelWidth - 2
	list := renderCheckpointList(m.summaries, m.cursor, listInnerHeight, m.filter, contentWidth)
	listPanel := panelStyle.Width(panelWidth).Height(listPanelOuter - 2).Render(list)

	return lipgloss.JoinVertical(lipgloss.Left, overviewPanel, listPanel, footerPanel)
}

func (m model) renderDetailScreen() string {
	if m.detail == nil {
		return dimStyle.Render("Loading checkpoint...")
	}

	panelWidth := m.width - 2

	header := renderDetailHeader(m.detail, panelWidth-2)
	headerPanel := panelStyle.Width(panelWidth).Render(header)
	headerHeight := lipgloss.Height(headerPanel)

	footerContent := dimStyle.Render("esc back  j/k scroll  tab switch  1-3 tabs  ? help  q quit")
	footerPanel := panelStyle.Width(panelWidth).Render(footerContent)
	footerHeight := lipgloss.Height(footerPanel)

	contentOuter := m.height - headerHeight - footerHeight
	if contentOuter < 5 {
		contentOuter = 5
	}
	contentInner := contentOuter - 2
	if contentInner < 3 {
		contentInner = 3
	}

	tabBar := renderTabBar(m.detailTab)
	tabBarHeight := lipgloss.Height(tabBar) + 1

	tabContentHeight := contentInner - tabBarHeight
	innerWidth := panelWidth - 2

	var tabContent string
	switch m.detailTab {
	case tabFiles:
		tabContent = renderTabFiles(m.detail.Diff, innerWidth, tabContentHeight)
	case tabTranscript:
		tabContent = renderTabTranscript(m.transcript, innerWidth, tabContentHeight, m.scrollOffset)
	case tabPlan:
		tabContent = renderTabPlan(m.detail.Plan, innerWidth, tabContentHeight)
	}

	if m.detailTab != tabTranscript {
		lines := strings.Split(tabContent, "\n")
		maxScroll := len(lines) - tabContentHeight
		if maxScroll < 0 {
			maxScroll = 0
		}
		scrollOffset := m.scrollOffset
		if scrollOffset > maxScroll {
			scrollOffset = maxScroll
		}
		end := scrollOffset + tabContentHeight
		if end > len(lines) {
			end = len(lines)
		}
		if scrollOffset < len(lines) {
			tabContent = strings.Join(lines[scrollOffset:end], "\n")
		}
	}

	combined := lipgloss.JoinVertical(lipgloss.Left, tabBar, "", tabContent)
	contentPanel := panelStyle.Width(panelWidth).Height(contentOuter - 2).Render(combined)

	return lipgloss.JoinVertical(lipgloss.Left, headerPanel, contentPanel, footerPanel)
}

func (m model) renderFooterContent() string {
	if m.filtering {
		return fmt.Sprintf("Filter: %s█", m.filter)
	}

	leftSide := fmt.Sprintf("%d checkpoints", m.stats.TotalCheckpoints)
	if m.stats.TotalTokens > 0 {
		leftSide += "  ·  " + formatTokensShort(m.stats.TotalTokens)
	}

	rightSide := "j/k navigate  enter select  / filter  ? help  q quit"

	padding := m.width - 6 - lipgloss.Width(leftSide) - lipgloss.Width(rightSide)
	if padding < 2 {
		padding = 2
	}

	return leftSide + strings.Repeat(" ", padding) + dimStyle.Render(rightSide)
}

func (m model) loadCheckpoints() tea.Cmd {
	return func() tea.Msg {
		summaries, err := m.loader.ListCheckpoints()
		if err != nil {
			return errMsg{err}
		}
		return checkpointsLoadedMsg{summaries}
	}
}

func (m model) loadDetail(id string) tea.Cmd {
	return func() tea.Msg {
		cp, err := m.loader.LoadCheckpoint(id)
		if err != nil {
			return errMsg{err}
		}
		return checkpointDetailLoadedMsg{cp}
	}
}

func (m model) loadTranscript(id string) tea.Cmd {
	return func() tea.Msg {
		entries, err := m.loader.LoadTranscript(id)
		if err != nil {
			return errMsg{err}
		}
		return transcriptLoadedMsg{entries}
	}
}
