import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import { DynamoDB } from "./constructs/Dynamodb";

export class StatefulStack extends Stack {
  public readonly DynamoDB: DynamoDB;

  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    this.DynamoDB = new DynamoDB(this, "DynamoDB");
  }
}
