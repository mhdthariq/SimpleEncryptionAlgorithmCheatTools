package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// ANSI Color codes
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorBold    = "\033[1m"

	// Bright colors
	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightBlue    = "\033[94m"
	ColorBrightMagenta = "\033[95m"
	ColorBrightCyan    = "\033[96m"
)

// Color helper functions
func Colorize(text, color string) string {
	return color + text + ColorReset
}

func Bold(text string) string {
	return ColorBold + text + ColorReset
}

func Red(text string) string {
	return ColorRed + text + ColorReset
}

func Green(text string) string {
	return ColorGreen + text + ColorReset
}

func Yellow(text string) string {
	return ColorYellow + text + ColorReset
}

func Blue(text string) string {
	return ColorBlue + text + ColorReset
}

func Cyan(text string) string {
	return ColorCyan + text + ColorReset
}

func Magenta(text string) string {
	return ColorMagenta + text + ColorReset
}

func BrightGreen(text string) string {
	return ColorBrightGreen + text + ColorReset
}

func BrightYellow(text string) string {
	return ColorBrightYellow + text + ColorReset
}

func BrightCyan(text string) string {
	return ColorBrightCyan + text + ColorReset
}

func BrightMagenta(text string) string {
	return ColorBrightMagenta + text + ColorReset
}

func BrightRed(text string) string {
	return ColorBrightRed + text + ColorReset
}

func BrightBlue(text string) string {
	return ColorBrightBlue + text + ColorReset
}

// Rainbow text effect
func Rainbow(text string) string {
	colors := []string{ColorRed, ColorYellow, ColorGreen, ColorCyan, ColorBlue, ColorMagenta}
	result := ""
	for i, char := range text {
		result += colors[i%len(colors)] + string(char) + ColorReset
	}
	return result
}

// Color combinations
func GreenBold(text string) string {
	return ColorBold + ColorGreen + text + ColorReset
}

func RedBold(text string) string {
	return ColorBold + ColorRed + text + ColorReset
}

func YellowBold(text string) string {
	return ColorBold + ColorYellow + text + ColorReset
}

func CyanBold(text string) string {
	return ColorBold + ColorCyan + text + ColorReset
}

func MagentaBold(text string) string {
	return ColorBold + ColorMagenta + text + ColorReset
}

func BlueBold(text string) string {
	return ColorBold + ColorBlue + text + ColorReset
}

func BrightGreenBold(text string) string {
	return ColorBold + ColorBrightGreen + text + ColorReset
}

func BrightYellowBold(text string) string {
	return ColorBold + ColorBrightYellow + text + ColorReset
}

func BrightCyanBold(text string) string {
	return ColorBold + ColorBrightCyan + text + ColorReset
}

func BrightMagentaBold(text string) string {
	return ColorBold + ColorBrightMagenta + text + ColorReset
}

func BrightRedBold(text string) string {
	return ColorBold + ColorBrightRed + text + ColorReset
}

func BrightBlueBold(text string) string {
	return ColorBold + ColorBrightBlue + text + ColorReset
}

// Decorative lines - responsive (no fixed width)
func PrintColorfulLine(length int) {
	// If length is too large (like 80), reduce it for mobile
	if length > 50 {
		length = 50
	}
	colors := []string{ColorRed, ColorYellow, ColorGreen, ColorCyan, ColorBlue, ColorMagenta}
	line := ""
	for i := 0; i < length; i++ {
		line += colors[i%len(colors)] + "═" + ColorReset
	}
	fmt.Println(line)
}

func PrintStarLine(length int) {
	// Reduce for mobile
	if length > 20 {
		length = 20
	}
	fmt.Println(BrightYellow(strings.Repeat("★ ", length/2)))
}

func PrintSectionHeader(title string) {
	// Colorful header - responsive
	PrintColorfulLine(50)
	fmt.Println(BrightMagentaBold(strings.ToUpper(title)))
	PrintColorfulLine(50)
}

// PrintSeparator - responsive separator line
func PrintSeparator() {
	fmt.Println(BrightCyan("─────────────────────────────────────────────"))
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}

func GetTableWriter() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	return t
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func GetInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
