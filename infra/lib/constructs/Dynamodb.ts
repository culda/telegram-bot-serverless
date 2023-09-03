import * as constructs from "constructs";
import * as cdk from "aws-cdk-lib";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";

export class DynamoDB extends constructs.Construct {
  public readonly users: dynamodb.Table;

  constructor(scope: cdk.Stack, id: string) {
    super(scope, id);

    this.users = new dynamodb.Table(this, "Users", {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      partitionKey: { name: "ID", type: dynamodb.AttributeType.STRING },
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });
  }
}
