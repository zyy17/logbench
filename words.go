package main

import (
	"math/rand"
	"strings"

	"github.com/zyy17/logbench/pkg/words"
)

// randomWords generates a string of random words with total size close to but not exceeding the specified size
func randomWords(size int) string {
	if size <= 0 {
		return ""
	}

	var result []string
	currentSize := 0

	for currentSize < size {
		word := randomWord()
		if word == "" {
			continue
		}

		// Check if adding this word (plus a space) would exceed the size
		newSize := currentSize
		if len(result) > 0 {
			newSize += 1 // space
		}
		newSize += len(word)

		result = append(result, word)
		currentSize = newSize
	}

	return strings.Join(result, " ")
}

func randomWord() string {
	return words.Words[rand.Intn(len(words.Words))]
}
