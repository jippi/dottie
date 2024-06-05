package tui

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type StyleChanger func(*lipgloss.Style)

var Bold = func(s *lipgloss.Style) {
	s.Bold(true)
}

type PrinterOption func(p *Printer)

// Printer mirrors the [fmt] package print/sprint functions, wraps them in a [lipgloss.Style]
// and an optional [WordWrap] configuration with a configured [BoxWidth].
//
// Additionally, [Printer*] methods writes to the configured [Writer] instead of [os.Stdout]
type Printer struct {
	boxWidth       int                // Max width for strings when using WrapMode
	writer         io.Writer          // Writer controls where implicit print output goes for [Print], [Printf], [Printfln] and [Println]
	renderer       *lipgloss.Renderer // The renderer responsible for providing the output and color management
	style          Style              // Style config
	textStyle      lipgloss.Style
	boxHeaderStyle lipgloss.Style
	boxBodyStyle   lipgloss.Style
}

func NewPrinter(style Style, renderer *lipgloss.Renderer, options ...PrinterOption) Printer {
	options = append([]PrinterOption{
		WitBoxWidth(80),
		WithStyle(style),
		WithRenderer(renderer),
	}, options...)

	printer := &Printer{}
	for _, option := range options {
		option(printer)
	}

	printer.boxHeaderStyle = style.BoxHeader()
	printer.boxBodyStyle = style.BoxBody()

	return *printer
}

// ----------------------------------------
// print to a specific io.Writer
// ----------------------------------------

// Fprint mirrors [fmt.Fprint] signature and behavior, with the configured style
// and (optional) word wrapping applied
func (p Printer) Fprint(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprint(w, p.Sprint(a...))
}

// Fprintf mirrors [fmt.Fprintf] signature and behavior, with the configured style
// and (optional) word wrapping applied
func (p Printer) Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
	return p.Fprint(w, p.Sprintf(format, a...))
}

// Fprintfln mirrors [fmt.Fprintfln] signature and behavior, with the configured style
// and (optional) word wrapping applied
func (p Printer) Fprintfln(w io.Writer, format string, a ...any) (n int, err error) {
	return p.Fprintln(w, p.Sprintf(format, a...))
}

// Fprintln mirrors [fmt.Fprintln] signature and behavior, with the configured style
// and (optional) word wrapping applied
func (p Printer) Fprintln(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprintln(w, p.printHelper(a...))
}

// -----------------------------------------------------
// Print to the default [p.writer] over [os.Stdout]
// -----------------------------------------------------

// Print mirrors [fmt.Print] signature and behavior, with the configured style
// and (optional) word wrapping applied.
//
// Instead of writing to [os.Stdout] it will write to the configured [io.Writer].
func (p Printer) Print(a ...any) (n int, err error) {
	return p.Fprint(p.writer, a...)
}

// Printf mirrors [fmt.Printf] signature and behavior, with the configured style
// and (optional) word wrapping applied.
//
// Instead of writing to [os.Stdout] it will write to the configured [io.Writer].
func (p Printer) Printf(format string, a ...any) (n int, err error) {
	return p.Fprintf(p.writer, format, a...)
}

// Printfln behaves like [fmt.Printf] but supports the [formatter] signature.
//
// This does *not* map to a Go native printer, but a mix for formatting + newline
func (p Printer) Printfln(format string, a ...any) (n int, err error) {
	return p.Fprintfln(p.writer, format, a...)
}

// Println mirrors [fmt.Println] signature and behavior, with the configured style
// and (optional) word wrapping applied.
//
// Instead of writing to [os.Stdout] it will write to the configured [io.Writer]
func (p Printer) Println(a ...any) (n int, err error) {
	return p.Fprintln(p.writer, a...)
}

// -----------------------------------------------------
// Return string
// -----------------------------------------------------

// Sprint mirrors [fmt.Sprint] signature and behavior, with the configured style
// and (optional) word wrapping applied.
func (p Printer) Sprint(a ...any) string {
	return p.render(fmt.Sprint(a...))
}

// Sprintf mirrors [fmt.Sprintf] signature and behavior, with the configured style
// and (optional) word wrapping applied.
func (p Printer) Sprintf(format string, a ...any) string {
	return p.render(fmt.Sprintf(format, a...))
}

// Sprintfln behaves like [fmt.Sprintln] but supports the [formatter] signature.
//
// This does *not* map to a Go native printer, but a mix for formatting + newline
func (p Printer) Sprintfln(format string, a ...any) string {
	return fmt.Sprintln(p.Sprintf(format, a...))
}

