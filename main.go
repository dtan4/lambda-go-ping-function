package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context) error {
	fmt.Println("ping function")

	return nil
}

func main() {
	lambda.Start(Handler)
}
