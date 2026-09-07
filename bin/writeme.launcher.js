#!/usr/bin/env node
"use strict";

// Thin launcher: resolves the platform-specific binary that npm installed as an
// optional dependency and hands off argv to it.
const { spawnSync } = require("node:child_process");

const PACKAGES = {
  "darwin-arm64": "writeme-cli-darwin-arm64",
  "darwin-x64": "writeme-cli-darwin-x64",
  "linux-arm64": "writeme-cli-linux-arm64",
  "linux-x64": "writeme-cli-linux-x64",
  "win32-arm64": "writeme-cli-win32-arm64",
  "win32-x64": "writeme-cli-win32-x64",
};

const platform = `${process.platform}-${process.arch}`;
const pkg = PACKAGES[platform];

if (!pkg) {
  console.error(
    `writeme: unsupported platform ${platform}\n` +
      `Build from source instead: go install github.com/dltldn333/WRITEME/cmd/writeme@latest`
  );
  process.exit(1);
}

const exe = process.platform === "win32" ? "writeme.exe" : "writeme";

let binary;
try {
  binary = require.resolve(`${pkg}/bin/${exe}`);
} catch {
  console.error(
    `writeme: the native binary for ${platform} is missing.\n` +
      `The optional dependency "${pkg}" failed to install.\n` +
      `Try: npm install writeme-cli --force`
  );
  process.exit(1);
}

const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit" });

if (result.error) {
  console.error(`writeme: ${result.error.message}`);
  process.exit(1);
}
process.exit(result.status ?? 1);
