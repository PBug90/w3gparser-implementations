#!/usr/bin/env tsx

import assert from 'assert/strict';
import { execFileSync } from 'child_process';
import path from 'path';
import fs from 'fs';
import { fileURLToPath } from 'url';
import type { ParserOutput as W3gjsOutput } from 'w3gjs';

const HERE        = path.dirname(fileURLToPath(import.meta.url));
const ROOT        = path.resolve(HERE, '..');
const REPLAYS_DIR = path.join(HERE, 'replays');
const OUTPUT_DIR  = path.join(HERE, 'output');

const EXE = process.platform === 'win32' ? '.exe' : '';
const GO_BIN   = process.env.GO_BIN   ?? path.join(ROOT, 'bin', `w3g-go${EXE}`);
const RUST_BIN = process.env.RUST_BIN ?? path.join(ROOT, 'w3grs', 'target', 'release', `parse${EXE}`);

const replays: string[] = process.argv.slice(2).length
  ? process.argv.slice(2).map(r => path.resolve(r))
  : fs.readdirSync(REPLAYS_DIR)
      .filter(f => f.endsWith('.w3g') || f.endsWith('.nwg'))
      .sort()
      .map(f => path.join(REPLAYS_DIR, f));

// ---------------------------------------------------------------------------
// Parser output type — camelCase, matching w3gjs shape, plus currentTimePlayed
// ---------------------------------------------------------------------------

type ParserOutput = W3gjsOutput & {
  parseTime: number;
  players: Array<W3gjsOutput['players'][number] & { currentTimePlayed: number }>;
};

// ---------------------------------------------------------------------------
// Normalizers
// ---------------------------------------------------------------------------

// Strip fields our parsers emit that w3gjs doesn't, for cross-tool comparison.
function normParser(p: ParserOutput): W3gjsOutput {
  const { parseTime: _, players, ...rest } = p;
  return {
    ...rest,
    players: players.map(({ currentTimePlayed: __, ...pl }) => pl),
  } as W3gjsOutput;
}

// JSON round-trip triggers Player.toJSON(), stripping w3gjs internals.
// Also fixes a known w3gjs typo in chat mode ("Obervers" → "Observers").
function normW3gjs(w: W3gjsOutput): Omit<W3gjsOutput, 'parseTime'> {
  const { parseTime: _, ...j } = JSON.parse(JSON.stringify(w)) as W3gjsOutput;
  if (j.chat?.length) {
    j.chat = j.chat.map(m => ({
      ...m,
      // w3gjs has a known typo: "Obervers" instead of "Observers"
      mode: m.mode === 'Obervers' ? 'Observers' : m.mode,
    }));
  }
  return j;
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function parse(bin: string, replayPath: string): ParserOutput {
  const stdout = execFileSync(bin, [replayPath], { encoding: 'utf8', timeout: 30_000 });
  return JSON.parse(stdout) as ParserOutput;
}

function stripParseTime(p: ParserOutput): Omit<ParserOutput, 'parseTime'> {
  const { parseTime: _, ...rest } = p;
  return rest;
}

function saveOutput(stem: string, tool: string, data: unknown): void {
  const dir = path.join(OUTPUT_DIR, stem);
  fs.mkdirSync(dir, { recursive: true });
  fs.writeFileSync(path.join(dir, `${tool}.json`), JSON.stringify(data, null, 2));
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

async function main(): Promise<void> {
  for (const replay of replays) {
    const stem = path.basename(replay, path.extname(replay));
    console.log(`\n${path.relative(HERE, replay)}`);

    const goOut   = parse(GO_BIN,   replay);
    const rustOut = parse(RUST_BIN, replay);
    saveOutput(stem, 'go',   goOut);
    saveOutput(stem, 'rust', rustOut);

    process.stdout.write('  Go == Rust ... ');
    assert.deepStrictEqual(stripParseTime(rustOut), stripParseTime(goOut));
    console.log('ok');

    process.stdout.write('  w3gjs x-check ... ');
    let W3GReplay: typeof import('w3gjs').default;
    try {
      ({ default: W3GReplay } = await import('w3gjs'));
    } catch {
      console.log('skipped (w3gjs not installed)');
      continue;
    }
    const w3gjsOut = await new W3GReplay().parse(replay);
    saveOutput(stem, 'w3gjs', JSON.parse(JSON.stringify(w3gjsOut)));
    assert.deepStrictEqual(normW3gjs(w3gjsOut), normParser(goOut));
    console.log('ok');
  }

  console.log(`\n${replays.length}/${replays.length} passed`);
}

main().catch(err => { console.error(err.message); process.exit(1); });
