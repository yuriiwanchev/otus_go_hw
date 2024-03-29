package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dlclark/regexp2"
)

func Top10(input string) []string {
	if input == "" {
		return []string{}
	}

	words := strings.Fields(input)

	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[word]++
	}

	wordList := make([]string, 0, 10)
	for word := range wordCount {
		wordList = append(wordList, word)
	}

	sort.Slice(wordList, func(i, j int) bool {
		if wordCount[wordList[i]] == wordCount[wordList[j]] {
			return wordList[i] < wordList[j]
		}
		return wordCount[wordList[i]] > wordCount[wordList[j]]
	})

	return wordList[:10]
}

func Top10Asterisk(input string) []string {
	if input == "" {
		return []string{}
	}

	var regexStringBuilder strings.Builder
	// Для нахождения небукв под конец слова
	regexStringBuilder.WriteString(`(?<=[\p{L}\p{N}])[^\p{L}\p{N}\s]+(?=[\s])|`)
	// Для нахождения небукв перед началом слова
	regexStringBuilder.WriteString(`(?<=[\s])[^\p{L}\p{N}\s]+(?=[\p{L}\p{N}])|`)
	// Для нахождения одной небуквы между пробелами
	regexStringBuilder.WriteString(`(?<=[\s])[^\p{L}\p{N}\s](?=[\s])|`)
	// Для нахождения одной небуквы под конец слова (даже если это конец текста)
	regexStringBuilder.WriteString(`(?<=[\p{L}\p{N}])[^\p{L}\p{N}]`)
	regexStringBuilder.WriteString(`(?<=[^\p{L}\p{N}\s])(?![\p{L}\p{N}\s])(?![^\p{L}\p{N}\s])`)

	lower := strings.ToLower(input)
	re := regexp2.MustCompile(regexStringBuilder.String(), 0)
	normalized, err := re.Replace(lower, "", 0, -1)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Work failed:", err)
		}
	}()

	return Top10(normalized)
}
