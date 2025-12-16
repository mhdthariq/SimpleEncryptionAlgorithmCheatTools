package algorithms

import (
	"bufio"
	"fmt"
	"strings"

	"example.com/cheat-encryption-algorithm/pkg/ui" // Assuming module name is cheat-encryption-algorithm

	"github.com/jedib0t/go-pretty/v6/table"
)

// RC4 implementation
type RC4 struct {
	plaintext    string
	key          string
	stateSize    int
	S            []int
	keystream    []int
	ciphertext   []int
	plaintextASC []int
	keyASC       []int
}

func NewRC4() *RC4 {
	return &RC4{}
}

func (r *RC4) GetName() string {
	return "RC4"
}

func (r *RC4) Run(reader *bufio.Reader, plaintext string, key string) {
	// RC4 Specific Input
	fmt.Print("Please input State Array Size (e.g., 4, 256): ")
	stateSizeStr, _ := reader.ReadString('\n')
	stateSizeStr = strings.TrimSpace(stateSizeStr)

	var stateSize int
	_, err := fmt.Sscanf(stateSizeStr, "%d", &stateSize)

	if err != nil || stateSize <= 0 {
		fmt.Println("Error: State size must be a positive number!")
		return
	}

	// Initialize
	r.plaintext = plaintext
	r.key = key
	r.stateSize = stateSize
	r.S = make([]int, stateSize)
	r.keystream = make([]int, 0)

	r.plaintextASC = make([]int, len(plaintext))
	for i, c := range plaintext {
		r.plaintextASC[i] = int(c)
	}

	r.keyASC = make([]int, len(key))
	for i, c := range key {
		r.keyASC[i] = int(c)
	}

	fmt.Println()
	fmt.Println("Processing...")
	fmt.Println()

	ui.PrintSectionHeader("RC4 ENCRYPTION ALGORITHM - STEP BY STEP")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightCyanBold("🔐 This demonstration will walk you through each step of the RC4 encryption process."))
	fmt.Println(ui.BrightYellow("📚 Follow the colors to understand each transformation!"))
	ui.PrintStarLine(40)
	fmt.Println()

	r.printSetupTable()
	r.runKSA()
	r.runPRGA()
	r.printXORAnalysis()

	ui.PrintSectionHeader("COMPLETE RESULTS")
	r.printResults()
}

func (r *RC4) printSetupTable() {
	ui.PrintSectionHeader("EXAMPLE SETUP")
	fmt.Println(ui.BrightCyanBold("✨ Initial parameters for the RC4 encryption:"))
	fmt.Println()

	t := ui.GetTableWriter()

	t.AppendHeader(table.Row{
		ui.MagentaBold("Parameter"),
		ui.CyanBold("Value"),
		ui.YellowBold("Description"),
	})

	asciiStr := "["
	for i, v := range r.plaintextASC {
		if i > 0 {
			asciiStr += ", "
		}
		asciiStr += fmt.Sprintf("%d", v)
	}
	asciiStr += "]"

	keyASCIIStr := "["
	for i, v := range r.keyASC {
		if i > 0 {
			keyASCIIStr += ", "
		}
		keyASCIIStr += fmt.Sprintf("%d", v)
	}
	keyASCIIStr += "]"

	t.AppendRows([]table.Row{
		{ui.BrightMagenta("Plaintext"), ui.BrightYellow("\"" + r.plaintext + "\""), ui.Cyan("Input message")},
		{ui.BrightMagenta("ASCII Values"), ui.BrightGreen(asciiStr), ui.Cyan("Plaintext in ASCII")},
		{ui.BrightMagenta("Key"), ui.BrightYellow("\"" + r.key + "\""), ui.Cyan("Encryption key")},
		{ui.BrightMagenta("Key ASCII"), ui.BrightGreen(keyASCIIStr), ui.Cyan("Key in ASCII")},
		{ui.BrightMagenta("State Array Size"), ui.BrightRed(fmt.Sprintf("%d", r.stateSize)), ui.Cyan("Array size")},
	})

	t.Render()
	fmt.Println()
}

