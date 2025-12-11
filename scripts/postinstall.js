import fs from "fs";
import path from "path";
import { execSync } from "child_process";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const packageDir = path.join(__dirname, "..");

const platform = process.platform;
const arch = process.arch;

const archMap = {
  x64: "amd64",
  arm64: "arm64",
};

const platformMap = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const goOS = platformMap[platform];
const goArch = archMap[arch];

if (!goOS || !goArch) {
  console.error(
    `Unsupported platform: ${platform} ${arch}. Skipping binary download.`,
  );
  process.exit(0);
}

const binaryName = platform === "win32" ? "fizzy.exe" : "fizzy";
const releaseTag = process.env.FIZZY_RELEASE_TAG || "latest";

const downloadUrl = `https://github.com/rogeriopvl/fizzy-cli/releases/download/${releaseTag}/fizzy-${goOS}-${goArch}${platform === "win32" ? ".exe" : ""}`;

const binDir = path.join(packageDir, "bin");
const binaryPath = path.join(binDir, binaryName);

if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

console.log(
  `Downloading fizzy binary for ${goOS}-${goArch} from ${downloadUrl}`,
);

try {
  execSync(`curl -L -o ${binaryPath} ${downloadUrl}`, { stdio: "inherit" });
  fs.chmodSync(binaryPath, 0o755);
  console.log(`Successfully installed fizzy binary to ${binaryPath}`);
} catch (error) {
  console.error(`Failed to download fizzy binary: ${error.message}`);
  process.exit(1);
}
