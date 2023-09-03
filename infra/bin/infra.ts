#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { CoreStage } from "../lib/core-stage";
import { runCommand } from "tools/node/runCommand";

async function init(): Promise<void> {
  const app = new cdk.App();

  runCommand({
    command: "yarn",
    args: ["botcore", "build"],
  });

  new CoreStage(app, "CoreStage");

  app.synth();
}

void init();