// Sprintln mirrors [fmt.Sprintln] signature and behavior, with the configured style
// and (optional) word wrapping applied.
func (p Printer) Sprintln(a ...any) string {
	return fmt.Sprintln(p.printHelper(a...))
}

// Create a visual box with the printer style
func (p Printer) Box(header string, bodies ...string) {
	body := strings.Join(bodies, " ")

	// Copy the box styles to avoid leaking changes to the styles
	headerStyle, bodyStyle := p.boxHeaderStyle, p.boxBodyStyle

	// If there are no body, just render the header box directly
	if len(body) == 0 {
		fmt.Fprintln(
			p.writer,
			headerStyle.
				Width(p.boxWidth-borderWidth).
				Border(headerOnlyBorder).
				Render(header),
		)

		return
	}

	// Render the header and body box
	boxHeader := headerStyle.Width(p.boxWidth - borderWidth).Render(header)
	boxBody := bodyStyle.Width(p.boxWidth - borderWidth).Render(body)

	// If a BoxWidth is set, the boxes will be aligned automatically to the max
	if p.boxWidth > 0 {
		fmt.Fprintln(
			p.writer,
			lipgloss.JoinVertical(
				lipgloss.Left,
				boxHeader,
				boxBody,
			),
		)

		return
	}

	// Compute the width of the header and body elements
	headerWidth := lipgloss.Width(boxHeader) - borderWidth
	bodyWidth := lipgloss.Width(boxBody) - borderWidth

	// Find the shortest box and (re)render it to the length of the longest one
	switch {
	case headerWidth > bodyWidth:
		boxBody = bodyStyle.Width(headerWidth).Render(body)

	case headerWidth < bodyWidth:
		boxHeader = headerStyle.Width(bodyWidth).Render(header)
	}

	fmt.Fprintln(
		p.writer,
		lipgloss.JoinVertical(lipgloss.Left, boxHeader, boxBody),
	)
}

// -----------------------------------------------------
// io.Writer
// -----------------------------------------------------

func (p Printer) Write(b []byte) (n int, err error) {
	return p.Print(string(b))
}

// -----------------------------------------------------
// Helper methods
// -----------------------------------------------------

// GetBoxWidth returns the configured [BoxWidth] for word wrapping
func (p Printer) BoxWidth() int {
	return p.boxWidth
}

// Writer returns the configured [io.Writer]
func (p Printer) Writer() io.Writer {
	return p.writer
}

func (p Printer) Copy(options ...PrinterOption) Printer {
	clone := &p

	for _, option := range options {
		option(clone)
	}

	return *clone
}

// TextStyle returns a *copy* of the current [lipgloss.Style]
func (p Printer) Style() lipgloss.Style {
	return p.textStyle
}

// ApplyTextStyle returns a new copy of [StylePrint] instance with the [Style] based on the callback changes
func (p Printer) ApplyStyle(callback StyleChanger) Printer {
	style := p.Style()
	callback(&style)

	return p.Copy(WithTextStyle(style))
}

func (p Printer) GetWriter() io.Writer {
	return p.writer
}

// -----------------------------------------------------
// internal helpers
// -----------------------------------------------------

func (p Printer) render(input string) string {
	return p.wrap(p.textStyle.Render(input))
}

func (p Printer) wrap(input string) string {
	return input
}

func (p Printer) printHelper(a ...any) string {
	var buff bytes.Buffer

	fmt.Fprintln(&buff, a...)

	out := buff.String()
	out, _ = strings.CutSuffix(out, "\n")

	return p.render(out)
}

// -----------------------------------------------------
// Printer options
// -----------------------------------------------------

func WithStyle(style Style) PrinterOption {
	return func(p *Printer) {
		p.style = style
		p.textStyle = p.renderer.NewStyle().Inherit(style.TextStyle())
	}
}

func WithRenderer(renderer *lipgloss.Renderer) PrinterOption {
	return func(p *Printer) {
		p.renderer = renderer
		p.writer = renderer.Output()
	}
}

func WithTextStyle(style lipgloss.Style) PrinterOption {
	return func(p *Printer) {
		p.textStyle = style
	}
}

func WithEmphasis(b bool) PrinterOption {
	return func(printer *Printer) {
		if b {
			printer.textStyle = printer.renderer.NewStyle().Inherit(printer.style.TextEmphasisStyle())

			return
		}

		printer.textStyle = printer.renderer.NewStyle().Inherit(printer.style.TextStyle())
	}
}

func WithWriter(w io.Writer) PrinterOption {
	return func(p *Printer) {
		p.writer = w
	}
}

func WitBoxWidth(i int) PrinterOption {
	return func(p *Printer) {
		p.boxWidth = i
	}
}
