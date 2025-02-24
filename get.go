package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/thisdougb/doocot/config"
)

const getUsage = `Get some data from doocot.
Usage:
  doocot get [options] <identifier>

Options:
  -raw        get the raw data from the backend

Example:
  $ doocot get slight-step-zoo-flock
  this is my secret text
`

func get(ctx context.Context, args []string) {

	getDataUrl, _ := url.JoinPath(config.StringValue("DOOCOT_HOST"), "/api/data")

	fs := flag.NewFlagSet("get", flag.ExitOnError)
	fs.Usage = func() { fmt.Print(getUsage) }

	raw := fs.Bool("raw", false, "return the raw data as stored remotely")
	fs.Parse(args)

	if *raw {
		getDataUrl, _ = url.JoinPath(getDataUrl, "/raw")
	}

	config.LogDebug(ctx,
		fmt.Sprintf("Using backend %s", config.StringValue("DOOCOT_HOST")))

	// really, there should be only one final-position argument
	if len(fs.Args()) != 1 {
		fs.Usage()
		os.Exit(1)
	}

	// leave it to the backend to validate the supplied identifier
	getDataUrl, _ = url.JoinPath(getDataUrl, fs.Args()[0])

	// create a context with reasonable timeout
	httpCtx, cncl := context.WithTimeout(ctx, time.Second*3)
	defer cncl()

	req, err := http.NewRequestWithContext(httpCtx, http.MethodGet, getDataUrl, nil)
	if err != nil {
		config.LogError(ctx, "Failed to create new http request.")
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		config.LogError(ctx, fmt.Sprintf("Failed making api request: %s", err.Error()))
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		config.LogError(ctx,
			fmt.Sprintf("Failure response code from backend (%s): %d %s",
				config.StringValue("DOOCOT_HOST"),
				resp.StatusCode,
				http.StatusText(resp.StatusCode)))
		os.Exit(1)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		config.LogError(ctx, fmt.Sprintf("Failed to make api request: %s", err.Error()))
		os.Exit(1)
	}

	fmt.Print(string(bodyBytes))
}
