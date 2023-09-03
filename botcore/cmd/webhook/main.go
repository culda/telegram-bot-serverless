package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/telegram-serverless/botcore/pkg/nlambda"
	"github.com/telegram-serverless/botcore/pkg/types"
	"go.uber.org/zap"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

var svc *nlambda.Services

func init() {
	var svcErr error
	if svc, svcErr = nlambda.Instance(); svcErr != nil {
		panic(fmt.Sprintf("unable to create nlambda instance: %v", svcErr.Error()))
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var update tg.Update
	err := json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		svc.Log.Error("Error unmarshalling request body", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Failed to parse request body",
		}, err
	}

	svc.Log.Info("Received update", zap.Any("update", update))

	res := types.SendMessage{
		Method: "sendMessage",
		ChatId: update.Message.Chat.ID,
		Text:   "Hi there, I'm a bot!",
	}

	body, err := json.Marshal(res)
	if err != nil {
		svc.Log.Error("Error marshalling response", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to generate response",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(nlambda.LambdaWrap(svc, handler))
}
