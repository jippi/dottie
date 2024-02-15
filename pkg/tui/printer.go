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

type PrinterOption func(p *Print)

// Printer mirrors the [fmt] package print/sprint functions, wraps them in a [lipgloss.Style]
// and an optional [WordWrap] configuration with a configured [MaxWidth].
//
// Additionally, [Print*] methods writes to the configured [Writer] instead of [os.Stdout]
type Printer interface {
	// ----------------------------------------
	// print to a specific io.Writer
	// ----------------------------------------

	// Fprint mirrors [fmt.Fprint] signature and behavior, with the configured style
	// and (optional) word wrapping applied
	Fprint(w io.Writer, a ...any) (n int, err error)

	// Fprintf mirrors [fmt.Fprintf] signature and behavior, with the configured style
	// and (optional) word wrapping applied
	Fprintf(w io.Writer, format string, a ...any) (n int, err error)

	// Fprintfln mirrors [fmt.Fprintfln] signature and behavior, with the configured style
	// and (optional) word wrapping applied
	Fprintfln(w io.Writer, format string, a ...any) (n int, err error)

	// Fprintln mirrors [fmt.Fprintln] signature and behavior, with the configured style
	// and (optional) word wrapping applied
	Fprintln(w io.Writer, a ...any) (n int, err error)

	// ----------------------------------------
	// print to the default io.Writer
	// ----------------------------------------

	// Print mirrors [fmt.Print] signature and behavior, with the configured style
	// and (optional) word wrapping applied.
	//
	// Instead of writing to [os.Stdout] it will write to the configured [io.Writer].
	Print(a ...any) (n int, err error)

	// Printf mirrors [fmt.Printf] signature and behavior, with the configured style
	// and (optional) word wrapping applied.
	//
	// Instead of writing to [os.Stdout] it will write to the configured [io.Writer].
	Printf(format string, a ...any) (n int, err error)

	// Printfln behaves like [fmt.Printf] but supports the [formatter] signature.
	//
	// This does *not* map to a Go native printer, but a mix for formatting + newline
	Printfln(format string, a ...any) (n int, err error)

	// Println mirrors [fmt.Println] signature and behavior, with the configured style
	// and (optional) word wrapping applied.
	//
	// Instead of writing to [os.Stdout] it will write to the configured [io.Writer].
	Println(a ...any) (n int, err error)

	// ----------------------------------------
	// return string
	// ----------------------------------------

	// Sprint mirrors [fmt.Sprint] signature and behavior, with the configured style
	// and (optional) word wrapping applied.
	Sprint(a ...any) string

	// Sprintf mirrors [fmt.Sprintf] signature and behavior, with the configured style
	// and (optional) word wrapping applied.
	Sprintf(format string, a ...any) string

	// Sprintfln behaves like [fmt.Sprintln] but supports the [formatter] signature.
	//
	// This does *not* map to a Go native printer, but a mix for formatting + newline
	Sprintfln(format string, a ...any) string

	// Sprintln mirrors [fmt.Sprintln] signature and behavior, with the configured style
	// and (optional) word wrapping applied.
	Sprintln(a ...any) string

	// ----------------------------------------
	// helper methods
	// ----------------------------------------

	Copy(options ...PrinterOption) Print

	// GetMaxWidth returns the configured [MaxWidth] for word wrapping
	MaxWidth() int

	// TextStyle returns a *copy* of the current [lipgloss.Style]
	Style() lipgloss.Style

	// ApplyTextStyle returns a new copy of [StylePrint] instance with the [Style] based on the callback changes
	ApplyStyle(changer StyleChanger) Print

	// Writer returns the configured [io.Writer]
	Writer() io.Writer

	// Create a visual box with the printer style
	Box(header string, body ...string)
}

type Print struct {
	boxWidth int                // Max width for strings when using WrapMode
	writer   io.Writer          // Writer controls where implicit print output goes for [Print], [Printf], [Printfln] and [Println]
	renderer *lipgloss.Renderer // The renderer responsible for providing the output and color management
	color    Color              // Color config

	textStyle    lipgloss.Style
	textEmphasis bool
	boxStyles    Box
}

func NewPrinter(color Color, renderer *lipgloss.Renderer, options ...PrinterOption) Print {
	options = append([]PrinterOption{
		WitBoxWidth(100),
		WithColor(color),
		WithRenderer(renderer),
		WithEmphasis(false),
	}, options...)

	printer := &Print{}
	for _, option := range options {
		option(printer)
	}

	printer.boxStyles = printer.color.BoxStyles(
		printer.renderer.NewStyle(),
		printer.renderer.NewStyle(),
	)

	return *printer
}

// -----------------------------------------------------
// Print to a user-provided io.Writer
// -----------------------------------------------------

func (p Print) Fprint(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprint(w, p.Sprint(a...))
}

func (p Print) Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
	return p.Fprint(w, p.Sprintf(format, a...))
}

