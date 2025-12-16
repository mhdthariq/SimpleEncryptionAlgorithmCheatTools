package algorithms

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"

	"example.com/cheat-encryption-algorithm/pkg/ui"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Vigenere implementation
type Vigenere struct {
	plaintext      string
	key            string
	mode           string // "encrypt" or "decrypt"
	ciphertext     string
	processedKey   string
	alphabet       string
	caseSensitive  bool
	preserveSpaces bool
}

func NewVigenere() *Vigenere {
	return &Vigenere{
		alphabet:       "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		caseSensitive:  false,
		preserveSpaces: true,
	}
}

func (v *Vigenere) GetName() string {
	return "Vigenère Cipher"
}

func (v *Vigenere) Run(reader *bufio.Reader, plaintext string, key string) {
	// Vigenère specific options
	fmt.Println()
	fmt.Print(ui.BrightCyan("Would you like to encrypt or decrypt? (e/d) [default: e]: "))
	modeStr, _ := reader.ReadString('\n')
	modeStr = strings.TrimSpace(strings.ToLower(modeStr))

	if modeStr == "" || modeStr == "e" {
		v.mode = "encrypt"
	} else if modeStr == "d" {
		v.mode = "decrypt"
	} else {
		fmt.Println(ui.Red("Invalid mode selected. Defaulting to encryption."))
		v.mode = "encrypt"
	}

	fmt.Print(ui.BrightCyan("Preserve spaces and punctuation? (y/n) [default: y]: "))
	preserveStr, _ := reader.ReadString('\n')
	preserveStr = strings.TrimSpace(strings.ToLower(preserveStr))
	v.preserveSpaces = preserveStr == "" || preserveStr == "y"

	// Initialize
	v.plaintext = plaintext
	v.key = key

	fmt.Println()
	fmt.Println(ui.BrightGreen("Processing..."))
	fmt.Println()

	ui.PrintSectionHeader("VIGENÈRE CIPHER - POLYALPHABETIC SUBSTITUTION")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightCyanBold("🔐 The Vigenère cipher uses a keyword to create multiple Caesar shifts."))
	fmt.Println(ui.BrightYellow("📚 Named after Blaise de Vigenère (16th century), it was considered unbreakable for centuries!"))
	fmt.Println(ui.BrightMagenta("✨ Each letter of the key determines how much to shift each letter of the plaintext."))
	ui.PrintStarLine(40)
	fmt.Println()

	v.printSetup()
	v.printTabulaRecta()
	v.printInteractiveLookup()
	v.processText()
	v.printStepByStep()
	v.printResults()
	v.printCryptanalysisInfo()
}

func (v *Vigenere) printSetup() {
	ui.PrintSectionHeader("SETUP INFORMATION")
	fmt.Println(ui.BrightCyanBold("📋 Configuration and input parameters:"))
	fmt.Println()

	t := ui.GetTableWriter()

	t.AppendHeader(table.Row{
		ui.MagentaBold("Parameter"),
		ui.CyanBold("Value"),
		ui.YellowBold("Description"),
	})

	modeDisplay := v.mode
	if v.mode == "encrypt" {
		modeDisplay = ui.BrightGreen("ENCRYPTION") + " 🔒"
	} else {
		modeDisplay = ui.BrightBlue("DECRYPTION") + " 🔓"
	}

	inputLabel := "Plaintext"
	if v.mode == "decrypt" {
		inputLabel = "Ciphertext"
	}

	t.AppendRow(table.Row{
		ui.Cyan("Mode"),
		modeDisplay,
		ui.Yellow("Operation to perform"),
	})

	t.AppendRow(table.Row{
		ui.Cyan(inputLabel),
		ui.BrightYellow(fmt.Sprintf("\"%s\"", v.plaintext)),
		ui.Yellow(fmt.Sprintf("Length: %d characters", len(v.plaintext))),
	})

	t.AppendRow(table.Row{
		ui.Cyan("Key"),
		ui.BrightGreen(fmt.Sprintf("\"%s\"", v.key)),
		ui.Yellow(fmt.Sprintf("Length: %d characters", len(v.key))),
	})

	t.AppendRow(table.Row{
		ui.Cyan("Alphabet"),
		ui.BrightMagenta(v.alphabet),
		ui.Yellow(fmt.Sprintf("%d letters", len(v.alphabet))),
	})

	preserveDisplay := ui.Red("No")
	if v.preserveSpaces {
		preserveDisplay = ui.Green("Yes")
	}

	t.AppendRow(table.Row{
		ui.Cyan("Preserve Spacing"),
		preserveDisplay,
		ui.Yellow("Keep spaces & punctuation"),
	})

	t.AppendSeparator()

	// Key preparation info
	cleanKey := v.cleanText(v.key)
	t.AppendRow(table.Row{
		ui.Magenta("Processed Key"),
		ui.BrightCyan(strings.ToUpper(cleanKey)),
		ui.Yellow("Letters only, uppercase"),
	})

	t.Render()
	fmt.Println()
}