func (r *RC4) runKSA() {
	ui.PrintSectionHeader("KEY-SCHEDULING ALGORITHM (KSA)")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightYellowBold("🔑 Phase 1: Scrambling the state array using the key"))
	ui.PrintSeparator()
	ui.PrintStarLine(40)
	fmt.Println()

	for i := 0; i < r.stateSize; i++ {
		r.S[i] = i
	}

	fmt.Println(ui.BrightMagentaBold("🎯 Initial Setup:"))
	fmt.Printf("%s = %s\n", ui.BrightYellowBold("S"), ui.BrightGreen(fmt.Sprintf("%v", r.S)))
	fmt.Printf("%s = %s\n", ui.BrightYellowBold("j"), ui.BrightGreen("0"))
	fmt.Println()

	j := 0
	keyLen := len(r.keyASC)

	for i := 0; i < r.stateSize; i++ {
		ui.PrintColorfulLine(50)
		fmt.Printf("%s %s (i = %s)\n", ui.BrightMagentaBold("🔄 KSA Iteration"), ui.BrightYellow(fmt.Sprintf("%d", i+1)), ui.BrightCyan(fmt.Sprintf("%d", i)))
		ui.PrintColorfulLine(50)

		fmt.Println(ui.BrightCyanBold("📊 Current State:"))
		fmt.Printf("  %s %s = %s  │  %s = %s\n",
			ui.Magenta("•"),
			ui.YellowBold("i"),
			ui.BrightYellow(fmt.Sprintf("%d", i)),
			ui.YellowBold("j"),
			ui.BrightYellow(fmt.Sprintf("%d", j)))
		fmt.Printf("  %s %s = %s\n\n", ui.Magenta("•"), ui.YellowBold("S"), ui.BrightGreen(fmt.Sprintf("%v", r.S)))

		oldJ := j
		j = (j + r.S[i] + r.keyASC[i%keyLen]) % r.stateSize

		fmt.Println(ui.BrightYellowBold("🧮 Calculation:"))
		fmt.Printf("  %s = (j + S[i] + K[i %% %s]) %% %s\n",
			ui.MagentaBold("j"),
			ui.BrightCyan(fmt.Sprintf("%d", keyLen)),
			ui.BrightCyan(fmt.Sprintf("%d", r.stateSize)))
		fmt.Printf("  %s = (%s + S[%s] + K[%s]) %% %s\n",
			ui.MagentaBold("j"),
			ui.BrightYellow(fmt.Sprintf("%d", oldJ)),
			ui.BrightCyan(fmt.Sprintf("%d", i)),
			ui.BrightCyan(fmt.Sprintf("%d", i%keyLen)),
			ui.BrightCyan(fmt.Sprintf("%d", r.stateSize)))
		fmt.Printf("  %s = (%s + %s + %s) %% %s\n",
			ui.MagentaBold("j"),
			ui.BrightYellow(fmt.Sprintf("%d", oldJ)),
			ui.Green(fmt.Sprintf("%d", r.S[i])),
			ui.Green(fmt.Sprintf("%d", r.keyASC[i%keyLen])),
			ui.BrightCyan(fmt.Sprintf("%d", r.stateSize)))
		fmt.Printf("  %s = %s %% %s = %s ✨\n\n",
			ui.MagentaBold("j"),
			ui.Yellow(fmt.Sprintf("%d", oldJ+r.S[i]+r.keyASC[i%keyLen])),
			ui.BrightCyan(fmt.Sprintf("%d", r.stateSize)),
			ui.BrightMagentaBold(fmt.Sprintf("%d", j)))

		oldS := make([]int, len(r.S))
		copy(oldS, r.S)

		r.S[i], r.S[j] = r.S[j], r.S[i]

		fmt.Printf("%s S[%s] ⇄ S[%s]:\n",
			ui.GreenBold("🔀 Swap"),
			ui.BrightYellow(fmt.Sprintf("%d", i)),
			ui.BrightYellow(fmt.Sprintf("%d", j)))
		fmt.Printf("  %s: S = %s\n", ui.RedBold("❌ Before"), ui.Red(fmt.Sprintf("%v", oldS)))
		fmt.Printf("  %s:  S = %s\n", ui.GreenBold("✅ After"), ui.BrightGreen(fmt.Sprintf("%v", r.S)))
		fmt.Println()
	}

	ui.PrintSeparator()
	resultTable := ui.GetTableWriter()
	resultTable.SetTitle(ui.BrightGreenBold("✓ KSA RESULT"))
	resultTable.AppendRow(table.Row{
		ui.BrightCyanBold("Final scrambled state array S"),
		ui.BrightGreen(fmt.Sprintf("%v", r.S)),
	})
	resultTable.Render()
	ui.PrintSeparator()
	fmt.Println()
}

