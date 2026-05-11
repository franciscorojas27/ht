#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const { spawnSync } = require("child_process");

const exe = process.platform === "win32" ? "ht.exe" : "ht";
const binPath = path.join(__dirname, exe);

if (!fs.existsSync(binPath)) {
  console.error("Error: Binario no encontrado.");
  console.error("Intenta reinstalar: npm install -g @franciscorojas27/ht");
  process.exit(1);
}

const result = spawnSync(binPath, process.argv.slice(2), { stdio: "inherit" });
process.exit(result.status ?? 0);