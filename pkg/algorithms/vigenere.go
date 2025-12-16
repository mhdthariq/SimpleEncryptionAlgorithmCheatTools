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
	fmt.Println(ui.BrightCyanBold("📊 The Vigenère square shows all possible Caesar shifts:"))
	fmt.Println(ui.BrightYellow("   • Each row is shifted one position from the previous"))
	fmt.Println(ui.BrightYellow("   • The key letter determines which row to use"))
	fmt.Println(ui.BrightYellow("   • The plaintext letter determines which column to use"))
	fmt.Println()

	// Print a compact version (first 13 rows for display purposes)
	displayRows := 13
	if displayRows > len(v.alphabet) {
		displayRows = len(v.alphabet)
	}

	// Header row
	fmt.Print(ui.MagentaBold("   │ "))
	for i := 0; i < len(v.alphabet); i++ {
		fmt.Print(ui.CyanBold(string(v.alphabet[i])) + " ")
	}
	fmt.Println()

	// Separator
	fmt.Print(ui.MagentaBold("───┼"))
	for i := 0; i < len(v.alphabet); i++ {
		fmt.Print("──")
	}
	fmt.Println()

	// Table rows
	for row := 0; row < displayRows; row++ {
		fmt.Print(ui.MagentaBold(fmt.Sprintf(" %c │ ", v.alphabet[row])))
		for col := 0; col < len(v.alphabet); col++ {
			shiftedIdx := (col + row) % len(v.alphabet)
			char := string(v.alphabet[shiftedIdx])

			// Highlight diagonal
			if row == col {
				fmt.Print(ui.BrightYellow(char) + " ")
			} else if row < 3 || col < 3 {
				fmt.Print(ui.BrightCyan(char) + " ")
			} else {
				fmt.Print(ui.Cyan(char) + " ")
			}
		}
		fmt.Println()
	}

	if displayRows < len(v.alphabet) {
		fmt.Println(ui.Yellow(fmt.Sprintf("   ... (showing %d of %d rows) ...", displayRows, len(v.alphabet))))
	}

	fmt.Println()
	fmt.Println(ui.BrightGreen("💡 How to use:"))
	fmt.Println(ui.Green("   1. Find the plaintext letter in the top row (column header)"))
	fmt.Println(ui.Green("   2. Find the key letter in the left column (row header)"))
	fmt.Println(ui.Green("   3. The intersection is your ciphertext letter!"))
	fmt.Println()
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
