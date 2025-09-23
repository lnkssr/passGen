package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
)

var (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars   = "0123456789"
	symbolChars  = "!@#$%^&*()-_=+[]{};:,.<>?/|"
	similarChars = "0O1lI5S"
)

type options struct {
	length     int
	count      int
	lower      bool
	upper      bool
	digits     bool
	symbols    bool
	all        bool
	noSimilar  bool
	custom     string
	jsonOutput bool
	help       bool
}

func main() {
	opts := parseFlags()

	if opts.help {
		printHelp()
		return
	}

	charset := buildCharset(opts)
	if charset == "" {
		fmt.Println("Error: no character set selected (use -h for help)")
		os.Exit(1)
	}

	var results []string
	for i := 0; i < opts.count; i++ {
		pass := generatePassword(opts.length, charset, opts)
		results = append(results, pass)
	}

	if opts.jsonOutput {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
	} else {
		for _, p := range results {
			fmt.Println(p)
		}
	}
}

func parseFlags() options {
	var opts options

	flag.IntVar(&opts.length, "l", 12, "length of password")
	flag.IntVar(&opts.count, "n", 1, "number of passwords to generate")
	flag.BoolVar(&opts.lower, "lower", false, "include lowercase letters")
	flag.BoolVar(&opts.upper, "upper", false, "include uppercase letters")
	flag.BoolVar(&opts.digits, "digits", false, "include digits")
	flag.BoolVar(&opts.symbols, "symbols", false, "include symbols")
	flag.BoolVar(&opts.all, "all", false, "use all available characters")
	flag.BoolVar(&opts.noSimilar, "no-similar", false, "exclude similar characters (0/O, 1/l/I, 5/S)")
	flag.StringVar(&opts.custom, "g", "", "custom range of symbols (e.g. A-F,0-5)")
	flag.BoolVar(&opts.jsonOutput, "json", false, "output in JSON format")
	flag.BoolVar(&opts.help, "h", false, "print help message")

	flag.Parse()
	return opts
}

func buildCharset(opts options) string {
	var charset strings.Builder

	if opts.all {
		charset.WriteString(lowerChars + upperChars + digitChars + symbolChars)
	} else {
		if opts.lower {
			charset.WriteString(lowerChars)
		}
		if opts.upper {
			charset.WriteString(upperChars)
		}
		if opts.digits {
			charset.WriteString(digitChars)
		}
		if opts.symbols {
			charset.WriteString(symbolChars)
		}
		if opts.custom != "" {
			charset.WriteString(parseRange(opts.custom))
		}
	}

	result := charset.String()

	if opts.noSimilar {
		for _, s := range similarChars {
			result = strings.ReplaceAll(result, string(s), "")
		}
	}

	return result
}

func generatePassword(length int, charset string, opts options) string {
	if len(charset) == 0 {
		return ""
	}

	var sb strings.Builder
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		sb.WriteByte(charset[n.Int64()])
	}
	return sb.String()
}

func parseRange(r string) string {
	var chars []rune
	parts := strings.Split(r, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")
			if len(bounds) == 2 {
				start := rune(bounds[0][0])
				end := rune(bounds[1][0])
				for c := start; c <= end; c++ {
					chars = append(chars, c)
				}
			}
		} else {
			chars = append(chars, []rune(part)...)
		}
	}
	return string(chars)
}

func printHelp() {
	fmt.Println(`Usage: passGen [FLAGS]... [OPTIONS]...

Flags:
  -l <num>          length of password (default 12)
  -n <num>          number of passwords to generate (default 1)
  --lower           include lowercase letters (a-z)
  --upper           include uppercase letters (A-Z)
  --digits          include digits (0-9)
  --symbols         include symbols (!@#$...)
  --all             use all available characters
  --no-similar      exclude similar characters (0/O, 1/l/I, 5/S)
  -g <range>        custom range of symbols, e.g. "A-F,0-5"
  --json            output as JSON array
  -h                print this help message

Examples:
  passGen -l 16 --all
  passGen -n 5 --lower --digits
  passGen -g "A-Z,0-9" -l 8
`)
}
