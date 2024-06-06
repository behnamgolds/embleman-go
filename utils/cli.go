package utils

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

const appname = "embleman"
const version = "0.1.0"

type action struct {
	sName string
	help  string
}

var actionStrs = []string{"--increase", "--decrease", "--check", "--clock", "--clear", "--version", "--help", "-v", "-h"}
var actions = map[string]action{
	"--increase": {
		"inc", "Increase emblem number:\n  --increase <path>",
	},
	"--decrease": {
		"dec", "Decrease emblem number:\n  --decrease <path>",
	},
	"--check": {
		"check", "Toggle check emblem:\n  --check <path>",
	},
	"--clock": {
		"clock", "Toggle clock emblem:\n  --clock <path>",
	},
	"--clear": {
		"clear", "Clear all emblems:\n  --clear <path(s)>",
	},
	"--version": {
		"version", "Show application version",
	},
	"--help": {
		"help", "Show this help",
	},
}

func ParseCmdArgs() []string {
	args := os.Args[1:]

	if len(args) == 0 || !slices.Contains(actionStrs, args[0]) {
		fatalParseArgs("No valid flag provided")
	}

	if args[0] == "--version" || args[0] == "-v" {
		printVersion()
		os.Exit(0)
	}

	if args[0] == "--help" || args[0] == "-h" {
		printHelp()
		os.Exit(0)
	}

	if len(args) == 1 {
		hint := fmt.Sprintf("Flag %s needs argument", args[0])
		fatalParseArgs(hint)
	}

	if args[0] != "--clear" && len(args) > 2 {
		fatalParseArgs("Only --clear can take multiple arguments")
	}

	for _, arg := range args[1:] {
		// check if no other argument has "-" or "--" prefix
		if strings.HasPrefix(arg, "-") {
			hint := fmt.Sprintf("Invalid argument: %s", arg)
			fatalParseArgs(hint)
		}

		if _, err := os.Stat(arg); err != nil {
			hint := fmt.Sprintf("Path does not exist: %s", arg)
			fatalParseArgs(hint)
		}
	}
	args[0] = actions[args[0]].sName
	return args
}

func fatalParseArgs(hint string) {
	printHelp()
	if hint != "" {
		fmt.Printf("\nHint: %s", hint)
	}
	fmt.Print("\n\n")
	os.Exit(2)
}

func printHelp() {
	fmt.Printf("%s: a set of emblem actions\n\n", appname)
	for key, value := range actions {
		fmt.Printf("%s: %s\n\n", key, value.help)
	}
}

func printVersion() {
	fmt.Printf("%s version %s\n", appname, version)
}