func (p Print) Fprintfln(w io.Writer, format string, a ...any) (n int, err error) {
	return p.Fprintln(w, p.Sprintf(format, a...))
}

func (p Print) Fprintln(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprintln(w, p.printHelper(a...))
}

// -----------------------------------------------------
// Print to the default [p.writer] over [os.Stdout]
// -----------------------------------------------------

func (p Print) Print(a ...any) (n int, err error) {
	return p.Fprint(p.writer, a...)
}

func (p Print) Printf(format string, a ...any) (n int, err error) {
	return p.Fprintf(p.writer, format, a...)
}

func (p Print) Printfln(format string, a ...any) (n int, err error) {
	return p.Fprintfln(p.writer, format, a...)
}

func (p Print) Println(a ...any) (n int, err error) {
	return p.Fprintln(p.writer, a...)
}

// -----------------------------------------------------
// Return string
// -----------------------------------------------------

func (p Print) Sprint(a ...any) string {
	return p.render(fmt.Sprint(a...))
}

func (p Print) Sprintf(format string, a ...any) string {
	return p.render(fmt.Sprintf(format, a...))
}

func (p Print) Sprintfln(format string, a ...any) string {
	return fmt.Sprintln(p.Sprintf(format, a...))
}

func (p Print) Sprintln(a ...any) string {
	return fmt.Sprintln(p.printHelper(a...))
}

func (p Print) Box(header string, bodies ...string) {
	body := strings.Join(bodies, " ")

	// Copy the box styles to avoid leaking changes to the styles
	styles := p.boxStyles.Copy()

	// If there are no body, just render the header box directly
	if len(body) == 0 {
		fmt.Fprintln(
			p.writer,
			styles.Header.
				Width(p.boxWidth-borderWidth).
				Border(headerOnlyBorder).
				Render(header),
		)

		return
	}

	// Render the header and body box
	boxHeader := styles.Header.Width(p.boxWidth - borderWidth).Render(header)
	boxBody := styles.Body.Width(p.boxWidth - borderWidth).Render(body)

	// If a maxWidth is set, the boxes will be aligned automatically to the max
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
		boxBody = styles.Body.Width(headerWidth).Render(body)

	case headerWidth < bodyWidth:
		boxHeader = styles.Header.Width(bodyWidth).Render(header)
	}

	fmt.Fprintln(
		p.writer,
		lipgloss.JoinVertical(lipgloss.Left, boxHeader, boxBody),
	)
}

// -----------------------------------------------------
// io.Writer
// -----------------------------------------------------

func (p Print) Write(b []byte) (n int, err error) {
	return p.Print(string(b))
}

// -----------------------------------------------------
// Helper methods
// -----------------------------------------------------

func (p Print) MaxWidth() int {
	return p.boxWidth
}

func (p Print) Writer() io.Writer {
	return p.writer
}

func (p Print) Copy(options ...PrinterOption) Print {
	clone := &p

	for _, option := range options {
		option(clone)
	}

	return *clone
}

func (p Print) Style() lipgloss.Style {
	return p.textStyle.Copy()
}

func (p Print) ApplyStyle(callback StyleChanger) Print {
	style := p.Style()
	callback(&style)

	return p.Copy(WithTextStyle(style))
}

func (p Print) GetWriter() io.Writer {
	return p.writer
}

// -----------------------------------------------------
// internal helpers
// -----------------------------------------------------

func (p Print) render(input string) string {
	return p.wrap(p.textStyle.Render(input))
}

func (p Print) wrap(input string) string {
	return input
}

func (p Print) printHelper(a ...any) string {
	var buff bytes.Buffer

	fmt.Fprintln(&buff, a...)

	out := buff.String()
	out, _ = strings.CutSuffix(out, "\n")

	return p.render(out)
}

// -----------------------------------------------------
// Printer options
// -----------------------------------------------------

func WithColor(color Color) PrinterOption {
	return func(p *Print) {
		p.color = color
	}
}

func WithRenderer(renderer *lipgloss.Renderer) PrinterOption {
	return func(p *Print) {
		p.renderer = renderer
		p.writer = renderer.Output()
	}
}

func WithTextStyle(style lipgloss.Style) PrinterOption {
	return func(p *Print) {
		p.textStyle = style
	}
}

func WithBoxStyle(style Box) PrinterOption {
	return func(p *Print) {
		p.boxStyles = style
	}
}

func WithEmphasis(b bool) PrinterOption {
	return func(printer *Print) {
		printer.textEmphasis = b

		if b {
			printer.textStyle = printer.color.TextEmphasisStyle(printer.renderer.NewStyle())

			return
		}

		printer.textStyle = printer.color.TextStyle(printer.renderer.NewStyle())
	}
}

func WithWriter(w io.Writer) PrinterOption {
	return func(p *Print) {
		p.writer = w
	}
}

func WitBoxWidth(i int) PrinterOption {
	return func(p *Print) {
		p.boxWidth = i
	}
}
