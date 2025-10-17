package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseRange(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"A-F", "ABCDEF"},
		{"0-3", "0123"},
		{"A-C,1-2", "ABC12"},
		{"X", "X"},
		{"a-c,Z", "abcZ"},
	}

	for _, tt := range tests {
		got := parseRange(tt.input)
		if got != tt.want {
			t.Errorf("parseRange(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBuildCharset_All(t *testing.T) {
	opts := options{all: true}
	charset := buildCharset(opts)

	if !strings.Contains(charset, "a") ||
		!strings.Contains(charset, "Z") ||
		!strings.Contains(charset, "0") ||
		!strings.Contains(charset, "!") {
		t.Errorf("buildCharset(all=true) missing expected characters: %q", charset)
	}
}

func TestBuildCharset_SpecificFlags(t *testing.T) {
	opts := options{
		lower:  true,
		upper:  true,
		digits: true,
		symbols: false,
	}
	charset := buildCharset(opts)

	if strings.Contains(charset, "!") {
		t.Error("expected no symbols in charset")
	}
	if !strings.Contains(charset, "a") || !strings.Contains(charset, "Z") || !strings.Contains(charset, "9") {
		t.Errorf("missing expected characters: %q", charset)
	}
}

func TestBuildCharset_NoSimilar(t *testing.T) {
	opts := options{
		all:       true,
		noSimilar: true,
	}
	charset := buildCharset(opts)
	for _, s := range similarChars {
		if strings.ContainsRune(charset, s) {
			t.Errorf("charset contains similar char %q", s)
		}
	}
}

func TestBuildCharset_CustomRange(t *testing.T) {
	opts := options{
		custom: "A-C,1-2",
	}
	charset := buildCharset(opts)
	want := "ABC12"
	for _, c := range want {
		if !strings.ContainsRune(charset, c) {
			t.Errorf("expected %q in charset %q", c, charset)
		}
	}
}

func TestGeneratePassword_LengthAndCharset(t *testing.T) {
	charset := "ABC123"
	length := 10
	pass := generatePassword(length, charset)

	if len(pass) != length {
		t.Fatalf("expected password length %d, got %d", length, len(pass))
	}
	for _, c := range pass {
		if !strings.ContainsRune(charset, c) {
			t.Fatalf("password contains invalid char %q", c)
		}
	}
}

func TestGeneratePassword_EmptyCharset(t *testing.T) {
	pass := generatePassword(10, "")
	if pass != "" {
		t.Errorf("expected empty password, got %q", pass)
	}
}

func TestJSONOutputFormat(t *testing.T) {
	results := []string{"abc", "def"}
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent failed: %v", err)
	}

	if !strings.Contains(string(data), "abc") || !strings.Contains(string(data), "def") {
		t.Errorf("JSON output missing expected values: %s", string(data))
	}
}
