package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/app"
)

var ginLambda *ginadapter.GinLambdaV2

func init() {
	application, err := app.Initialize()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	ginLambda = ginadapter.NewV2(application.Router)
}

// Handler proxies API Gateway v2 (HTTP API) events to the Gin router.
func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
