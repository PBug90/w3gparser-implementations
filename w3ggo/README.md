# w3ggo

A Warcraft III replay (`.w3g`) parser written in Go. Parses replay files into structured data: players, heroes, actions, resource transfers, chat, APM, map info, and more.

Output shape mirrors [w3gjs](https://github.com/PBug90/w3gjs) for cross-implementation compatibility.

## Install

```sh
go get github.com/PBug90/w3gparser-implementations/w3ggo
```

## Usage

### Parse from file path

```go
package main

import (
	"fmt"
	"log"

	w3g "github.com/PBug90/w3gparser-implementations/w3ggo"
)

func main() {
	out, err := w3g.ParseFile("replay.w3g")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Map:     ", out.Map.File)
	fmt.Println("Matchup: ", out.Matchup)
	fmt.Printf("Duration: %ds\n", out.Duration/1000)

	for _, p := range out.Players {
		fmt.Printf("  %s (%s) — APM %d\n", p.Name, p.Race, p.APM)
	}
}
```

### Parse from bytes

```go
data, _ := os.ReadFile("replay.w3g")
out := w3g.ParseBytes(data)
```

### Serialize to JSON

`ParserOutput` is JSON-tagged, so standard `encoding/json` works:

```go
enc := json.NewEncoder(os.Stdout)
enc.SetIndent("", "  ")
enc.Encode(out)
```

## CLI

A ready-made CLI is included under `cmd/parse`:

```sh
go run github.com/PBug90/w3gparser-implementations/w3ggo/cmd/parse replay.w3g
```

## License

MIT
