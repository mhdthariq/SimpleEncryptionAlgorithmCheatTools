package algorithms

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"strings"

	"example.com/cheat-encryption-algorithm/pkg/ui"
	"github.com/jedib0t/go-pretty/v6/table"
)

type ChaCha20 struct {
	key        []byte
	nonce      []byte
	counter    uint32
	initial    [16]uint32
	state      [16]uint32
	keystream  []byte
	plaintext  []byte
	ciphertext []byte
}

func NewChaCha20() *ChaCha20 {
	return &ChaCha20{}
}

func (c *ChaCha20) GetName() string {
	return "ChaCha20"
}

func (c *ChaCha20) Run(reader *bufio.Reader, plaintextInput string, keyInput string) {
	// ChaCha20 Specific Input: Nonce
	// Key needs to be 32 bytes. We will pad or hash. For visual simplicity, we'll just pad.
	// Nonce needs to be 12 bytes.

	// Process Key (32 bytes)
	keyBytes := make([]byte, 32)
	copy(keyBytes, []byte(keyInput))

	fmt.Print("Please input Nonce (optional, default 0): ")
	nonceInput, _ := reader.ReadString('\n')
	nonceInput = strings.TrimSpace(nonceInput)

	nonceBytes := make([]byte, 12)
	copy(nonceBytes, []byte(nonceInput))

	// Initial Counter
	var counter uint32 = 1

	c.key = keyBytes
	c.nonce = nonceBytes
	c.counter = counter
	c.plaintext = []byte(plaintextInput)

	ui.PrintSectionHeader("CHACHA20 ENCRYPT ALGORITHM - STEP BY STEP")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightCyanBold("🔐 This demonstration will walk you through each step of the ChaCha20 encryption process."))
	fmt.Println(ui.BrightYellow("📚 Follow the colors to understand each transformation!"))
	ui.PrintStarLine(40)
	fmt.Println()

	c.printSetup()
	c.initializeState()
	c.runQuarterRounds()
	c.generateKeystreamAndEncrypt()

	ui.PrintSectionHeader("COMPLETE RESULTS")
	c.printResults()
}

func (c *ChaCha20) printSetup() {
	ui.PrintSectionHeader("SETUP & CONSTANTS")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightCyanBold("✨ Initial parameters for ChaCha20 encryption:"))
	ui.PrintStarLine(40)
	fmt.Println()

	t := ui.GetTableWriter()

	t.AppendHeader(table.Row{
		ui.MagentaBold("Component"),
		ui.CyanBold("Value (Bytes)"),
		ui.YellowBold("Hex Representation"),
	})
	t.AppendRows([]table.Row{
		{ui.BrightMagenta("Key"), ui.BrightYellow(fmt.Sprintf("%q (padded to 32B)", string(c.key))), ui.BrightGreen(fmt.Sprintf("%x...", c.key[:8]))},
		{ui.BrightMagenta("Nonce"), ui.BrightYellow(fmt.Sprintf("%q (padded to 12B)", string(c.nonce))), ui.BrightGreen(fmt.Sprintf("%x", c.nonce))},
		{ui.BrightMagenta("Counter"), ui.BrightCyan(fmt.Sprintf("%d", c.counter)), ui.BrightGreen(fmt.Sprintf("%08x", c.counter))},
		{ui.BrightMagenta("Constant"), ui.BrightYellow("\"expand 32-byte k\""), ui.BrightGreen("65787061 6e642033 322d6279 7465206b")},
	})
	t.Render()
	fmt.Println()
}

func (c *ChaCha20) initializeState() {
	// Constants "expand 32-byte k"
	c.initial[0] = 0x61707865
	c.initial[1] = 0x3320646e
	c.initial[2] = 0x79622d32
	c.initial[3] = 0x6b206574

	// Key
	for i := 0; i < 8; i++ {
		c.initial[4+i] = binary.LittleEndian.Uint32(c.key[i*4 : (i+1)*4])
	}

	// Counter
	c.initial[12] = c.counter

	// Nonce
	for i := 0; i < 3; i++ {
		c.initial[13+i] = binary.LittleEndian.Uint32(c.nonce[i*4 : (i+1)*4])
	}

	c.state = c.initial // Copy to working state

	ui.PrintSectionHeader("INITIAL STATE MATRIX (4x4)")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightYellowBold("🎯 Building the initial 4x4 state matrix from constants, key, counter, and nonce:"))
	ui.PrintStarLine(40)
	fmt.Println()
	c.printMatrix(c.initial, ui.BrightMagentaBold("📊 Initial State"))
	fmt.Println(ui.BrightCyan("  🔹 Row 0: Constants ('expand 32-byte k')"))
	fmt.Println(ui.BrightCyan("  🔹 Rows 1-2: Key (256 bits)"))
	fmt.Println(ui.BrightCyan("  🔹 Row 3: Counter + Nonce (96 bits)"))
	fmt.Println()
}