func (r *RC4) runPRGA() {
	ui.PrintSectionHeader("PSEUDO-RANDOM GENERATION ALGORITHM (PRGA)")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightYellowBold("🔐 Phase 2: Generating keystream and encrypting plaintext"))
	ui.PrintSeparator()
	ui.PrintStarLine(40)
	fmt.Println()

	fmt.Println(ui.BrightMagentaBold("🎯 Initial Setup:"))
	fmt.Printf("  %s %s = %s (from KSA result)\n", ui.Magenta("•"), ui.YellowBold("S"), ui.BrightGreen(fmt.Sprintf("%v", r.S)))
	fmt.Printf("  %s %s = %s\n", ui.Magenta("•"), ui.YellowBold("i"), ui.BrightGreen("0"))
	fmt.Printf("  %s %s = %s\n", ui.Magenta("•"), ui.YellowBold("j"), ui.BrightGreen("0"))
	fmt.Println()

	i := 0
	j := 0
	r.ciphertext = make([]int, len(r.plaintextASC))

	for byteIdx, plaintextByte := range r.plaintextASC {
		ui.PrintColorfulLine(50)
		fmt.Printf("%s %s: Encrypt Byte %s (%s)\n",
			ui.BrightMagentaBold("🔐 Step"),
			ui.BrightCyan(fmt.Sprintf("%d", byteIdx+1)),
			ui.BrightYellowBold(fmt.Sprintf("'%c'", rune(plaintextByte))),
			ui.BrightYellow(fmt.Sprintf("%d", plaintextByte)))
		ui.PrintColorfulLine(50)

		// Step 1: Update Indices
		fmt.Printf("%s PRGA Step %s: %s\n",
			ui.BlueBold("►"),
			ui.BrightCyan(fmt.Sprintf("%d.1", byteIdx+1)),
			ui.BrightCyanBold("Update Indices"))

		fmt.Println(ui.BrightYellowBold("📊 Current State:"))
		fmt.Printf("  %s %s = %s  │  %s = %s\n",
			ui.Cyan("•"),
			ui.YellowBold("i"),
			ui.BrightYellow(fmt.Sprintf("%d", i)),
			ui.YellowBold("j"),
			ui.BrightYellow(fmt.Sprintf("%d", j)))
		fmt.Printf("  %s %s = %s\n", ui.Cyan("•"), ui.YellowBold("S"), ui.Green(fmt.Sprintf("%v", r.S)))
		fmt.Printf("  %s Plaintext byte: %s (%s)\n\n",
			ui.Cyan("•"),
			ui.BrightYellowBold(fmt.Sprintf("%d", plaintextByte)),
			ui.BrightYellowBold(fmt.Sprintf("'%c'", rune(plaintextByte))))

		oldI := i
		oldJ := j
		i = (i + 1) % r.stateSize
		j = (j + r.S[i]) % r.stateSize

		fmt.Println(ui.BrightMagentaBold("🧮 Index Updates:"))
		fmt.Printf("  %s = (i + 1) %% %s = (%s + 1) %% %s = %s ✨\n",
			ui.MagentaBold("i"),
			ui.Cyan(fmt.Sprintf("%d", r.stateSize)),
			ui.Yellow(fmt.Sprintf("%d", oldI)),
			ui.Cyan(fmt.Sprintf("%d", r.stateSize)),
			ui.BrightMagentaBold(fmt.Sprintf("%d", i)))
		fmt.Printf("  %s = (j + S[i]) %% %s = (%s + %s) %% %s = %s ✨\n\n",
			ui.MagentaBold("j"),
			ui.Cyan(fmt.Sprintf("%d", r.stateSize)),
			ui.Yellow(fmt.Sprintf("%d", oldJ)),
			ui.Green(fmt.Sprintf("%d", r.S[i])),
			ui.Cyan(fmt.Sprintf("%d", r.stateSize)),
			ui.BrightMagentaBold(fmt.Sprintf("%d", j)))

		// Step 2: Swap Operation
		fmt.Printf("%s PRGA Step %s: %s\n",
			ui.BlueBold("►"),
			ui.BrightCyan(fmt.Sprintf("%d.2", byteIdx+1)),
			ui.BrightCyanBold("Swap Operation"))

		oldS := make([]int, len(r.S))
		copy(oldS, r.S)

		r.S[i], r.S[j] = r.S[j], r.S[i]

		fmt.Printf("%s S[%s] ⇄ S[%s]:\n",
			ui.GreenBold("🔀 Swap"),
			ui.BrightYellow(fmt.Sprintf("%d", i)),
			ui.BrightYellow(fmt.Sprintf("%d", j)))
		fmt.Printf("  %s: S = %s\n", ui.RedBold("❌ Before"), ui.Red(fmt.Sprintf("%v", oldS)))
		fmt.Printf("  %s:  S = %s\n\n", ui.GreenBold("✅ After"), ui.BrightGreen(fmt.Sprintf("%v", r.S)))

		// Step 3: Generate Keystream
		fmt.Printf("%s PRGA Step %s: %s\n",
			ui.BlueBold("►"),
			ui.BrightCyan(fmt.Sprintf("%d.3", byteIdx+1)),
			ui.BrightCyanBold("Generate Keystream"))

		t := (r.S[i] + r.S[j]) % r.stateSize
		keystreamByte := r.S[t]

		fmt.Println(ui.BrightYellowBold("🧮 Calculation:"))
		fmt.Printf("  %s = (S[%s] + S[%s]) %% %s = (%s + %s) %% %s = %s ✨\n",
			ui.CyanBold("t"),
			ui.BrightYellow(fmt.Sprintf("%d", i)),
			ui.BrightYellow(fmt.Sprintf("%d", j)),
			ui.Cyan(fmt.Sprintf("%d", r.stateSize)),
			ui.Green(fmt.Sprintf("%d", r.S[i])),
			ui.Green(fmt.Sprintf("%d", r.S[j])),
			ui.Cyan(fmt.Sprintf("%d", r.stateSize)),
			ui.BrightCyanBold(fmt.Sprintf("%d", t)))
		fmt.Printf("  %s = S[t] = S[%s] = %s 🔑\n\n",
			ui.GreenBold("Keystream byte"),
			ui.BrightCyan(fmt.Sprintf("%d", t)),
			ui.BrightGreenBold(fmt.Sprintf("%d", keystreamByte)))

		r.keystream = append(r.keystream, keystreamByte)

		// Step 4: Encryption
		fmt.Printf("%s PRGA Step %s: %s\n",
			ui.BlueBold("►"),
			ui.BrightCyan(fmt.Sprintf("%d.4", byteIdx+1)),
			ui.BrightCyanBold("Encryption"))

		ciphertextByte := plaintextByte ^ keystreamByte
		r.ciphertext[byteIdx] = ciphertextByte

		fmt.Println(ui.BrightMagentaBold("⚡ XOR Encryption:"))
		fmt.Printf("  %s %s %s = %s\n",
			ui.BrightYellowBold(fmt.Sprintf("%d", plaintextByte)),
			ui.MagentaBold("⊕ XOR ⊕"),
			ui.BrightGreenBold(fmt.Sprintf("%d", keystreamByte)),
			ui.BrightRedBold(fmt.Sprintf("%d", ciphertextByte)))

		// Specific requested output format
		fmt.Printf("\n%s Ciphertext byte %s = %s 🎯\n\n",
			ui.GreenBold("🎉 Result"),
			ui.BrightRedBold(fmt.Sprintf("%d", ciphertextByte)),
			ui.BrightRedBold(fmt.Sprintf("'%c'", rune(ciphertextByte))))
	}
}

