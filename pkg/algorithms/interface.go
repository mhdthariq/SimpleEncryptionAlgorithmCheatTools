package algorithms

import "bufio"

// Algorithm interface for extensibility
type Algorithm interface {
	GetName() string
	// Run handles the full interaction: asking for specific params, encrypting, and showing steps
	Run(reader *bufio.Reader, plaintext string, key string)
}