func (c *ChaCha20) printMatrix(state [16]uint32, title string) {
	fmt.Println(title + ":")
	for i := 0; i < 4; i++ {
		fmt.Printf("%s[%s, %s, %s, %s]%s\n",
			ui.ColorBrightGreen,
			ui.ColorBrightYellow+fmt.Sprintf("%08x", state[i*4+0])+ui.ColorBrightGreen,
			ui.ColorBrightYellow+fmt.Sprintf("%08x", state[i*4+1])+ui.ColorBrightGreen,
			ui.ColorBrightYellow+fmt.Sprintf("%08x", state[i*4+2])+ui.ColorBrightGreen,
			ui.ColorBrightYellow+fmt.Sprintf("%08x", state[i*4+3])+ui.ColorBrightGreen,
			ui.ColorReset)
	}
	fmt.Println()
}

// ARX helper
func rotl(v uint32, n uint) uint32 {
	return (v << n) | (v >> (32 - n))
}

func quarterRound(x *[16]uint32, a, b, c, d int) {
	x[a] += x[b]
	x[d] ^= x[a]
	x[d] = rotl(x[d], 16)
	x[c] += x[d]
	x[b] ^= x[c]
	x[b] = rotl(x[b], 12)
	x[a] += x[b]
	x[d] ^= x[a]
	x[d] = rotl(x[d], 8)
	x[c] += x[d]
	x[b] ^= x[c]
	x[b] = rotl(x[b], 7)
}

func (c *ChaCha20) runQuarterRounds() {
	ui.PrintSectionHeader("THE QUARTER ROUND (Visualization)")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightYellowBold("🔑 Phase 1: Applying the ChaCha20 quarter-round function"))
	ui.PrintSeparator()
	ui.PrintStarLine(40)
	fmt.Println()

	fmt.Println(ui.BrightCyanBold("ℹ️  ChaCha20 runs 20 rounds (10 loops of 2 rounds each: Column rounds + Diagonal rounds)."))
	fmt.Println(ui.BrightCyanBold("ℹ️  Each round consists of 4 Quarter Rounds."))
	fmt.Println()
	fmt.Println(ui.BrightMagentaBold("⚙️  Calculating 20 rounds..."))
	fmt.Println()

	// Copy state to work on
	x := c.initial

	for i := 0; i < 10; i++ { // 10 iterations * 2 rounds = 20 rounds
		// Odd round (Column)
		quarterRound(&x, 0, 4, 8, 12)
		quarterRound(&x, 1, 5, 9, 13)
		quarterRound(&x, 2, 6, 10, 14)
		quarterRound(&x, 3, 7, 11, 15)

		// Even round (Diagonal)
		quarterRound(&x, 0, 5, 10, 15)
		quarterRound(&x, 1, 6, 11, 12)
		quarterRound(&x, 2, 7, 8, 13)
		quarterRound(&x, 3, 4, 9, 14)
	}

	c.state = x
	fmt.Println()
	fmt.Println(ui.BrightGreenBold("✅ 20 rounds completed successfully!"))
	fmt.Println()
	c.printMatrix(c.state, ui.BrightMagentaBold("📊 State After 20 Rounds"))
}

