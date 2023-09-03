import { StackProps, Stage } from "aws-cdk-lib";
import { Construct } from "constructs";
import { StatefulStack } from "./stateful-stack";
import { ApiStack } from "./api-stack";

export class CoreStage extends Stage {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const statefulStack = new StatefulStack(this, "StatefulStack");

    new ApiStack(this, "ApiStack", {
      DDB: statefulStack.DynamoDB,
    });
  }
}
