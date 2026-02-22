package game

import "testing"

func TestParseInput_Empty(t *testing.T) {
	verb, noun := ParseInput("")
	if verb != "" || noun != "" {
		t.Errorf("Expected empty verb and noun, got verb=%q noun=%q", verb, noun)
	}
}

func TestParseInput_SingleWord(t *testing.T) {
	verb, noun := ParseInput("look")
	if verb != "look" {
		t.Errorf("Expected verb 'look', got %q", verb)
	}
	if noun != "" {
		t.Errorf("Expected empty noun, got %q", noun)
	}
}

func TestParseInput_TwoWords(t *testing.T) {
	verb, noun := ParseInput("go north")
	if verb != "go" {
		t.Errorf("Expected verb 'go', got %q", verb)
	}
	if noun != "north" {
		t.Errorf("Expected noun 'north', got %q", noun)
	}
}

func TestParseInput_MultiWordNoun(t *testing.T) {
	verb, noun := ParseInput("take rusty old key")
	if verb != "take" {
		t.Errorf("Expected verb 'take', got %q", verb)
	}
	if noun != "rusty old key" {
		t.Errorf("Expected noun 'rusty old key', got %q", noun)
	}
}

func TestParseInput_ExtraWhitespace(t *testing.T) {
	verb, noun := ParseInput("  go   north  ")
	if verb != "go" {
		t.Errorf("Expected verb 'go', got %q", verb)
	}
	if noun != "north" {
		t.Errorf("Expected noun 'north', got %q", noun)
	}
}
