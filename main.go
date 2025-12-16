package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"example.com/cheat-encryption-algorithm/pkg/algorithms"
	"example.com/cheat-encryption-algorithm/pkg/ui"

	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		ui.ClearScreen()

		// Welcome banner
		ui.PrintSectionHeader("ENCRYPTION ALGORITHM CALCULATOR - MODULAR DESIGN")
		ui.PrintStarLine(40)
		fmt.Println(ui.BrightCyanBold("🔐 Welcome! Choose an encryption algorithm to visualize step-by-step."))
		fmt.Println(ui.BrightYellow("📚 Learn cryptography through beautiful, colorful visualizations!"))
		ui.PrintStarLine(40)
		fmt.Println()

		// Algorithm selection
		algoTable := ui.GetTableWriter()
		algoTable.SetTitle(ui.BrightMagentaBold("✨ Available Algorithms"))
		algoTable.AppendHeader(table.Row{
			ui.CyanBold("#"),
			ui.MagentaBold("Algorithm"),
			ui.YellowBold("Description"),
		})
		algoTable.AppendRows([]table.Row{
			{ui.BrightYellow("1"), ui.BrightGreen("🔐 RC4"), ui.Cyan("Stream cipher encryption")},
			{ui.BrightYellow("2"), ui.BrightGreen("⚡ ChaCha20"), ui.Cyan("Modern stream cipher (Visualized)")},
			{ui.BrightYellow("3"), ui.BrightGreen("📜 Vigenère"), ui.Cyan("Classic polyalphabetic cipher")},
			{ui.BrightYellow("4"), ui.BrightRed("🚪 Exit"), ui.Cyan("Exit the program")},
		})
		algoTable.Render()
		fmt.Println()

		fmt.Print(ui.BrightCyanBold("👉 Select option: "))
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "4" || strings.ToLower(choice) == "exit" {
			fmt.Println()
			ui.PrintColorfulLine(50)
			fmt.Println(ui.BrightGreenBold("✨ Thank you for using the Encryption Algorithm Calculator!"))
			fmt.Println(ui.BrightCyanBold("🎓 Keep learning and stay curious about cryptography!"))
			fmt.Println(ui.Rainbow("★ ★ ★  G O O D B Y E  ★ ★ ★"))
			ui.PrintColorfulLine(50)
			fmt.Println()
			break
		}

		var algo algorithms.Algorithm

		switch choice {
		case "1", "rc4":
			algo = algorithms.NewRC4()
		case "2", "chacha20":
			algo = algorithms.NewChaCha20()
		case "3", "vigenere":
			algo = algorithms.NewVigenere()
		default:
			fmt.Println()
			fmt.Println(ui.RedBold("❌ Invalid choice! Please select a valid option."))
			fmt.Println(ui.YellowBold("💡 Hint: Choose 1, 2, 3, or 4"))
			fmt.Println()
			fmt.Print(ui.Cyan("Press Enter to try again..."))
			reader.ReadString('\n')
			continue
		}

		// Common Inputs
		fmt.Println()
		ui.PrintColorfulLine(50)
		fmt.Println(ui.BrightMagentaBold(fmt.Sprintf("═══ %s Encryption ═══", algo.GetName())))
		ui.PrintColorfulLine(50)
		fmt.Println()

		fmt.Print(ui.BrightYellowBold("📝 Please input Plaintext: "))
		plaintext, _ := reader.ReadString('\n')
		plaintext = strings.TrimSpace(plaintext)

		fmt.Print(ui.BrightYellowBold("🔑 Please input Key: "))
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		// Validation
		if plaintext == "" || key == "" {
			fmt.Println()
			fmt.Println(ui.RedBold("❌ Error: Plaintext and key cannot be empty!"))
			fmt.Println(ui.YellowBold("💡 Both fields are required for encryption."))
			fmt.Println()
			fmt.Print(ui.Cyan("Press Enter to continue..."))
			reader.ReadString('\n')
			continue
		}

		// Run Algorithm
		algo.Run(reader, plaintext, key)

		// Wait before clearing
		fmt.Println()
		ui.PrintColorfulLine(50)
		fmt.Println(ui.BrightGreenBold("✅ Encryption process complete!"))
		fmt.Println(ui.BrightCyanBold("🎉 Successfully encrypted your message!"))
		fmt.Println(ui.Rainbow("★ ★ ★  COMPLETE  ★ ★ ★"))
		ui.PrintColorfulLine(50)
		fmt.Println()
		fmt.Print(ui.YellowBold("👉 Press Enter to return to main menu..."))
		reader.ReadString('\n')
	}
}
