package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
	"unicode/utf8"
)

var (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars   = "0123456789"
	symbolChars  = "!@#$%^&*()-_=+[]{};:,.<>?/|"
	similarChars = "0O1lI5S"
)

// Command describes a CLI subcommand.
type Command struct {
	Name  string
	Usage string
	Long  string
	Flag  flag.FlagSet
	Run   func(cmd *Command, args []string)
}

var commands []*Command

// Register all subcommands here
func init() {
	commands = []*Command{
		cmdGen,
		cmdCharset,
		cmdHelp,
	}
}

//
// ─── SUBCOMMAND: GEN ──────────────────────────────────────────────────────────────
//

var cmdGen = &Command{
	Name:  "gen",
	Usage: "generate passwords",
	Long: `
Generate random passwords with custom rules.

Usage:
  passgen gen [flags]

Examples:
  passgen gen -l 16 -all
  passgen gen -n 5 -lower -digits
  passgen gen -g "A-Z,0-9" -l 8
`,
}

type genOptions struct {
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
}

func init() {
	cmdGen.Flag.IntVar(&genOpts.length, "l", 12, "length of password")
	cmdGen.Flag.IntVar(&genOpts.count, "n", 1, "number of passwords")
	cmdGen.Flag.BoolVar(&genOpts.lower, "lower", false, "include lowercase letters")
	cmdGen.Flag.BoolVar(&genOpts.upper, "upper", false, "include uppercase letters")
	cmdGen.Flag.BoolVar(&genOpts.digits, "digits", false, "include digits")
	cmdGen.Flag.BoolVar(&genOpts.symbols, "symbols", false, "include symbols")
	cmdGen.Flag.BoolVar(&genOpts.all, "all", false, "use all available characters")
	cmdGen.Flag.BoolVar(&genOpts.noSimilar, "no-similar", false, "exclude similar characters")
	cmdGen.Flag.StringVar(&genOpts.custom, "g", "", "custom range (e.g. A-F,0-5)")
	cmdGen.Flag.BoolVar(&genOpts.jsonOutput, "json", false, "output JSON")
	cmdGen.Run = runGen
}

var genOpts genOptions

func runGen(cmd *Command, args []string) {
	cmd.Flag.Parse(args)

	charset := buildCharset(genOpts)
	if charset == "" {
		fmt.Println("Error: no character set selected")
		os.Exit(1)
	}

	var results []string
	for i := 0; i < genOpts.count; i++ {
		pass := generatePassword(genOpts.length, charset)
		results = append(results, pass)
	}

	if genOpts.jsonOutput {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
	} else {
		for _, p := range results {
			fmt.Println(p)
		}
	}
}

//
// ─── SUBCOMMAND: CHARSET ─────────────────────────────────────────────────────────
//

var cmdCharset = &Command{
	Name:  "charset",
	Usage: "show built-in character sets",
	Long: `
Show character sets used for password generation.

Usage:
  passgen charset
`,
	Run: func(cmd *Command, args []string) {
		fmt.Println("lower:   ", lowerChars)
		fmt.Println("upper:   ", upperChars)
		fmt.Println("digits:  ", digitChars)
		fmt.Println("symbols: ", symbolChars)
		fmt.Println("similar: ", similarChars)
	},
}

//
// ─── SUBCOMMAND: HELP ────────────────────────────────────────────────────────────
//

var cmdHelp = &Command{
	Name:  "help",
	Usage: "show help for a command",
	Long: `
Show detailed help for a specific subcommand.

Usage:
  passgen help [command]
`,
	Run: func(cmd *Command, args []string) {
		if len(args) == 0 {
			printMainHelp()
			return
		}
		name := args[0]
		for _, c := range commands {
			if c.Name == name {
				fmt.Println(strings.TrimSpace(c.Long))
				return
			}
		}
		fmt.Printf("Unknown command: %s\n", name)
	},
}

//
// ─── MAIN ENTRYPOINT ─────────────────────────────────────────────────────────────
//

func main() {
	flag.Usage = printMainHelp
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printMainHelp()
		return
	}

	cmdName := args[0]
	for _, cmd := range commands {
		if cmd.Name == cmdName {
			cmd.Run(cmd, args[1:])
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %q\n", cmdName)
	fmt.Println("Run 'passgen help' for usage.")
	os.Exit(1)
}

//
// ─── UTIL FUNCTIONS ─────────────────────────────────────────────────────────────
//

func printMainHelp() {
	fmt.Println(`passgen - password generator CLI

Usage:
  passgen <command> [flags]

Available commands:`)
	for _, c := range commands {
		fmt.Printf("  %-10s %s\n", c.Name, c.Usage)
	}
	fmt.Println("\nUse 'passgen help <command>' for more details.")
}

func buildCharset(opts genOptions) string {
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

func generatePassword(length int, charset string) string {
	if len(charset) == 0 {
		return ""
	}
	var sb strings.Builder
	for range make([]struct{}, length) {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		sb.WriteByte(charset[n.Int64()])
	}
	return sb.String()
}

func parseRange(r string) string {
	var chars []rune
	for _, part := range strings.Split(r, ",") {
		part = strings.TrimSpace(part)
		if !strings.Contains(part, "-") {
			chars = append(chars, []rune(part)...)
			continue
		}
		bounds := strings.Split(part, "-")
		if len(bounds) != 2 {
			continue
		}
		start, _ := utf8.DecodeLastRuneInString(bounds[0])
		end, _ := utf8.DecodeLastRuneInString(bounds[1])
		for j := start; j <= end; j++ {
			chars = append(chars, j)
		}
	}
	return string(chars)
}
