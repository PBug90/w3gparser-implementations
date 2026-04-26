# w3gparser-implementations

Parallel Go and Rust implementations of a Warcraft III replay (`.w3g`) parser, cross-validated against [w3gjs](https://github.com/pbug90/w3gjs).

Both parsers were written entirely by [Claude](https://claude.ai) (Anthropic). Both produce identical JSON output and are verified against each other and against w3gjs on every push via GitHub Actions.

---

## Repository layout

```
w3ggo/                  Go parser (module: w3ggo)
w3grs/                  Rust parser (crate: w3grs)
validate/               Cross-tool validation harness (Node/TypeScript)
  replays/              Sample .w3g files used for validation
  output/               Per-replay JSON dumps (gitignored, for local inspection)
.github/workflows/      CI pipeline
```

---

## Parsers

Both parsers expose a CLI binary that accepts a single `.w3g` path and writes JSON to stdout.

### Go

```bash
cd w3ggo
go test ./...                         # unit tests
go build -o ../bin/w3g-go ./cmd/parse # build CLI
../bin/w3g-go path/to/replay.w3g
```

Requires Go 1.21+.

### Rust

```bash
cd w3grs
cargo test --all                      # unit tests
cargo build --release                 # build CLI → target/release/parse
./target/release/parse path/to/replay.w3g
```

Requires a stable Rust toolchain.

---

## JSON output format

Both parsers emit camelCase JSON matching the w3gjs output shape. The top-level structure is:

```jsonc
{
  "id": "...",
  "gamename": "...",
  "players": [ { "id": 1, "name": "...", "apm": 142, "heroes": [...], ... } ],
  "matchup": "HvO",
  "duration": 1800000,
  "settings": { "observerMode": "NONE", "fixedTeams": true, ... },
  "chat": [...],
  "map": { "path": "...", "file": "...", "checksum": "...", "checksumSha1": "..." },
  "parseTime": 12,
  ...
}
```

The only fields present in our output but absent from w3gjs are `parseTime` (top-level) and `currentTimePlayed` (per player).

---

## Cross-tool validation

The `validate/` harness runs all three parsers against every replay in `validate/replays/` and asserts:

1. **Go == Rust** — byte-for-byte identical JSON (excluding `parseTime`)
2. **w3gjs x-check** — our output matches w3gjs after stripping timing fields and a [known w3gjs typo](https://github.com/angehung/w3gjs) (`"Obervers"` → `"Observers"` in chat mode)

```bash
cd validate
npm install
npm run validate                      # all replays in validate/replays/
npx tsx validate.ts path/to/any.w3g  # single replay
```

JSON output for each replay is written to `validate/output/<name>/{go,rust,w3gjs}.json` for manual inspection.

The validation script locates the parser binaries at `bin/w3g-go` and `w3grs/target/release/parse` by default; override with the `GO_BIN` and `RUST_BIN` environment variables.

---

## CI

Three jobs run on every push and pull request:

| Job | What it does |
|-----|-------------|
| `test-go` | `go test ./...` in `w3ggo/` |
| `test-rust` | `cargo test --all` in `w3grs/` |
| `cross-validate` | Builds both CLIs, installs w3gjs, runs `validate.ts` against all sample replays |

`cross-validate` only runs after both unit-test jobs pass.