func (c *ChaCha20) generateKeystreamAndEncrypt() {
	ui.PrintSectionHeader("KEYSTREAM GENERATION & XOR")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightYellowBold("🔐 Phase 2: Generating keystream and encrypting plaintext"))
	ui.PrintSeparator()
	ui.PrintStarLine(40)
	fmt.Println()

	// Final add
	fmt.Println(ui.BrightMagentaBold("🎯 Final Step: Add Final State to Initial State (Word-wise)"))

	var block [64]byte
	for i := 0; i < 16; i++ {
		c.state[i] = c.state[i] + c.initial[i]
		binary.LittleEndian.PutUint32(block[i*4:], c.state[i])
	}

	c.printMatrix(c.state, ui.BrightMagentaBold("📊 Final Keystream Block (State + Initial)"))

	fmt.Println(ui.BrightGreenBold("🔑 Resulting Keystream Block (64 bytes):"))
	fmt.Printf("%s%x%s\n", ui.ColorBrightYellow, block, ui.ColorReset)
	fmt.Println()

	ui.PrintSectionHeader("ENCRYPTION")
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightCyanBold("⚡ Performing byte-wise XOR between plaintext and keystream:"))
	ui.PrintStarLine(40)
	fmt.Println()

	c.keystream = block[:len(c.plaintext)] // Truncate if plaintext < 64
	c.ciphertext = make([]byte, len(c.plaintext))

	fmt.Println(ui.BrightMagentaBold("🔐 Byte-wise XOR Operations:"))
	fmt.Println()

	for i, b := range c.plaintext {
		// handle multi-block? For this simple tool, we assume inputs < 64 bytes or just truncation for demo.
		// If input > 64, we would need next block.
		// For safety in this "cheat" demo, we just use the first block.
		k := block[i%64]
		c.ciphertext[i] = b ^ k

		fmt.Printf("%s %s: %s (%s) %s %s = %s\n",
			ui.CyanBold("📍 Index"),
			ui.BrightCyan(fmt.Sprintf("%d", i)),
			ui.BrightYellowBold(fmt.Sprintf("'%c'", b)),
			ui.BrightYellow(fmt.Sprintf("%02x", b)),
			ui.MagentaBold("⊕ XOR ⊕"),
			ui.BrightGreen(fmt.Sprintf("%02x", k)),
			ui.BrightRedBold(fmt.Sprintf("%02x", c.ciphertext[i])))
		fmt.Printf("%s Ciphertext byte %s = %s 🎯\n\n",
			ui.GreenBold("🎉 Result"),
			ui.BrightRedBold(fmt.Sprintf("%d", c.ciphertext[i])),
			ui.BrightRedBold(fmt.Sprintf("'%c'", rune(c.ciphertext[i]))))
	}
	fmt.Println()
}

func (c *ChaCha20) printResults() {
	ui.PrintStarLine(40)
	fmt.Println(ui.BrightGreenBold("🎉 Complete encryption results:"))
	ui.PrintStarLine(40)
	fmt.Println()

	finalTable := ui.GetTableWriter()
	finalTable.SetTitle(ui.BrightGreenBold("🎯 Final Results"))

	// Build plaintext string
	plaintextStr := string(c.plaintext)
	ciphertextStr := string(c.ciphertext)

	finalTable.AppendHeader(table.Row{
		ui.MagentaBold("Type"),
		ui.CyanBold("Text"),
		ui.YellowBold("Value (Hex)"),
	})
	finalTable.AppendRows([]table.Row{
		{ui.BrightYellow("Plaintext"), ui.BrightYellow(fmt.Sprintf("\"%s\"", plaintextStr)), ui.BrightGreen(fmt.Sprintf("%x", c.plaintext))},
		{ui.BrightGreen("Keystream"), ui.Green("-"), ui.BrightGreen(fmt.Sprintf("%x", c.keystream))},
		{ui.BrightRed("Ciphertext"), ui.BrightRed(fmt.Sprintf("\"%s\"", ciphertextStr)), ui.BrightRed(fmt.Sprintf("%x", c.ciphertext))},
	})
	finalTable.Render()
	fmt.Println()

	// Summary checklist
	summaryCheckTable := ui.GetTableWriter()
	summaryCheckTable.SetTitle(ui.BrightGreenBold("✅ SUMMARY"))
	summaryCheckTable.AppendRows([]table.Row{
		{ui.BrightGreen("✓"), ui.Cyan("State Initialization: Built 4x4 matrix from constants, key, counter, and nonce")},
		{ui.BrightGreen("✓"), ui.Cyan("Quarter Rounds: Applied 20 rounds of ARX operations")},
		{ui.BrightGreen("✓"), ui.Cyan("Keystream Generation: Added initial state to final state")},
		{ui.BrightGreen("✓"), ui.Cyan("Encryption: XORed plaintext with keystream to produce ciphertext")},
	})
	summaryCheckTable.Render()

	fmt.Println()
	ui.PrintColorfulLine(80)
	fmt.Println(ui.Rainbow("★ ★ ★  C H A C H A 2 0   E N C R Y P T I O N   C O M P L E T E  ★ ★ ★"))
	ui.PrintColorfulLine(80)
	fmt.Println()
}