func (v *Vigenere) printTabulaRecta() {
	ui.PrintSectionHeader("TABULA RECTA (VIGENÈRE TABLE)")
	fmt.Println(ui.BrightCyanBold("📊 The complete Vigenère square (26×26) - All possible Caesar shifts:"))
	fmt.Println(ui.BrightYellow("   • Each row is shifted one position from the previous"))
	fmt.Println(ui.BrightYellow("   • The key letter determines which ROW to use"))
	fmt.Println(ui.BrightYellow("   • The plaintext letter determines which COLUMN to use"))
	fmt.Println(ui.BrightMagenta("   • The intersection gives you the ciphertext letter!"))
	fmt.Println()

	// Get first few characters for demonstration
	cleanKey := v.cleanText(v.key)
	cleanPlain := v.cleanText(v.plaintext)

	// Collect demo positions (first 3 characters)
	demoPositions := make([]struct {
		plainChar rune
		keyChar   rune
		plainPos  int
		keyPos    int
		cipherPos int
	}, 0)

	if len(cleanKey) > 0 && len(cleanPlain) > 0 {
		for i := 0; i < len(cleanPlain) && i < 3; i++ {
			plainChar := rune(cleanPlain[i])
			keyChar := rune(cleanKey[i%len(cleanKey)])
			plainPos := strings.IndexRune(v.alphabet, plainChar)
			keyPos := strings.IndexRune(v.alphabet, keyChar)

			if plainPos >= 0 && keyPos >= 0 {
				var cipherPos int
				if v.mode == "encrypt" {
					cipherPos = (plainPos + keyPos) % len(v.alphabet)
				} else {
					cipherPos = (plainPos - keyPos + len(v.alphabet)) % len(v.alphabet)
				}

				demoPositions = append(demoPositions, struct {
					plainChar rune
					keyChar   rune
					plainPos  int
					keyPos    int
					cipherPos int
				}{plainChar, keyChar, plainPos, keyPos, cipherPos})
			}
		}
	}

	// Print header with column numbers
	fmt.Print(ui.BrightBlueBold("     "))
	for i := 0; i < len(v.alphabet); i++ {
		if i < 10 {
			fmt.Print(ui.Yellow(fmt.Sprintf("%d ", i)))
		} else {
			fmt.Print(ui.Yellow(fmt.Sprintf("%d", i)))
		}
	}
	fmt.Println()

	// Header row with letters
	fmt.Print(ui.MagentaBold("   │ "))
	for i := 0; i < len(v.alphabet); i++ {
		// Highlight columns used in demo
		highlighted := false
		for _, demo := range demoPositions {
			if demo.plainPos == i {
				fmt.Print(ui.BrightGreenBold(string(v.alphabet[i])) + " ")
				highlighted = true
				break
			}
		}
		if !highlighted {
			fmt.Print(ui.CyanBold(string(v.alphabet[i])) + " ")
		}
	}
	fmt.Println(ui.Yellow(" ← Plaintext (column)"))

	// Separator
	fmt.Print(ui.MagentaBold("───┼"))
	for i := 0; i < len(v.alphabet); i++ {
		fmt.Print("──")
	}
	fmt.Println()

	// Print all 26 rows
	for row := 0; row < len(v.alphabet); row++ {
		// Row header with highlighting
		rowHighlighted := false
		for _, demo := range demoPositions {
			if demo.keyPos == row {
				fmt.Print(ui.BrightMagentaBold(fmt.Sprintf(" %c │ ", v.alphabet[row])))
				rowHighlighted = true
				break
			}
		}
		if !rowHighlighted {
			fmt.Print(ui.MagentaBold(fmt.Sprintf(" %c │ ", v.alphabet[row])))
		}

		// Row cells
		for col := 0; col < len(v.alphabet); col++ {
			shiftedIdx := (col + row) % len(v.alphabet)
			char := string(v.alphabet[shiftedIdx])

			// Check if this cell is used in encryption demo
			isDemo := false
			demoNum := 0
			for idx, demo := range demoPositions {
				if demo.keyPos == row && demo.plainPos == col {
					isDemo = true
					demoNum = idx + 1
					break
				}
			}

			if isDemo {
				// Highlight the actual encryption cells with background effect
				if demoNum == 1 {
					fmt.Print(ui.BrightRedBold(char) + ui.Red("●"))
				} else if demoNum == 2 {
					fmt.Print(ui.BrightYellowBold(char) + ui.Yellow("●"))
				} else {
					fmt.Print(ui.BrightCyanBold(char) + ui.Cyan("●"))
				}
			} else if row == col {
				// Diagonal (same letter - no shift)
				fmt.Print(ui.Yellow(char) + " ")
			} else if demo := findDemoInRow(row, col, demoPositions); demo != nil {
				// Same row or column as demo
				fmt.Print(ui.BrightCyan(char) + " ")
			} else {
				// Regular cell
				fmt.Print(ui.Cyan(char) + " ")
			}
		}

		// Row annotation
		if rowHighlighted {
			fmt.Print(ui.BrightMagenta(fmt.Sprintf(" ← Key '%c' (shift +%d)", v.alphabet[row], row)))
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println(ui.BrightGreen("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(ui.BrightGreenBold("💡 HOW TO READ THE TABLE:"))
	fmt.Println(ui.Green("   1️⃣  Find your PLAINTEXT letter in the TOP ROW (column header) — shown in GREEN"))
	fmt.Println(ui.Green("   2️⃣  Find your KEY letter in the LEFT COLUMN (row header) — shown in MAGENTA"))
	fmt.Println(ui.Green("   3️⃣  Follow the column DOWN and the row ACROSS to find the intersection"))
	fmt.Println(ui.Green("   4️⃣  The letter at the intersection is your CIPHERTEXT letter!"))
	fmt.Println()

	// Show live examples if we have demo data
	if len(demoPositions) > 0 {
		fmt.Println(ui.BrightYellowBold("🎯 LIVE EXAMPLES FROM YOUR TEXT:"))
		fmt.Println()

		for idx, demo := range demoPositions {
			cipherChar := rune(v.alphabet[demo.cipherPos])

			var color func(string) string
			var symbol string
			if idx == 0 {
				color = ui.BrightRedBold
				symbol = "●"
			} else if idx == 1 {
				color = ui.BrightYellowBold
				symbol = "●"
			} else {
				color = ui.BrightCyanBold
				symbol = "●"
			}

			fmt.Println(color(fmt.Sprintf("   Example %d %s:", idx+1, symbol)))
			if v.mode == "encrypt" {
				fmt.Println(ui.BrightGreen(fmt.Sprintf("      Plaintext:  '%c' (column %d)", demo.plainChar, demo.plainPos)))
				fmt.Println(ui.BrightMagenta(fmt.Sprintf("      Key:        '%c' (row %d, shift +%d)", demo.keyChar, demo.keyPos, demo.keyPos)))
				fmt.Println(color(fmt.Sprintf("      ➜ Ciphertext: '%c' (position %d) %s", cipherChar, demo.cipherPos, symbol)))
				fmt.Println(ui.Yellow(fmt.Sprintf("      Formula: (%d + %d) mod 26 = %d", demo.plainPos, demo.keyPos, demo.cipherPos)))
			} else {
				fmt.Println(ui.BrightGreen(fmt.Sprintf("      Ciphertext: '%c' (column %d)", demo.plainChar, demo.plainPos)))
				fmt.Println(ui.BrightMagenta(fmt.Sprintf("      Key:        '%c' (row %d, shift -%d)", demo.keyChar, demo.keyPos, demo.keyPos)))
				fmt.Println(color(fmt.Sprintf("      ➜ Plaintext:  '%c' (position %d) %s", cipherChar, demo.cipherPos, symbol)))
				fmt.Println(ui.Yellow(fmt.Sprintf("      Formula: (%d - %d + 26) mod 26 = %d", demo.plainPos, demo.keyPos, demo.cipherPos)))
			}
			fmt.Println()
		}
		fmt.Println(ui.BrightCyan("   💫 Look for the colored markers (●) in the table above to see these encryptions!"))
	}

	fmt.Println()
	fmt.Println(ui.BrightGreen("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()
}

// Helper function to check if a cell is in the same row or column as demo
func findDemoInRow(row, col int, demos []struct {
	plainChar rune
	keyChar   rune
	plainPos  int
	keyPos    int
	cipherPos int
}) *struct {
	plainChar rune
	keyChar   rune
	plainPos  int
	keyPos    int
	cipherPos int
} {
	for _, demo := range demos {
		if demo.keyPos == row || demo.plainPos == col {
			return &demo
		}
	}
	return nil
}

func (v *Vigenere) printInteractiveLookup() {
	ui.PrintSectionHeader("INTERACTIVE TABLE LOOKUP GUIDE")
	fmt.Println(ui.BrightCyanBold("🎓 Let's practice using the Vigenère table step-by-step!"))
	fmt.Println()

	cleanKey := v.cleanText(v.key)
	cleanPlain := v.cleanText(v.plaintext)

	if len(cleanKey) == 0 || len(cleanPlain) == 0 {
		return
	}

	// Take first character as example
	plainChar := rune(cleanPlain[0])
	keyChar := rune(cleanKey[0])
	plainPos := strings.IndexRune(v.alphabet, plainChar)
	keyPos := strings.IndexRune(v.alphabet, keyChar)

	if plainPos < 0 || keyPos < 0 {
		return
	}

	var cipherPos int
	if v.mode == "encrypt" {
		cipherPos = (plainPos + keyPos) % len(v.alphabet)
	} else {
		cipherPos = (plainPos - keyPos + len(v.alphabet)) % len(v.alphabet)
	}
	cipherChar := rune(v.alphabet[cipherPos])

	fmt.Println(ui.BrightYellowBold(fmt.Sprintf("📝 Example: Encrypting '%c' with key '%c'", plainChar, keyChar)))
	fmt.Println()

	// Step 1: Find column
	fmt.Println(ui.BrightGreenBold("STEP 1: Find the PLAINTEXT letter in the top row (columns)"))
	fmt.Println()
	fmt.Print(ui.Cyan("   Columns: "))
	for i := 0; i < len(v.alphabet); i++ {
		if i == plainPos {
			fmt.Print(ui.BrightGreenBold(fmt.Sprintf("[%c]", v.alphabet[i])) + " ")
		} else if i == plainPos-1 || i == plainPos+1 {
			fmt.Print(ui.Green(string(v.alphabet[i])) + " ")
		} else {
			fmt.Print(ui.Cyan(string(v.alphabet[i])) + " ")
		}
	}
	fmt.Println()
	fmt.Print(ui.Green("            "))
	for i := 0; i < plainPos; i++ {
		fmt.Print("  ")
	}
	fmt.Println(ui.BrightGreen("▲"))
	fmt.Println(ui.BrightGreenBold(fmt.Sprintf("   Found '%c' at column %d!", plainChar, plainPos)))
	fmt.Println()

	// Step 2: Find row
	fmt.Println(ui.BrightMagentaBold("STEP 2: Find the KEY letter in the left column (rows)"))
	fmt.Println()

	startRow := max(0, keyPos-2)
	endRow := min(len(v.alphabet), keyPos+3)

	for i := startRow; i < endRow; i++ {
		if i == keyPos {
			fmt.Print(ui.BrightMagentaBold(fmt.Sprintf("   ► %c ◄", v.alphabet[i])))
			fmt.Println(ui.BrightMagenta(fmt.Sprintf(" ← KEY letter (row %d)", i)))
		} else {
			fmt.Println(ui.Magenta(fmt.Sprintf("     %c", v.alphabet[i])))
		}
	}
	fmt.Println()
	fmt.Println(ui.BrightMagentaBold(fmt.Sprintf("   Found '%c' at row %d!", keyChar, keyPos)))
	fmt.Println()

	// Step 3: Find intersection
	fmt.Println(ui.BrightCyanBold("STEP 3: Find the INTERSECTION of column and row"))
	fmt.Println()
	fmt.Println(ui.Cyan("   Visual representation:"))
	fmt.Println()

	// ASCII art table lookup
	fmt.Print(ui.Yellow("        "))
	for i := max(0, plainPos-3); i < min(len(v.alphabet), plainPos+4); i++ {
		if i == plainPos {
			fmt.Print(ui.BrightGreenBold(fmt.Sprintf("%c  ", v.alphabet[i])))
		} else {
			fmt.Print(ui.Yellow(fmt.Sprintf("%c  ", v.alphabet[i])))
		}
	}
	fmt.Println(ui.Yellow(" ← Columns (Plaintext)"))

	for row := max(0, keyPos-2); row < min(len(v.alphabet), keyPos+3); row++ {
		if row == keyPos {
			fmt.Print(ui.BrightMagentaBold(fmt.Sprintf("   %c →  ", v.alphabet[row])))
		} else {
			fmt.Print(ui.Magenta(fmt.Sprintf("   %c    ", v.alphabet[row])))
		}

		for col := max(0, plainPos-3); col < min(len(v.alphabet), plainPos+4); col++ {
			shiftedIdx := (col + row) % len(v.alphabet)
			char := string(v.alphabet[shiftedIdx])

			if row == keyPos && col == plainPos {
				// The intersection!
				fmt.Print(ui.BrightRedBold(fmt.Sprintf("[%s]", char)))
			} else if row == keyPos {
				fmt.Print(ui.BrightMagenta(fmt.Sprintf(" %s ", char)))
			} else if col == plainPos {
				fmt.Print(ui.BrightGreen(fmt.Sprintf(" %s ", char)))
			} else {
				fmt.Print(ui.Cyan(fmt.Sprintf(" %s ", char)))
			}
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println(ui.BrightRedBold(fmt.Sprintf("   ✨ RESULT: The intersection is '%c'!", cipherChar)))
	fmt.Println()

	// Show the calculation
	fmt.Println(ui.BrightYellowBold("STEP 4: Verify with mathematics"))
	fmt.Println()
	if v.mode == "encrypt" {
		fmt.Println(ui.Yellow(fmt.Sprintf("   Plaintext '%c' = position %d", plainChar, plainPos)))
		fmt.Println(ui.Yellow(fmt.Sprintf("   Key '%c'       = position %d (shift)", keyChar, keyPos)))
		fmt.Println(ui.Yellow("   ───────────────────────────────"))
		fmt.Println(ui.BrightYellow(fmt.Sprintf("   (%d + %d) mod 26 = %d", plainPos, keyPos, cipherPos)))
		fmt.Println(ui.BrightCyan(fmt.Sprintf("   Position %d = letter '%c'", cipherPos, cipherChar)))
	} else {
		fmt.Println(ui.Yellow(fmt.Sprintf("   Ciphertext '%c' = position %d", plainChar, plainPos)))
		fmt.Println(ui.Yellow(fmt.Sprintf("   Key '%c'        = position %d (shift)", keyChar, keyPos)))
		fmt.Println(ui.Yellow("   ───────────────────────────────"))
		fmt.Println(ui.BrightYellow(fmt.Sprintf("   (%d - %d + 26) mod 26 = %d", plainPos, keyPos, cipherPos)))
		fmt.Println(ui.BrightCyan(fmt.Sprintf("   Position %d = letter '%c'", cipherPos, cipherChar)))
	}

	fmt.Println()
	fmt.Println(ui.BrightGreen("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(ui.BrightGreenBold("🎯 SUMMARY:"))
	if v.mode == "encrypt" {
		fmt.Println(ui.Green(fmt.Sprintf("   To encrypt '%c' with key '%c':", plainChar, keyChar)))
		fmt.Println(ui.Green(fmt.Sprintf("   1. Find column '%c' (plaintext)", plainChar)))
		fmt.Println(ui.Green(fmt.Sprintf("   2. Find row '%c' (key)", keyChar)))
		fmt.Println(ui.Green(fmt.Sprintf("   3. Read intersection: '%c' (ciphertext)", cipherChar)))
	} else {
		fmt.Println(ui.Green(fmt.Sprintf("   To decrypt '%c' with key '%c':", plainChar, keyChar)))
		fmt.Println(ui.Green(fmt.Sprintf("   1. Find column '%c' (ciphertext)", plainChar)))
		fmt.Println(ui.Green(fmt.Sprintf("   2. Find row '%c' (key)", keyChar)))
		fmt.Println(ui.Green(fmt.Sprintf("   3. Read intersection: '%c' (plaintext)", cipherChar)))
	}
	fmt.Println(ui.BrightGreen("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (v *Vigenere) cleanText(text string) string {
	// Remove non-alphabetic characters and convert to uppercase
	var result strings.Builder
	for _, ch := range text {
		if unicode.IsLetter(ch) {
			result.WriteRune(unicode.ToUpper(ch))
		}
	}
	return result.String()
}

func (v *Vigenere) processText() {
	cleanKey := v.cleanText(v.key)
	if len(cleanKey) == 0 {
		fmt.Println(ui.Red("Error: Key must contain at least one letter!"))
		v.ciphertext = v.plaintext
		return
	}

	var result strings.Builder
	keyIndex := 0

	for _, ch := range v.plaintext {
		if !unicode.IsLetter(ch) {
			if v.preserveSpaces {
				result.WriteRune(ch)
			}
			continue
		}

		isUpper := unicode.IsUpper(ch)
		ch = unicode.ToUpper(ch)

		// Find position in alphabet
		plainPos := strings.IndexRune(v.alphabet, ch)
		if plainPos == -1 {
			// Non-standard letter, keep as is
			result.WriteRune(ch)
			continue
		}

		// Get key character
		keyChar := rune(cleanKey[keyIndex%len(cleanKey)])
		keyPos := strings.IndexRune(v.alphabet, keyChar)

		// Apply Vigenère formula
		var cipherPos int
		if v.mode == "encrypt" {
			cipherPos = (plainPos + keyPos) % len(v.alphabet)
		} else {
			cipherPos = (plainPos - keyPos + len(v.alphabet)) % len(v.alphabet)
		}

		cipherChar := rune(v.alphabet[cipherPos])

		// Preserve original case
		if !isUpper && !v.caseSensitive {
			cipherChar = unicode.ToLower(cipherChar)
		}

		result.WriteRune(cipherChar)
		keyIndex++
	}

	v.ciphertext = result.String()
	v.processedKey = cleanKey
}

func (v *Vigenere) printStepByStep() {
	ui.PrintSectionHeader("STEP-BY-STEP ENCRYPTION PROCESS")
	fmt.Println(ui.BrightCyanBold("🔍 Detailed character-by-character transformation:"))
	fmt.Println()

	// Key repetition visualization
	cleanKey := v.processedKey
	fmt.Println(ui.BrightMagentaBold("Step 1: Key Repetition Pattern"))
	fmt.Println(ui.Yellow("The key is repeated to match the length of the plaintext:"))
	fmt.Println()

	// Build repeated key
	var repeatedKey strings.Builder
	keyIndex := 0

	inputLabel := "Input"
	outputLabel := "Output"
	if v.mode == "decrypt" {
		inputLabel = "Cipher"
		outputLabel = "Plain"
	}

	fmt.Print(ui.Cyan(fmt.Sprintf("%-12s: ", inputLabel)))
	for _, ch := range v.plaintext {
		if unicode.IsLetter(ch) {
			fmt.Print(ui.BrightYellow(string(unicode.ToUpper(ch))))
			repeatedKey.WriteRune(rune(cleanKey[keyIndex%len(cleanKey)]))
			keyIndex++
		} else {
			fmt.Print(ui.Yellow(string(ch)))
			if v.preserveSpaces {
				repeatedKey.WriteRune(' ')
			}
		}
	}
	fmt.Println()

	fmt.Print(ui.Magenta(fmt.Sprintf("%-12s: ", "Key")))
	for _, ch := range repeatedKey.String() {
		if ch == ' ' {
			fmt.Print(ui.Magenta(string(ch)))
		} else {
			fmt.Print(ui.BrightMagenta(string(ch)))
		}
	}
	fmt.Println()
	fmt.Println()

	// Character-by-character transformation
	fmt.Println(ui.BrightMagentaBold("Step 2: Character Transformation"))
	fmt.Println(ui.Yellow("Each letter is shifted by the corresponding key letter:"))
	fmt.Println()

	t := ui.GetTableWriter()
	t.AppendHeader(table.Row{
		ui.CyanBold("Pos"),
		ui.MagentaBold(inputLabel),
		ui.GreenBold("Key"),
		ui.YellowBold("Shift"),
		ui.BlueBold("Calculation"),
		ui.BrightGreenBold(outputLabel),
	})

	keyIndex = 0
	outputIndex := 0

	for i, ch := range v.plaintext {
		if !unicode.IsLetter(ch) {
			if v.preserveSpaces {
				t.AppendRow(table.Row{
					ui.Cyan(fmt.Sprintf("%d", i+1)),
					ui.Yellow(string(ch)),
					ui.Magenta("—"),
					ui.Yellow("—"),
					ui.Blue("(preserved)"),
					ui.Green(string(ch)),
				})
			}
			continue
		}

		isUpper := unicode.IsUpper(ch)
		ch = unicode.ToUpper(ch)

		plainPos := strings.IndexRune(v.alphabet, ch)
		if plainPos == -1 {
			continue
		}

		keyChar := rune(cleanKey[keyIndex%len(cleanKey)])
		keyPos := strings.IndexRune(v.alphabet, keyChar)

		var cipherPos int
		var formula string
		if v.mode == "encrypt" {
			cipherPos = (plainPos + keyPos) % len(v.alphabet)
			formula = fmt.Sprintf("(%d + %d) mod 26 = %d", plainPos, keyPos, cipherPos)
		} else {
			cipherPos = (plainPos - keyPos + len(v.alphabet)) % len(v.alphabet)
			formula = fmt.Sprintf("(%d - %d) mod 26 = %d", plainPos, keyPos, cipherPos)
		}

		cipherChar := rune(v.alphabet[cipherPos])
		if !isUpper {
			cipherChar = unicode.ToLower(cipherChar)
		}

		operation := "+"
		if v.mode == "decrypt" {
			operation = "−"
		}

		t.AppendRow(table.Row{
			ui.Cyan(fmt.Sprintf("%d", i+1)),
			ui.BrightYellow(string(ch)) + ui.Yellow(fmt.Sprintf(" (%d)", plainPos)),
			ui.BrightMagenta(string(keyChar)) + ui.Magenta(fmt.Sprintf(" (%d)", keyPos)),
			ui.Yellow(fmt.Sprintf("%s%d", operation, keyPos)),
			ui.Blue(formula),
			ui.BrightGreen(string(cipherChar)) + ui.Green(fmt.Sprintf(" (%d)", cipherPos)),
		})

		keyIndex++
		outputIndex++

		// Show first 15 rows for readability
		if outputIndex >= 15 && len(v.plaintext) > 20 {
			t.AppendSeparator()
			t.AppendRow(table.Row{
				ui.Yellow("..."),
				ui.Yellow("..."),
				ui.Yellow("..."),
				ui.Yellow("..."),
				ui.Yellow(fmt.Sprintf("(showing first 15 of %d letters)", keyIndex)),
				ui.Yellow("..."),
			})
			break
		}
	}

	t.Render()
	fmt.Println()

	// Visual transformation flow
	fmt.Println(ui.BrightMagentaBold("Step 3: Complete Transformation"))
	fmt.Println()

	maxLineLength := 60
	inputText := v.plaintext
	outputText := v.ciphertext

	for i := 0; i < len(inputText); i += maxLineLength {
		end := i + maxLineLength
		if end > len(inputText) {
			end = len(inputText)
		}

		segment := inputText[i:end]
		outputSegment := outputText[i:end]

		fmt.Print(ui.Cyan(fmt.Sprintf("%s: ", inputLabel)))
		for _, ch := range segment {
			if unicode.IsLetter(ch) {
				fmt.Print(ui.BrightYellow(string(ch)))
			} else {
				fmt.Print(ui.Yellow(string(ch)))
			}
		}
		fmt.Println()

		fmt.Print(ui.Yellow("          "))
		for range segment {
			fmt.Print(ui.Green("↓"))
		}
		fmt.Println()

		fmt.Print(ui.Green(fmt.Sprintf("%s: ", outputLabel)))
		for _, ch := range outputSegment {
			if unicode.IsLetter(ch) {
				fmt.Print(ui.BrightGreen(string(ch)))
			} else {
				fmt.Print(ui.Green(string(ch)))
			}
		}
		fmt.Println()
		fmt.Println()
	}
}

func (v *Vigenere) printResults() {
	ui.PrintSectionHeader("FINAL RESULTS")

	inputLabel := "Original Plaintext"
	outputLabel := "Encrypted Ciphertext"
	if v.mode == "decrypt" {
		inputLabel = "Original Ciphertext"
		outputLabel = "Decrypted Plaintext"
	}

	t := ui.GetTableWriter()
	t.AppendHeader(table.Row{
		ui.MagentaBold("Item"),
		ui.CyanBold("Value"),
	})

	t.AppendRow(table.Row{
		ui.Magenta(inputLabel),
		ui.BrightYellow(fmt.Sprintf("\"%s\"", v.plaintext)),
	})

	t.AppendRow(table.Row{
		ui.Magenta("Key Used"),
		ui.BrightGreen(fmt.Sprintf("\"%s\"", v.key)),
	})

	t.AppendSeparator()

	t.AppendRow(table.Row{
		ui.BrightMagenta(outputLabel),
		ui.BrightCyan(fmt.Sprintf("\"%s\"", v.ciphertext)),
	})

	t.Render()
	fmt.Println()

	// Verification section
	fmt.Println(ui.BrightGreenBold("✅ Verification:"))
	fmt.Println()

	if v.mode == "encrypt" {
		fmt.Println(ui.Green("To decrypt this message, use:"))
		fmt.Println(ui.BrightCyan(fmt.Sprintf("   • Ciphertext: %s", v.ciphertext)))
		fmt.Println(ui.BrightGreen(fmt.Sprintf("   • Key: %s", v.key)))
		fmt.Println(ui.BrightYellow("   • Mode: Decrypt"))
	} else {
		fmt.Println(ui.Green("To encrypt this message back, use:"))
		fmt.Println(ui.BrightCyan(fmt.Sprintf("   • Plaintext: %s", v.ciphertext)))
		fmt.Println(ui.BrightGreen(fmt.Sprintf("   • Key: %s", v.key)))
		fmt.Println(ui.BrightYellow("   • Mode: Encrypt"))
	}

	fmt.Println()
}

func (v *Vigenere) printCryptanalysisInfo() {
	ui.PrintSectionHeader("CRYPTANALYSIS INSIGHTS")
	fmt.Println(ui.BrightCyanBold("🔬 Security Analysis and Historical Context:"))
	fmt.Println()

	fmt.Println(ui.BrightYellowBold("📜 Historical Significance:"))
	fmt.Println(ui.Yellow("   • Invented in the 1550s, popularized by Blaise de Vigenère"))
	fmt.Println(ui.Yellow("   • Nicknamed 'le chiffre indéchiffrable' (the indecipherable cipher)"))
	fmt.Println(ui.Yellow("   • Remained unbroken for over 300 years!"))
	fmt.Println(ui.Yellow("   • Finally cracked by Charles Babbage and Friedrich Kasiski in the 19th century"))
	fmt.Println()

	fmt.Println(ui.BrightMagentaBold("💪 Strengths:"))
	fmt.Println(ui.Magenta("   • Resists simple frequency analysis (unlike Caesar cipher)"))
	fmt.Println(ui.Magenta("   • Each letter can be encrypted differently each time it appears"))
	fmt.Println(ui.Magenta("   • Simple to implement with pen and paper"))
	fmt.Println(ui.Magenta("   • No special equipment needed"))
	fmt.Println()

	fmt.Println(ui.BrightRedBold("⚠️  Weaknesses:"))
	fmt.Println(ui.Red("   • Vulnerable to Kasiski examination (finds key length)"))
	fmt.Println(ui.Red("   • Vulnerable to Friedman test (statistical key length detection)"))
	fmt.Println(ui.Red("   • Short keys are particularly weak"))
	fmt.Println(ui.Red("   • Repeated key allows pattern analysis"))
	fmt.Println(ui.Red("   • Not secure against modern computational attacks"))
	fmt.Println()

	// Key strength analysis
	keyLength := len(v.processedKey)
	fmt.Println(ui.BrightCyanBold("🔑 Your Key Analysis:"))

	var strength string
	var strengthColor func(string) string
	var advice string

	if keyLength < 4 {
		strength = "VERY WEAK"
		strengthColor = ui.BrightRedBold
		advice = "Use a longer key (at least 12 characters) for better security."
	} else if keyLength < 8 {
		strength = "WEAK"
		strengthColor = ui.RedBold
		advice = "Consider using a longer key (12+ characters) for improved security."
	} else if keyLength < 12 {
		strength = "MODERATE"
		strengthColor = ui.YellowBold
		advice = "Good length! For maximum security, use 15+ character random keys."
	} else if keyLength < 20 {
		strength = "STRONG"
		strengthColor = ui.GreenBold
		advice = "Excellent key length! Random keys of this length are quite secure."
	} else {
		strength = "VERY STRONG"
		strengthColor = ui.BrightGreenBold
		advice = "Outstanding! When key length equals plaintext, it becomes a One-Time Pad."
	}

	fmt.Println(ui.Cyan(fmt.Sprintf("   • Key Length: %d characters", keyLength)))
	fmt.Println(ui.Cyan(fmt.Sprintf("   • Text Length: %d characters", len(v.cleanText(v.plaintext)))))
	fmt.Printf("   • Strength: %s\n", strengthColor(strength))
	fmt.Println(ui.Yellow(fmt.Sprintf("   • Advice: %s", advice)))
	fmt.Println()

	// Attack methods
	fmt.Println(ui.BrightBlueBold("🎯 Known Attack Methods:"))
	fmt.Println()

	attackTable := ui.GetTableWriter()
	attackTable.AppendHeader(table.Row{
		ui.CyanBold("Attack Method"),
		ui.YellowBold("Description"),
		ui.MagentaBold("Effectiveness"),
	})

	attackTable.AppendRow(table.Row{
		ui.Cyan("Kasiski Examination"),
		ui.Yellow("Finds repeated sequences to deduce key length"),
		ui.Red("High for short keys"),
	})

	attackTable.AppendRow(table.Row{
		ui.Cyan("Friedman Test"),
		ui.Yellow("Uses Index of Coincidence to estimate key length"),
		ui.Red("High for any key"),
	})

	attackTable.AppendRow(table.Row{
		ui.Cyan("Frequency Analysis"),
		ui.Yellow("After finding key length, analyze each Caesar shift"),
		ui.Red("High once key length known"),
	})

	attackTable.AppendRow(table.Row{
		ui.Cyan("Brute Force"),
		ui.Yellow("Try all possible keys"),
		ui.Green("Low for long keys"),
	})

	attackTable.AppendRow(table.Row{
		ui.Cyan("Known Plaintext"),
		ui.Yellow("If part of plaintext is known, key can be derived"),
		ui.BrightRed("Very High"),
	})

	attackTable.Render()
	fmt.Println()

	fmt.Println(ui.BrightGreenBold("💡 Modern Recommendations:"))
	fmt.Println(ui.Green("   • For actual security, use modern algorithms (AES, ChaCha20)"))
	fmt.Println(ui.Green("   • Vigenère is great for learning cryptography concepts"))
	fmt.Println(ui.Green("   • Never use Vigenère for real-world sensitive data"))
	fmt.Println(ui.Green("   • If you must use it, use truly random keys as long as the message"))
	fmt.Println()

	fmt.Println(ui.BrightCyanBold("🎓 Learning Exercise:"))
	fmt.Println(ui.Cyan("   Try encrypting the same message with different key lengths and"))
	fmt.Println(ui.Cyan("   observe how the ciphertext changes. Notice patterns with short keys!"))
	fmt.Println()
}