func (r *RC4) printXORAnalysis() {
	ui.PrintSectionHeader("XOR OPERATION ANALYSIS")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightCyanBold("🔬 This section provides a detailed bit-level analysis of the XOR encryption"))
	fmt.Println(ui.BrightCyan("   operation performed during PRGA. Each plaintext byte is XORed with its"))
	fmt.Println(ui.BrightCyan("   corresponding keystream byte to produce the ciphertext."))
	ui.PrintStarLine(40)
	fmt.Println()

	for i := 0; i < len(r.plaintextASC); i++ {
		plainByte := r.plaintextASC[i]
		keyByte := r.keystream[i]
		cipherByte := r.ciphertext[i]

		ui.PrintColorfulLine(50)
		fmt.Printf("%s %s: Plaintext %s (%s) %s Keystream (%s) = Ciphertext (%s)\n",
			ui.BrightMagentaBold("🔍 XOR Analysis"),
			ui.BrightCyan(fmt.Sprintf("%d", i+1)),
			ui.BrightYellowBold(fmt.Sprintf("'%c'", rune(plainByte))),
			ui.BrightYellow(fmt.Sprintf("%d", plainByte)),
			ui.MagentaBold("⊕"),
			ui.BrightGreen(fmt.Sprintf("%d", keyByte)),
			ui.BrightRedBold(fmt.Sprintf("%d", cipherByte)))
		ui.PrintColorfulLine(50)

		plainBin := fmt.Sprintf("%08b", plainByte)
		keyBin := fmt.Sprintf("%08b", keyByte)
		resultBin := fmt.Sprintf("%08b", cipherByte)

		// Binary representation table
		binTable := ui.GetTableWriter()
		binTable.SetTitle(ui.BrightCyanBold("Binary Representation"))
		binTable.AppendHeader(table.Row{
			ui.MagentaBold("Type"),
			ui.YellowBold("Decimal"),
			ui.CyanBold("Character"),
			ui.GreenBold("Binary (8-bit)"),
		})
		binTable.AppendRows([]table.Row{
			{ui.BrightYellow("Plaintext"), ui.BrightYellow(fmt.Sprintf("%d", plainByte)), ui.BrightYellow(fmt.Sprintf("'%c'", rune(plainByte))), ui.Cyan(plainBin)},
			{ui.BrightGreen("Keystream"), ui.BrightGreen(fmt.Sprintf("%d", keyByte)), ui.Green("-"), ui.Cyan(keyBin)},
			{ui.Magenta(strings.Repeat("─", 10)), ui.Magenta(strings.Repeat("─", 7)), ui.Magenta(strings.Repeat("─", 11)), ui.Magenta(strings.Repeat("─", 8))},
			{ui.BrightRed("Ciphertext"), ui.BrightRed(fmt.Sprintf("%d", cipherByte)), ui.BrightRed(fmt.Sprintf("'%c'", rune(cipherByte))), ui.Cyan(resultBin)},
		})
		binTable.Render()
		fmt.Println()

		// Bit-by-bit XOR table
		xorTable := ui.GetTableWriter()
		xorTable.SetTitle(ui.BrightMagentaBold("⚡ Bit-by-Bit XOR Operation"))

		headerRow := table.Row{ui.CyanBold("Source")}
		for bit := 7; bit >= 0; bit-- {
			headerRow = append(headerRow, ui.MagentaBold(fmt.Sprintf("Bit %d", bit)))
		}
		xorTable.AppendHeader(headerRow)

		plainRow := table.Row{ui.BrightYellow(fmt.Sprintf("Plain (%d)", plainByte))}
		for bit := 7; bit >= 0; bit-- {
			plainRow = append(plainRow, ui.BrightCyan(string(plainBin[7-bit])))
		}
		xorTable.AppendRow(plainRow)

		keyRow := table.Row{ui.BrightGreen(fmt.Sprintf("Key (%d)", keyByte))}
		for bit := 7; bit >= 0; bit-- {
			keyRow = append(keyRow, ui.BrightCyan(string(keyBin[7-bit])))
		}
		xorTable.AppendRow(keyRow)

		xorTable.AppendSeparator()

		xorRow := table.Row{ui.BrightRed(fmt.Sprintf("Result (%d)", cipherByte))}
		for bit := 7; bit >= 0; bit-- {
			xorRow = append(xorRow, ui.BrightYellow(string(resultBin[7-bit])))
		}
		xorTable.AppendRow(xorRow)

		xorTable.Render()
		fmt.Println()

		// XOR Truth Table Reminder
		fmt.Println(ui.BrightCyanBold("📘 XOR Truth Table Reference:"))
		fmt.Printf("  %s  │  %s\n", ui.BrightYellow("0 ⊕ 0 = 0"), ui.BrightYellow("0 ⊕ 1 = 1"))
		fmt.Printf("  %s  │  %s\n", ui.BrightYellow("1 ⊕ 0 = 1"), ui.BrightYellow("1 ⊕ 1 = 0"))
		fmt.Println()
	}

	// Summary of XOR properties
	ui.PrintSectionHeader("XOR PROPERTIES SUMMARY")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightMagentaBold("🔬 Understanding the mathematical properties that make XOR perfect for encryption:"))
	ui.PrintStarLine(40)
	fmt.Println()

	propertiesTable := ui.GetTableWriter()
	propertiesTable.SetTitle(ui.BrightGreenBold("✨ Key Properties of XOR in RC4"))
	propertiesTable.AppendHeader(table.Row{
		ui.MagentaBold("Property"),
		ui.CyanBold("Description"),
		ui.YellowBold("Example"),
	})
	propertiesTable.AppendRows([]table.Row{
		{ui.BrightGreen("Reversibility"), ui.Cyan("A ⊕ B ⊕ B = A"), ui.BrightYellow("Cipher ⊕ Key = Plain")},
		{ui.BrightGreen("Commutativity"), ui.Cyan("A ⊕ B = B ⊕ A"), ui.BrightYellow("Plain ⊕ Key = Key ⊕ Plain")},
		{ui.BrightGreen("Self-Inverse"), ui.Cyan("A ⊕ A = 0"), ui.BrightYellow("Same input bits cancel out")},
		{ui.BrightGreen("Identity"), ui.Cyan("A ⊕ 0 = A"), ui.BrightYellow("XOR with 0 preserves value")},
	})
	propertiesTable.Render()
	fmt.Println()
}

