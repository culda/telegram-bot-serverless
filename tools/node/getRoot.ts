import { existsSync } from "fs";
import { dirname, resolve } from "path";

export function getRoot(startDir: string = process.cwd()): string | never {
  let currentDir = startDir;

  while (currentDir !== "/") {
    if (existsSync(resolve(currentDir, ".yarnrc.yml"))) {
      console.log("Found root directory", currentDir);
      return currentDir;
    }

    currentDir = dirname(currentDir);
  }

  throw new Error("root directory not found");
}
