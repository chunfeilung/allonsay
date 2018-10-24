package main

import (
	"math"
	"os"
	"os/exec"
	"strings"
)

// -----------------------------------------------------------------------------

const EnglishVoice = "Samantha"
const DutchVoice = "Claire"
const CantoneseVoice = "Sin-ji"

// -----------------------------------------------------------------------------

// This is a poor man’s enum

const EnglishLanguage = 0
const DutchLanguage = 1

// -----------------------------------------------------------------------------

// These are all taken from Wikipedia. Manual training might yield slightly
// better results if a user often reads niche content, but it’s not really worth
// the effort for a simple tool like this.

var EnglishFrequencies = FrequencyTable{
	'a': 8.167,
	'b': 1.492,
	'c': 2.782,
	'd': 4.253,
	'e': 12.702,
	'f': 2.228,
	'g': 2.015,
	'h': 6.094,
	'i': 6.966,
	'j': 0.153,
	'k': 0.772,
	'l': 4.025,
	'm': 2.406,
	'n': 6.749,
	'o': 7.507,
	'p': 1.929,
	'q': 0.095,
	'r': 5.987,
	's': 6.327,
	't': 9.056,
	'u': 2.758,
	'v': 0.978,
	'w': 2.360,
	'x': 0.150,
	'y': 1.974,
	'z': 0.074,
}

var DutchFrequencies = FrequencyTable{
	'a': 7.486,
	'b': 1.584,
	'c': 1.242,
	'd': 5.933,
	'e': 18.91,
	'f': 0.805,
	'g': 3.403,
	'h': 2.380,
	'i': 6.499,
	'j': 1.46,
	'k': 2.248,
	'l': 3.568,
	'm': 2.213,
	'n': 10.032,
	'o': 6.063,
	'p': 1.57,
	'q': 0.009,
	'r': 6.411,
	's': 3.73,
	't': 6.79,
	'u': 1.99,
	'v': 2.85,
	'w': 1.52,
	'x': 0.036,
	'y': 0.035,
	'z': 1.39,
}

var LetterFrequencies = map[Language]FrequencyTable{
	EnglishLanguage: EnglishFrequencies,
	DutchLanguage:   DutchFrequencies,
}

// -----------------------------------------------------------------------------

// BilingualSentence represents a sentence that contains one or two languages
type BilingualSentence = []SentencePart

// FrequencyTable stores relative letter frequencies within a string
type FrequencyTable = map[rune]float64

// Language is one of the supported languages
type Language = int

// SentencePart is a sequence of characters that can be pronounced correctly
// using a single system voice
type SentencePart = string

// Voice is the name of one of macOS’ system voices
type Voice = string

// -----------------------------------------------------------------------------

func main() {
	if len(os.Args) <= 1 {
		// Don’t bother the user
		return
	}

	var input = os.Args[1]

	if len(input) == 0 {
		// Still don’t bother the user
		return
	}

	var bilingualSentence = convertStringToBilingualSentence(input)

	// Each call to exec.Command creates a new shell, which would introduce
	// awkward pauses whenever a sentence continues in a different language.
	// Grouping everything into a “sh -c <sequence of commands>” circumvents
	// that issue.
	var commands = []string{}
	for _, sentencePart := range bilingualSentence {
		var voice = getVoice(sentencePart, input)
		var command = "say -v " + voice + " \"" + sentencePart + "\""
		commands = append(commands, command)
	}
	exec.Command("sh", "-c", strings.Join(commands, ";")).Output()
}

// -----------------------------------------------------------------------------

// getVoice returns the name of the system voice that should be used for the
// sentence part.
func getVoice(part SentencePart, sentence string) Voice {
	if isCharacterChinese([]rune(part)[0]) {
		return CantoneseVoice
	}

	// We’re assuming that any text surrounding sequences of Chinese characters
	// will be in a single language. This assumption is clearly wrong in many
	// cases, but it also lowers the chance that single-language texts are
	// misclassified (keep in mind that we determine the language using letter
	// frequency, which doesn’t work well with short texts!).
	var language = guessLanguage(sentence)

	if language == DutchLanguage {
		return DutchVoice
	}

	return EnglishVoice
}

