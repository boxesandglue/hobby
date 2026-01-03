// Command hobby runs Lua scripts with the hobby library.
package main

import (
	"fmt"
	"os"

	"github.com/boxesandglue/hobby"
	lua "github.com/speedata/go-lua"
)

// Version is set via ldflags at build time
var Version = "dev"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: hobby <script.lua>\n")
		os.Exit(1)
	}

	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Printf("hobby %s\n", Version)
		return
	}

	l := lua.NewState()

	lua.OpenLibraries(l)
	hobby.Open(l)

	if err := lua.DoFile(l, os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
