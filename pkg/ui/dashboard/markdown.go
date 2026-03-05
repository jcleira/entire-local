package dashboard

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
)

func markdownStyle() glamour.TermRendererOption {
	s := styles.DarkStyleConfig
	s.H1.Prefix = ""
	s.H2.Prefix = ""
	s.H3.Prefix = ""
	s.H4.Prefix = ""
	s.H5.Prefix = ""
	s.H6.Prefix = ""
	s.Document.Margin = nil
	return glamour.WithStyles(s)
}

func renderMarkdown(content string, width int) string {
	r, err := glamour.NewTermRenderer(
		markdownStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return wordWrap(content, width)
	}

	out, err := r.Render(content)
	if err != nil {
		return wordWrap(content, width)
	}

	return strings.TrimRight(out, "\n")
}
