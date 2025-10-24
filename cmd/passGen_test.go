package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseRange(t *testing.T) {
	cases := map[string]string{
		"A-F":     "ABCDEF",
		"0-3":     "0123",
		"A-C,1-2": "ABC12",
		"X":       "X",
		"a-c,Z":   "abcZ",
		"":        "",
	}

	for input, want := range cases {
		t.Run(input, func(t *testing.T) {
			if got := parseRange(input); got != want {
				t.Errorf("parseRange(%q) = %q, want %q", input, got, want)
			}
		})
	}
}

func TestBuildCharset_All(t *testing.T) {
	opts := genOptions{all: true}
	got := buildCharset(opts)

	checks := []string{"a", "Z", "0", "!"}
	for _, c := range checks {
		if !strings.Contains(got, c) {
			t.Errorf("charset missing expected char %q: %q", c, got)
		}
	}
}

func TestBuildCharset_SpecificFlags(t *testing.T) {
	opts := genOptions{
		lower:  true,
		upper:  true,
		digits: true,
	}
	got := buildCharset(opts)

	if strings.Contains(got, "!") {
		t.Error("expected no symbols in charset")
	}
	for _, c := range []string{"a", "Z", "9"} {
		if !strings.Contains(got, c) {
			t.Errorf("missing expected char %q in %q", c, got)
		}
	}
}

func TestBuildCharset_NoSimilar(t *testing.T) {
	opts := genOptions{all: true, noSimilar: true}
	got := buildCharset(opts)

	for _, s := range similarChars {
		if strings.ContainsRune(got, s) {
			t.Errorf("charset contains similar char %q", s)
		}
	}
}

func TestBuildCharset_CustomRange(t *testing.T) {
	opts := genOptions{custom: "A-C,1-2"}
	got := buildCharset(opts)

	for _, c := range "ABC12" {
		if !strings.ContainsRune(got, c) {
			t.Errorf("expected %q in charset %q", c, got)
		}
	}
}

func TestGeneratePassword(t *testing.T) {
	t.Run("valid charset", func(t *testing.T) {
		charset := "ABC123"
		length := 10
		pass := generatePassword(length, charset)

		if len(pass) != length {
			t.Errorf("expected length %d, got %d", length, len(pass))
		}
		for _, c := range pass {
			if !strings.ContainsRune(charset, c) {
				t.Errorf("password contains invalid char %q", c)
			}
		}
	})

	t.Run("empty charset", func(t *testing.T) {
		if got := generatePassword(10, ""); got != "" {
			t.Errorf("expected empty password, got %q", got)
		}
	})
}

func TestJSONOutputFormat(t *testing.T) {
	results := []string{"abc", "def"}
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent failed: %v", err)
	}

	got := string(data)
	for _, s := range results {
		if !strings.Contains(got, s) {
			t.Errorf("JSON output missing %q: %s", s, got)
		}
	}
}