// -----------------------------------------------------------------------------

// guessLanguage returns the language that’s most likely to be used in the
// sentence.
func guessLanguage(sentence string) Language {
	var sentenceFrequencies = calculateLetterFrequencies(sentence)
	var differences = make(map[Language]float64)

	for language, languageFrequencies := range LetterFrequencies {
		var difference = float64(0)

		for letter, frequency := range languageFrequencies {
			difference += math.Pow(frequency-sentenceFrequencies[letter], 2)
		}

		differences[language] = difference / 26
	}

	var closestLanguage = EnglishLanguage
	var lowestDifference = math.MaxFloat64

	for language, difference := range differences {
		if difference < lowestDifference {
			closestLanguage = language
			lowestDifference = difference
		}
	}

	return closestLanguage
}

// -----------------------------------------------------------------------------

// calculateLetterFrequencies determines the relative frequency of characters
// a-z in a sentence. This function is case-insensitive.
func calculateLetterFrequencies(sentence string) FrequencyTable {
	var absoluteFrequencies = make(map[rune]int)
	var relativeFrequencies = make(FrequencyTable)
	var total = 0

	for _, char := range "abcdefghijklmnopqrstuvwxyz" {
		absoluteFrequencies[char] = 0
	}

	for _, char := range onlyAlphabeticCharacters(strings.ToLower(sentence)) {
		absoluteFrequencies[char] += 1
		total += 1
	}

	for char, count := range absoluteFrequencies {
		relativeFrequencies[char] = float64(count*100) / float64(total)
	}

	return relativeFrequencies
}

// -----------------------------------------------------------------------------

// convertStringToBilingualSentence splits a sentence into parts, such that only
// a single voice is needed for each sentence part.
func convertStringToBilingualSentence(sentence string) BilingualSentence {
	var bilingualSentence = BilingualSentence{}
	var sentencePart = ""
	var isCurrentlyChinese = false

	for _, char := range sentence {
		var isCharacterChinese = isCharacterChinese(char)

		if isCharacterChinese != isCurrentlyChinese {
			if len(sentencePart) > 0 {
				bilingualSentence = append(bilingualSentence, sentencePart)
			}
			sentencePart = ""
		}

		sentencePart += string(char)
		isCurrentlyChinese = isCharacterChinese
	}

	if len(sentencePart) > 0 {
		bilingualSentence = append(bilingualSentence, sentencePart)
	}

	return bilingualSentence
}

// -----------------------------------------------------------------------------

// onlyAlphabeticCharacters returns all characters from sentence that are not
// within the range A-Z or a-z.
func onlyAlphabeticCharacters(sentence string) []rune {
	var onlyAlphabeticCharacters = []rune{}

	for _, char := range sentence {
		if isCharacterAlphabetic(char) {
			onlyAlphabeticCharacters = append(onlyAlphabeticCharacters, char)
		}
	}

	return onlyAlphabeticCharacters
}

// -----------------------------------------------------------------------------

// isCharacterAlphabetic returns true if the passed rune is A-Z or a-z.
func isCharacterAlphabetic(char rune) bool {
	return char >= 0x41 && char <= 0x5a || char >= 0x61 && char <= 0x7a
}

// -----------------------------------------------------------------------------

// isCharacterChinese returns true if the passed rune is a Chinese character.
func isCharacterChinese(char rune) bool {
	// CJK Unified Ideographs + Extensions A–D
	return char >= 0x4e00 && char <= 0x9fff ||
		char >= 0x3400 && char <= 0x4dbf ||
		char >= 0x20000 && char <= 0x2a6df ||
		char >= 0x2a700 && char <= 0x2b73f ||
		char >= 0x2b740 && char <= 0x2b81f
}
