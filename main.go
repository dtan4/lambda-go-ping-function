package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/pkg/errors"
)

const (
	defaultPingURL = "https://checkip.amazonaws.com"
)

var (
	HTTPClient *http.Client
	PingURL    string
	Stdout     io.Writer
)

func RealHandler(ctx context.Context) error {
	req, err := http.NewRequest(http.MethodGet, PingURL, nil)
	if err != nil {
		return errors.Wrap(err, "cannot create HTTP request")
	}

	req.WithContext(ctx)

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "cannot access to the URL: %s", PingURL)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected response status: %d, body: %q", resp.StatusCode, body)
	}

	fmt.Fprintln(Stdout, body)

	return nil
}

func Handler(ctx context.Context) error {
	HTTPClient = xray.Client(http.DefaultClient)
	PingURL = defaultPingURL
	Stdout = os.Stdout

	return RealHandler(ctx)
}

func main() {
	lambda.Start(Handler)
}
