import { spawnSync } from "child_process";
import { existsSync } from "fs";
import { resolve, dirname } from "path";
import { getRoot } from "./getRoot";

export async function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

export async function randomSleep(max = 500): Promise<void> {
  const randomTime = Math.floor(Math.random() * Math.floor(max));
  return sleep(randomTime);
}

export function trueKeysFromBooleanMap(map: Record<string, boolean>): string[] {
  const filtered = Object.entries(map).filter(([_, value]) => Boolean(value));

  return Array.from(filtered, (element) => element[0]);
}

type PpRunCommand = {
  args?: string[];
  command: string;
  cwd?: string;
};

type TpSpawnSyncResult = {
  status: number;
};

export function runCommand({ command, cwd, args = [] }: PpRunCommand): void {
  const fullCmd = `${command} ${args.join(" ")}`;
  console.info("Executing command", fullCmd);

  const childProcess = spawnSync(command, args, {
    cwd: getRoot(),
    stdio: "inherit",
  });

  const result: TpSpawnSyncResult = {
    status: childProcess.status ?? 1,
  };

  if (result.status !== 0) {
    console.error(
      `Command "${fullCmd}" failed with exit code ${result.status}`
    );
  }

  function handleSignal(signal: NodeJS.Signals): void {
    if (childProcess.pid) {
      console.error(
        `Command "${fullCmd}" received ${signal}. Killing child process ${childProcess.pid}...`
      );
      process.kill(childProcess.pid, signal);
    }
  }

  process.on("SIGTERM", () => {
    handleSignal("SIGTERM");
    process.exit();
  });

  process.on("SIGINT", () => {
    handleSignal("SIGINT");
    process.exit();
  });
}
