// Package text contains basic text manipulations for displaying text in the command line
// from the use of the lipgloss package
package text

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Text conatins all the functions related to text manipulations
type Text struct{}

// Style is a struct that contains the styling properties
type Style struct {
	Color   lipgloss.TerminalColor
	Padding P
	Align   lipgloss.Position
	Bold    bool
}

// P is a struct that is used to represent the padding
type P struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

// H heading
func (Text) H(style Style, strs ...string,
) string {
	return lipgloss.NewStyle().
		Bold(style.Bold).
		Foreground(style.Color).
		PaddingTop(style.Padding.Top).
		PaddingLeft(style.Padding.Left).
		PaddingRight(style.Padding.Right).
		PaddingBottom(style.Padding.Bottom).
		Align(style.Align).
		Render(strs...)
}

// P paragraph
func (Text) P(style Style, strs ...string,
) string {
	return lipgloss.NewStyle().
		Bold(style.Bold).
		Foreground(style.Color).
		PaddingTop(style.Padding.Top).
		PaddingLeft(style.Padding.Left).
		PaddingRight(style.Padding.Right).
		PaddingBottom(style.Padding.Bottom).
		Align(style.Align).
		Render(strs...)
}

// Error is a function that is used to display errors to the standered output
func (Text) Error(strs ...string) {
	fmt.Println(
		Text{}.P(Style{
			Bold:  false,
			Color: lipgloss.Color("#D72023"),
			Padding: P{
				Left: 1,
				Top:  1,
			},
			Align: lipgloss.Left,
		}, strs...),
	)
}

// ErrorWithPadding is a function that is used to display the error messages with
// the ability to control the pading
func (Text) ErrorWithPadding(padding P, strs ...string) {
	fmt.Println(
		Text{}.P(Style{
			Bold:    false,
			Color:   lipgloss.Color("#D72023"),
			Padding: padding,
			Align:   lipgloss.Left,
		}, strs...),
	)
}
