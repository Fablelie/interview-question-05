package database

import (
	"fmt"
	"strings"
)

// TicketNumberGenerator handles ticket number generation (A0 to Z9)
type TicketNumberGenerator struct{}

// GetNextNumber returns the next ticket number in sequence
// A0, A1, ..., A9, B0, B1, ..., Z9
func (g *TicketNumberGenerator) GetNextNumber(currentNumber string) (string, error) {
	if len(currentNumber) != 2 {
		return "", fmt.Errorf("invalid ticket number format: %s", currentNumber)
	}

	letter := rune(currentNumber[0])
	digit := rune(currentNumber[1])

	// Increment digit
	digit++

	// If digit exceeds 9, reset to 0 and increment letter
	if digit > '9' {
		digit = '0'
		letter++
	}

	// If letter exceeds Z, reset to A
	if letter > 'Z' {
		letter = 'A'
	}

	return string(letter) + string(digit), nil
}

// IsValidTicketNumber validates if the number is in correct format
func (g *TicketNumberGenerator) IsValidTicketNumber(number string) bool {
	if len(number) != 2 {
		return false
	}

	letter := rune(number[0])
	digit := rune(number[1])

	return letter >= 'A' && letter <= 'Z' && digit >= '0' && digit <= '9'
}

// FormatTicketNumber ensures proper formatting
func (g *TicketNumberGenerator) FormatTicketNumber(number string) string {
	return strings.ToUpper(number)
}
