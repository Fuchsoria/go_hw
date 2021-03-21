package hw03frequencyanalysis

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

var (
	wordRegex      = regexp.MustCompile(`[^\s]+`)
	clearWordRegex = regexp.MustCompile(`[^a-zA-Zа-яА-Яё]+`)
)

type Word struct {
	key   string
	count int
}

func clearWord(word string) (string, error) {
	str := clearWordRegex.ReplaceAllString(word, "")

	if len(str) == 0 {
		return "", errors.New("Word is empty")
	}

	return strings.ToLower(str), nil
}

func Top10(text string) []string {
	wordsByString := wordRegex.FindAllString(text, -1)
	maps := make(map[string]int)
	words := []Word{}
	result := []string{}

	for _, word := range wordsByString {
		word, wordError := clearWord(word)

		if wordError != nil {
			continue
		}

		_, isWordExist := maps[word]

		if isWordExist {
			maps[word]++
		} else {
			maps[word] = 1
		}
	}

	for k, v := range maps {
		words = append(words, Word{key: k, count: v})
	}

	sort.Slice(words, func(a, b int) bool {
		aWord, bWord := words[a], words[b]

		if aWord.count == bWord.count {
			return aWord.key < bWord.key
		}

		return aWord.count > bWord.count
	})

	for i, word := range words {
		if i < 10 {
			result = append(result, word.key)
		} else {
			break
		}
	}

	return result
}
