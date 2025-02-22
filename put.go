package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/thisdougb/doocot/config"
)

const putUsage = `Put some data into doocot.
Usage:
  doocot put [options] <data>

Options:
  -v           verbose commentary
  -words       return a word passphrase link rather than hex string
  -once        expire the data after it has been read once
  -json        output in json format

Example:
  $ doocot put -words this is my secret text
  doocot get slight-step-zoo-flock
  curl https://doocot.sh/collect/slight-step-zoo-flock
`

func put(ctx context.Context, args []string) {

	type RequestData struct {
		DataValue     string `json:"data_value"`
		Once          bool   `json:"once"`
		Words         bool   `json:"words"`
		ClientVersion string `json:"client_version"`
	}

	putDataUrl, _ := url.JoinPath(config.StringValue("DOOCOT_HOST"), "/api/data")

	fs := flag.NewFlagSet("share", flag.ExitOnError)
	fs.Usage = func() { fmt.Print(putUsage) }

	verbose := fs.Bool("v", false, "verbose")
	words := fs.Bool("words", false, "words")
	once := fs.Bool("once", false, "once")
	jsonOut := fs.Bool("json", false, "json")
	fs.Parse(args)

	if *verbose {
		ctx = config.SetContextDebug(ctx, *verbose)
	}

	config.LogDebug(ctx,
		fmt.Sprintf("Using backend %s", config.StringValue("DOOCOT_HOST")))

	if len(fs.Args()) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	// the data to send is the rest of the command line after the flags args
	value := strings.Join(fs.Args(), " ")

	var requestData RequestData

	requestData.DataValue = value
	requestData.Once = *once
	requestData.Words = *words
	requestData.ClientVersion = Version

	requestDataJsonBytes, err := json.Marshal(requestData)
	if err != nil {
		config.LogError(ctx, "Failed to create json data for api request.")
		os.Exit(1)
	}

	// create a context with reasonable timeout
	httpCtx, cncl := context.WithTimeout(ctx, time.Second*3)
	defer cncl()

	req, err := http.NewRequestWithContext(
		httpCtx, http.MethodPost, putDataUrl,
		bytes.NewBuffer(requestDataJsonBytes))

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

	if resp.StatusCode != http.StatusCreated {
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

	type RespData struct {
		Id      string `json:"id"`
		Url     string `json:"url"`
		Expires string `json:"expires"`
	}

	respData := &RespData{}
	err = json.Unmarshal(bodyBytes, respData)
	if err != nil {
		config.LogError(ctx, fmt.Sprintf("Failed to read api response: %s", err.Error()))
		os.Exit(1)
	}

	if *jsonOut {
		fmt.Print(string(bodyBytes))
	} else {
		// print in a way that makes copy/paste easy for the user
		fmt.Printf("doocot get %s\n", respData.Id)
		fmt.Printf("curl %s\n", respData.Url)
	}
}
