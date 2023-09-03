# Telegram Bot with Serverless Architecture on AWS

This project is a Telegram bot that runs on a serverless architecture, hosted on AWS. It uses AWS API Gateway to receive webhook requests from Telegram and forwards them to an AWS Lambda function. The Lambda function is integrated with Amazon DynamoDB for database operations. The entire infrastructure is provisioned using the AWS Cloud Development Kit (CDK).

## Features

- **Serverless**: Built on top of AWS Lambda for auto-scaling and pay-as-you-go pricing.
- **Webhook Based**: Utilizes Telegram's webhook API for real-time updates.
- **AWS API Gateway**: Serves as the entry point for all webhook requests.
- **DynamoDB**: For stateful interactions and data persistence.
- **Infrastructure as Code**: Uses AWS CDK to make deployment and management easier.

## Requirements

- AWS Account
- AWS CLI configured with necessary access rights
- Node.js (optional, needed for CDK)
- Python 3.6+ (for the Lambda function)
- AWS CDK Toolkit (if you're using CDK)

## Quick Start

### Deployment

Make sure your AWS account is linked in the terminal session.

```bash
aws configure
```

Run the following command to deploy all stacks via AWS CDK.

```bash
make ops
```

Set Up Telegram Webhook
After deploying the stack, you'll receive an API Gateway URL. You need to set this URL as your Telegram bot's webhook.

To set the webhook, make an API call to Telegram using the setWebhook method. Replace <API_TOKEN> with your Telegram bot token and <API_GATEWAY_URL> with the URL obtained from the stack deployment.

```bash
curl -F "url=<API_GATEWAY_URL>" https://api.telegram.org/bot<API_TOKEN>/setWebhook
```

Once the webhook is set, your Telegram bot should start receiving updates, which will be processed by the AWS Lambda function.

### License

This project is open-source and available under the MIT License.

### Support and Feedback

For support, issues, or feedback, please open an issue on this repository.
