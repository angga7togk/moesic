#!/usr/bin/env node

import os from "node:os";
import fs from "node:fs";
import { spawn } from "child_process";
import path from "path";
import { pipeline } from "stream";
import { promisify } from "util";
import { Readable } from "stream";

const streamPipeline = promisify(pipeline);

const homePath = os.homedir();
const moesicPath = path.join(homePath, ".moesic");
const binPath = path.join(moesicPath, "bin");
const platform = process.platform;
const arch = process.arch;

async function installMoesic(force = false) {
  const moesicFile = path.join(moesicPath, "moesic");

  if (fs.existsSync(moesicFile) && !force) {
    return;
  }

  console.log("Installing Moesic...");
  fs.mkdirSync(moesicPath, { recursive: true });
  console.log("'~/.moesic/' folder created.");
  console.log(`Your platform: ${platform}/${arch}`);

  let filename;
  if (platform === "linux") {
    filename = "moesic-linux";
  } else if (platform === "darwin") {
    filename = "moesic-macos";
  } else if (platform === "win32") {
    filename = "moesic-windows.exe";
  } else {
    console.error(`Unsupported Platform: ${platform}.`);
    process.exit(1);
  }

  console.log("Fetching latest release info...");
  let latestTag;
  try {
    const releaseInfo = await fetch(
      "https://api.github.com/repos/angga7togk/moesic/releases/latest"
    ).then((r) => r.json());
    latestTag = releaseInfo.tag_name;
  } catch (error) {
    console.error("Failed to get latest release.");
    process.exit(1);
  }

  console.log("Downloading moesic...");
  try {
    const res = await fetch(
      `https://github.com/angga7togk/moesic/releases/download/${latestTag}/${filename}`
    );
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);

    await streamPipeline(
      Readable.fromWeb(res.body),
      fs.createWriteStream(moesicFile)
    );

    if (platform !== "win32") {
      fs.chmodSync(moesicFile, 0o755);
    }

    console.log("Moesic downloaded.");
  } catch (error) {
    console.error(error);
    console.error("Failed to download moesic.");
    process.exit(1);
  }
  console.log("Moesic installed.");
}

async function installYtDlp(force = false) {
  const targetPath = path.join(binPath, "yt-dlp");

  if (fs.existsSync(targetPath) && !force) {
    return;
  }

  console.log("Installing yt-dlp...");
  fs.mkdirSync(binPath, { recursive: true });
  console.log("'~/.moesic/bin' folder created.");

  console.log(`Your platform: ${platform}/${arch}`);

  const baseYtDlpUrl =
    "https://github.com/yt-dlp/yt-dlp/releases/latest/download/";
  let ytDlpUrl;
  switch (platform) {
    case "win32":
      ytDlpUrl = baseYtDlpUrl + "yt-dlp.exe";
      break;
    case "darwin":
      ytDlpUrl = baseYtDlpUrl + "yt-dlp_macos";
      break;
    case "linux":
      ytDlpUrl = baseYtDlpUrl + "yt-dlp_linux";
      break;
    default:
      console.error(`Unsupported Platform: ${platform}.`);
      process.exit(1);
  }

  console.log(`Downloading from ${ytDlpUrl}...`);
  try {
    const res = await fetch(ytDlpUrl);
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);

    await streamPipeline(res.body, fs.createWriteStream(targetPath));

    if (platform !== "win32") {
      fs.chmodSync(targetPath, 0o755);
    }

    console.log("yt-dlp downloaded and installed.");
  } catch (error) {
    console.error("Failed to download yt-dlp:", error.message);
    process.exit(1);
  }
}

await installMoesic();
await installYtDlp();

const args = process.argv.slice(2);
spawn(path.join(moesicPath, "moesic"), args, {
  cwd: moesicPath,
  stdio: "inherit",
});
