import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { RetentionDays } from "aws-cdk-lib/aws-logs";
import * as apigw from "aws-cdk-lib/aws-apigateway";
import { getFunctionBuildDir } from "../utils/lambda";
import { DynamoDB } from "./constructs/Dynamodb";

type PpApiStack = StackProps & {
  DDB: DynamoDB;
};

export class ApiStack extends Stack {
  constructor(scope: Construct, id: string, props: PpApiStack) {
    super(scope, id, props);

    const { DDB } = props;

    const whLambda = new lambda.Function(this, "bot-webhook", {
      code: lambda.Code.fromAsset(getFunctionBuildDir("webhook")),
      runtime: lambda.Runtime.PROVIDED_AL2,
      handler: "bootstrap",
      logRetention: RetentionDays.ONE_WEEK,
    });

    DDB.users.grantReadWriteData(whLambda);

    const api = new apigw.RestApi(this, "api", {
      description: "bot api gateway",
      deployOptions: {
        stageName: "dev",
      },
      // 👇 enable CORS
      defaultCorsPreflightOptions: {
        allowHeaders: [
          "Content-Type",
          "X-Amz-Date",
          "Authorization",
          "X-Api-Key",
        ],
        allowMethods: ["OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"],
        allowCredentials: true,
        allowOrigins: ["http://localhost:3000"],
      },
    });

    const whResource = api.root.addResource("webhook");
    const whIntegration = new apigw.LambdaIntegration(whLambda);
    whResource.addMethod("POST", whIntegration);
  }
}
