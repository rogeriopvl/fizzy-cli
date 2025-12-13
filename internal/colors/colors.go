package colors

import "github.com/charmbracelet/lipgloss"

type ColorDef struct {
	Name      string
	CSSValue  string
	TermColor lipgloss.Color
}

var (
	Blue   = ColorDef{"Blue", "var(--color-card-default)", lipgloss.Color("12")}
	Gray   = ColorDef{"Gray", "var(--color-card-1)", lipgloss.Color("8")}
	Tan    = ColorDef{"Tan", "var(--color-card-2)", lipgloss.Color("180")}
	Yellow = ColorDef{"Yellow", "var(--color-card-3)", lipgloss.Color("11")}
	Lime   = ColorDef{"Lime", "var(--color-card-4)", lipgloss.Color("10")}
	Aqua   = ColorDef{"Aqua", "var(--color-card-5)", lipgloss.Color("14")}
	Violet = ColorDef{"Violet", "var(--color-card-6)", lipgloss.Color("177")}
	Purple = ColorDef{"Purple", "var(--color-card-7)", lipgloss.Color("135")}
	Pink   = ColorDef{"Pink", "var(--color-card-8)", lipgloss.Color("205")}
)

var All = []ColorDef{Blue, Gray, Tan, Yellow, Lime, Aqua, Violet, Purple, Pink}

func ByName(name string) *ColorDef {
	for _, c := range All {
		if c.Name == name {
			return &c
		}
	}
	return nil
}
