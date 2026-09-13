#!/usr/bin/env node
"use strict";

// Runs the Go binary built for this platform, shipped inside the package under dist/.
const { spawnSync } = require("node:child_process");
const { existsSync } = require("node:fs");
const path = require("node:path");

const SUPPORTED = new Set([
  "darwin-arm64",
  "darwin-x64",
  "linux-arm64",
  "linux-x64",
  "win32-arm64",
  "win32-x64",
]);

const platform = `${process.platform}-${process.arch}`;

if (!SUPPORTED.has(platform)) {
  console.error(
    `writeme: unsupported platform ${platform}\n` +
      `Install from source instead: go install github.com/dltldn333/WRITEME/cmd/writeme@latest`
  );
  process.exit(1);
}

const exe = process.platform === "win32" ? "writeme.exe" : "writeme";
const binary = path.join(__dirname, "..", "dist", platform, exe);

if (!existsSync(binary)) {
  console.error(`writeme: missing binary ${binary}\nReinstall writeme-cli.`);
  process.exit(1);
}

const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit" });

if (result.error) {
  console.error(`writeme: ${result.error.message}`);
  process.exit(1);
}
if (result.signal) {
  process.kill(process.pid, result.signal);
}
process.exit(result.status ?? 1);
