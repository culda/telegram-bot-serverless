import { exec, execSync } from "child_process";
import inquirer from "inquirer";

type Stack = "ApiStack" | "StatefulStack";
type Operation = "deploy" | "destroy";
type Answers = {
  operation: Operation;
  stacks: Stack[];
};

function getStack(s: string) {
  return `CoreStage/${s}`;
}

function getCmd(op: Operation, stack: string) {
  return `cdk ${op} ${getStack(stack)} --require-approval never`;
}

async function ops() {
  const responses = await inquirer.prompt<Answers>([
    {
      type: "list",
      name: "operation",
      message: "Operation:",
      choices: ["deploy", "destroy"],
    },
    {
      type: "checkbox",
      name: "stacks",
      message: "Stacks:",
      choices: ["ApiStack", "StatefulStack"],
    },
  ]);

  for (const stack of responses.stacks) {
    const cmd = getCmd(responses.operation, stack);
    console.log(`Executing: ${cmd}...`);
    execSync(cmd, { stdio: "inherit" });
  }
}

void ops();