func (r *RC4) printResults() {
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightGreenBold("🎉 Complete encryption results with step-by-step breakdown:"))
	ui.PrintStarLine(40)
	fmt.Println()

	// Encryption summary table
	summaryTable := ui.GetTableWriter()
	summaryTable.SetTitle(ui.BrightMagentaBold("📊 Encryption Summary"))
	summaryTable.AppendHeader(table.Row{
		ui.CyanBold("Step"),
		ui.YellowBold("Input"),
		ui.GreenBold("Keystream"),
		ui.RedBold("Output"),
		ui.MagentaBold("Operation"),
	})

	for i := 0; i < len(r.plaintextASC); i++ {
		inputChar := string(rune(r.plaintextASC[i]))
		outputChar := string(rune(r.ciphertext[i]))
		operation := fmt.Sprintf("%d ⊕ %d = %d", r.plaintextASC[i], r.keystream[i], r.ciphertext[i])
		summaryTable.AppendRow(table.Row{
			ui.BrightCyan(fmt.Sprintf("%d", i+1)),
			ui.BrightYellow(fmt.Sprintf("'%c' (%d)", inputChar[0], r.plaintextASC[i])),
			ui.BrightGreen(fmt.Sprintf("%d", r.keystream[i])),
			ui.BrightRed(fmt.Sprintf("'%c' (%d)", outputChar[0], r.ciphertext[i])),
			ui.Magenta(operation),
		})
	}
	summaryTable.Render()
	fmt.Println()

	// Final results table
	finalTable := ui.GetTableWriter()
	finalTable.SetTitle(ui.BrightGreenBold("🎯 Final Results"))

	plaintextStr := fmt.Sprintf("\"%s\"", r.plaintext)

	// Build ciphertext string
	ciphertextChars := ""
	for _, v := range r.ciphertext {
		ciphertextChars += string(rune(v))
	}

	cipherStr := "["
	for i, v := range r.ciphertext {
		if i > 0 {
			cipherStr += ", "
		}
		cipherStr += fmt.Sprintf("%d", v)
	}
	cipherStr += "]"

	keystreamStr := "["
	for i, v := range r.keystream {
		if i > 0 {
			keystreamStr += ", "
		}
		keystreamStr += fmt.Sprintf("%d", v)
	}
	keystreamStr += "]"

	finalTable.AppendHeader(table.Row{
		ui.MagentaBold("Type"),
		ui.CyanBold("Value"),
	})
	finalTable.AppendRows([]table.Row{
		{ui.BrightYellow("Plaintext"), ui.BrightYellow(plaintextStr + fmt.Sprintf(" → %v", r.plaintextASC))},
		{ui.BrightRed("Ciphertext"), ui.BrightRed(fmt.Sprintf("\"%s\" → %s", ciphertextChars, cipherStr))},
		{ui.BrightGreen("Keystream"), ui.BrightGreen(keystreamStr)},
	})
	finalTable.Render()
	fmt.Println()

	// Verification table
	decrypted := make([]int, len(r.ciphertext))
	for i := 0; i < len(r.ciphertext); i++ {
		decrypted[i] = r.ciphertext[i] ^ r.keystream[i]
	}

	decryptedStr := ""
	for _, v := range decrypted {
		decryptedStr += string(rune(v))
	}

	verifyTable := ui.GetTableWriter()
	verifyTable.SetTitle(ui.BrightCyanBold("🔓 Verification (Decryption)"))
	verifyTable.AppendHeader(table.Row{
		ui.MagentaBold("Component"),
		ui.YellowBold("Value"),
	})
	verifyTable.AppendRows([]table.Row{
		{ui.BrightRed("Ciphertext"), ui.Red(fmt.Sprintf("\"%s\" → %s", ciphertextChars, cipherStr))},
		{ui.BrightGreen("Keystream"), ui.Green(keystreamStr + " (same as encryption)")},
		{ui.BrightYellow("Decrypted"), ui.BrightGreen(fmt.Sprintf("%v → \"%s\" ✓", decrypted, decryptedStr))},
	})
	verifyTable.Render()
	fmt.Println()

	// Summary checklist
	summaryCheckTable := ui.GetTableWriter()
	summaryCheckTable.SetTitle(ui.BrightGreenBold("✅ SUMMARY"))
	summaryCheckTable.AppendRows([]table.Row{
		{ui.BrightGreen("✓"), ui.Cyan("KSA Phase: Successfully scrambled the state array")},
		{ui.BrightGreen("✓"), ui.Cyan("PRGA Phase: Generated keystream and encrypted plaintext")},
		{ui.BrightGreen("✓"), ui.Cyan("XOR Operations: Performed bit-level encryption")},
		{ui.BrightGreen("✓"), ui.Cyan("Reversibility: Verified encryption/decryption process")},
	})
	summaryCheckTable.Render()

	fmt.Println()
	ui.PrintColorfulLine(80)
	fmt.Println(ui.Rainbow("★ ★ ★  R C 4   E N C R Y P T I O N   C O M P L E T E  ★ ★ ★"))
	ui.PrintColorfulLine(80)
	fmt.Println()
}
