// Cross-compiles the Go CLI into dist/<platform>-<arch>/ for the npm package.
// `node scripts/build.mjs --check` only verifies every binary exists.
import { execFileSync } from "node:child_process";
import { existsSync, rmSync } from "node:fs";
import { join } from "node:path";

const TARGETS = [
  { dir: "darwin-arm64", goos: "darwin", goarch: "arm64" },
  { dir: "darwin-x64", goos: "darwin", goarch: "amd64" },
  { dir: "linux-arm64", goos: "linux", goarch: "arm64" },
  { dir: "linux-x64", goos: "linux", goarch: "amd64" },
  { dir: "win32-arm64", goos: "windows", goarch: "arm64" },
  { dir: "win32-x64", goos: "windows", goarch: "amd64" },
];

const binaryPath = (t) => join("dist", t.dir, t.goos === "windows" ? "writeme.exe" : "writeme");

if (process.argv.includes("--check")) {
  const missing = TARGETS.map(binaryPath).filter((p) => !existsSync(p));
  if (missing.length > 0) {
    console.error(`Missing binaries, run "npm run build" first:\n  ${missing.join("\n  ")}`);
    process.exit(1);
  }
  process.exit(0);
}

rmSync("dist", { recursive: true, force: true });

for (const t of TARGETS) {
  const out = binaryPath(t);
  execFileSync("go", ["build", "-trimpath", "-ldflags=-s -w", "-o", out, "./cmd/writeme"], {
    stdio: "inherit",
    env: { ...process.env, CGO_ENABLED: "0", GOOS: t.goos, GOARCH: t.goarch },
  });
  // stderr keeps stdout clean for `npm publish --json`, which changesets parses.
  console.error(`built ${out}`);
}
