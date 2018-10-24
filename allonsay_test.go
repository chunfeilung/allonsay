package main

import (
	"os"
	"testing"
)

// -----------------------------------------------------------------------------

func TestMainFunction(t *testing.T) {
	var strings = []struct {
		description string
		arguments   []string
	}{
		{"No argument", []string{"allonsay"}},
		{"Empty argument", []string{"allonsay", ""}},
		{"Empty-ish argument", []string{"allonsay", " "}},
		{"Two arguments", []string{"allonsay", " ", "Allons-y!"}},
	}

	for _, tt := range strings {
		t.Run(string(tt.description), func(t *testing.T) {
			// Arrange
			os.Args = tt.arguments

			// Act & Assert
			main()
			// Just don’t break for now
		})
	}
}

// -----------------------------------------------------------------------------

func TestGetVoice(t *testing.T) {
	var strings = []struct {
		input string
		voice Voice
	}{
		{"University of Hong Kong", EnglishVoice},
		{"Nederlandse Spoorwegen", DutchVoice},
		{"香港大學", CantoneseVoice},
	}

	for _, tt := range strings {
		t.Run(string(tt.input), func(t *testing.T) {
			// Arrange
			var expected = tt.voice

			// Act
			var actual = getVoice(tt.input, tt.input)

			// Assert
			if actual != expected {
				t.Error(tt.input)
			}
		})
	}
}

// -----------------------------------------------------------------------------

func TestGuessLanguage(t *testing.T) {
	var strings = []struct {
		label    string
		input    string
		language Language
	}{
		{
			"English 1",
			`
			Harvard University is devoted to excellence in teaching, learning,
			and research, and to developing leaders in many disciplines who make
			a difference globally. The University, which is based in Cambridge
			and Boston, Massachusetts, has an enrollment of over 20,000 degree
			candidates
			`,
			EnglishLanguage,
		},
		{
			"English 2",
			`
			Royal Schiphol Group brings the world within reach. We connect the
			Netherlands with the rest of the world, thus creating value for the
			economy and society.
			`,
			EnglishLanguage,
		},
		{
			"Dutch 1",
			`
			De Universiteit van Amsterdam is een van de meest vooraanstaande
			onderzoeksuniversiteiten in Europa, een moderne universiteit met een
			lange en rijke geschiedenis
			`,
			DutchLanguage,
		},
		{
			"Dutch 2",
			`
			TNO verbindt mensen en kennis om innovaties te creëren die de
			concurrentiekracht van bedrijven en het welzijn van de samenleving
			duurzaam versterken.
			`,
			DutchLanguage,
		},
		{
			"Dutch 3 (names and English loanwords)",
			`
			InTraffic is in 2003 opgericht als ‘joint venture’ van het
			IT-bedrijf ICT Automatisering en ingenieursbureau Movares. De
			filosofie achter InTraffic is dat alleen een gespecialiseerd bedrijf
			kan beschikken over een hoge mate van domeinkennis in combinatie met
			IT-expertise.
			`,
			DutchLanguage,
		},
	}

	for _, tt := range strings {
		t.Run(string(tt.label), func(t *testing.T) {
			// Arrange
			var expected = tt.language

			// Act
			var actual = guessLanguage(tt.input)

			// Assert
			if actual != expected {
				t.Error("Guessed wrong language")
			}
		})
	}
}

// -----------------------------------------------------------------------------

func TestCalculateLetterFrequencies(t *testing.T) {
	// Arrange
	var input = "Kowloon Tong (九龍塘)"

	// Act
	var actual = calculateLetterFrequencies(input)

	// Assert
	if actual['k'] != float64(1*100)/float64(11) ||
		actual['n'] != float64(2*100)/float64(11) ||
		actual['o'] != float64(4*100)/float64(11) ||
		actual['t'] != float64(1*100)/float64(11) {
		t.Error("Letter frequency mismatch")
	}
}

// -----------------------------------------------------------------------------

func TestConvertStringToBilingualSentence(t *testing.T) {
	// Arrange
	var input = "You can change trains at 美孚 station"
	var expected = BilingualSentence{
		"You can change trains at ",
		"美孚",
		" station",
	}

	// Act
	var actual = convertStringToBilingualSentence(input)

	// Assert
	if len(actual) != len(expected) {
		t.Error("Resulting BilingualSentence does not have expected length")
	}
	for i := 0; i < 3; i++ {
		if actual[i] != expected[i] {
			t.Error("Result does not exactly match expected BilingualSentence")
		}
	}
}

// -----------------------------------------------------------------------------

func TestOnlyAlphabeticCharacters(t *testing.T) {
	// Arrange
	var input = "Kennedy Town (堅尼地城) – Chai Wan (柴灣)"
	var expected = "KennedyTownChaiWan"

	// Act
	var actual = string(onlyAlphabeticCharacters(input))

	// Assert
	if actual != expected {
		t.Error("Output still contains non-alphabetic characters")
	}
}

// -----------------------------------------------------------------------------

func TestIsCharacterChinese(t *testing.T) {
	var possiblyChineseCharacters = []struct {
		description string
		char        rune
		isCjk       bool
	}{
		{"Uppercase A", 'A', false},
		{"Lowercase A", 'a', false},
		{"Circled small letter A", 'ⓐ', false},
		{"Mathematical double-struck A", '𝔸', false},
		{"Hebrew alef", 'א', false},
		{"Traditional xué", '學', true},
		{"Simplified xué", '学', true},
	}

	for _, tt := range possiblyChineseCharacters {
		t.Run(string(tt.char), func(t *testing.T) {
			// Arrange
			var expected = tt.isCjk

			// Act
			var actual = isCharacterChinese(tt.char)

			// Assert
			if actual != expected {
				t.Error(tt.description)
			}
		})
	}
}
