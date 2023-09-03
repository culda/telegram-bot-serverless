import path from "path";
import { getRoot } from "tools/node/getRoot";

export function getFunctionBuildDir(functionName: string): string {
  return path.join(getRoot(), "botcore", "build", functionName);
}
