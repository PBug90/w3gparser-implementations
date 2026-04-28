package main

import (
	"encoding/json"
	"fmt"
	"os"

	w3g "github.com/PBug90/w3gparser-implementations/w3ggo"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: parse <replay.w3g>")
		os.Exit(1)
	}
	result, err := w3g.ParseFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, "error encoding JSON:", err)
		os.Exit(1)
	}
}
