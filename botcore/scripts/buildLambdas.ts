import type { ExecSyncOptions } from "child_process";
import { execSync } from "child_process";
import path from "path";
import fs from "fs-extra";
import glob from "glob";

export const buildDir = path.join(__dirname, "..", "build");

function prepareLambdas(binaryPath: string, i = 0, filePaths = [""]): void {
  const curr = i + 1;
  const total = filePaths.length;
  const filename = path.basename(binaryPath);
  const buildPath = path.join(buildDir, filename);

  console.info(`${curr} of ${total} - ${filename}: creating build folder`);
  if (!fs.existsSync(buildPath)) {
    fs.mkdirSync(buildPath, { recursive: true });
  }
  fs.copyFileSync(binaryPath, path.join(buildPath, "bootstrap"));
}

type TpFnNames = "./..." | string[];
export function buildLambdas(cmdFolderNames: TpFnNames): void {
  const start = Date.now();
  const binDir = path.join(__dirname, "..", "bin");
  const isAll = cmdFolderNames === "./...";

  // Create bin folder
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  const options: ExecSyncOptions = {
    env: {
      ...process.env,
      CGO_ENABLED: "0",
      GOARCH: "amd64",
      GOOS: "linux",
    },
    stdio: "inherit",
  };

  console.info(`go build for ${JSON.stringify(cmdFolderNames)}`);

  /**
   * -buildvcs=false removes the commit hash, introduced in Go 1.18 which caused the build hash to change on every commit and meant that it was always a new deployment in CI
   * -ldflags="-s -w" reduces bundle size by not including symbols and DWARF debugging info
   * -tags lambda.norpc removes rpc libraries from the golambda sdk, which are not needed as we use a custom runtime
   */
  const flags = [
    "-mod=vendor",
    "-buildvcs=false",
    '-ldflags="-s -w"',
    "-tags lambda.norpc",
    "-trimpath",
    "-o",
  ].join(" ");
  const goCmd = `go build ${flags}`;

  if (isAll) {
    execSync(`${goCmd} ${binDir} ${cmdFolderNames}`, {
      cwd: path.join(__dirname, ".."),
      ...options,
    });

    // Iterate over each build file and create a zip archive from it
    const filePaths = glob.sync(`${binDir}/*`, { nodir: true });
    filePaths.forEach(prepareLambdas);
  } else {
    const cmdFolderPaths = cmdFolderNames.map((cmdFolderName) => {
      execSync(`${goCmd} ${binDir}`, {
        cwd: path.join(__dirname, "..", "cmd", cmdFolderName),
        ...options,
      });
      return path.join(binDir, cmdFolderName);
    });
    cmdFolderPaths.forEach(prepareLambdas);
  }

  const duration = Date.now() - start;
  console.info(`took ${(duration / 1000).toFixed(2)}s`);
}
