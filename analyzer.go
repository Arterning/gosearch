package main

import (
	"strings"
	"unicode"
)

// analyze tokenizes and normalizes text for indexing
func analyze(text string) []string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Split into tokens
	var tokens []string
	var currentToken strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			currentToken.WriteRune(r)
		} else {
			if currentToken.Len() > 0 {
				token := currentToken.String()
				if len(token) >= 2 { // Filter out single character tokens
					tokens = append(tokens, token)
				}
				currentToken.Reset()
			}
		}
	}

	// Add last token if exists
	if currentToken.Len() > 0 {
		token := currentToken.String()
		if len(token) >= 2 {
			tokens = append(tokens, token)
		}
	}

	return tokens
}
