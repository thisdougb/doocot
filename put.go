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
  -v             verbose commentary
  -create <n>    generate a new secret of length n (0 < n <= 100)
  -once          expire the data immediately after first access
  -words         return a word passphrase link rather than hex string
  -json          output in json format

Example:
  $ doocot put -words this is my secret text
  doocot get slight-step-zoo-flock
  curl https://doocot.sh/collect/slight-step-zoo-flock
`

func put(ctx context.Context, args []string) {

	// we want the samllest request data size, so omit empty or unused fields
	type RequestData struct {
		DataValue     string `json:"data_value,omitempty"`
		Create        int    `json:"create,omitempty,omitzero"`
		Once          *bool  `json:"once,omitempty"`  // as ptr lets us omitempty
		Words         *bool  `json:"words,omitempty"` // as ptr lets us omitempty
		Lang          string `json:"lang"`
		ClientVersion string `json:"client_version"`
	}

	putDataUrl, _ := url.JoinPath(config.StringValue("DOOCOT_HOST"), "/api/data")

	fs := flag.NewFlagSet("put", flag.ExitOnError)
	fs.Usage = func() { fmt.Print(putUsage) }

	verbose := fs.Bool("v", false, "enable verbose commentary")
	create := fs.Int("create", 0, "remote creates a rand data value of length n")
	once := fs.Bool("once", false, "expire data after it is read once")
	words := fs.Bool("words", false, "return word passphrase instead of hex str")
	lang := fs.String("lang", "", "use this langauge code")
	jsonOut := fs.Bool("json", false, "output in json format")
	fs.Parse(args)

	if *verbose {
		ctx = config.SetContextDebug(ctx, *verbose)
	}

	config.LogDebug(ctx, fmt.Sprintf("doocot version: %s", Version))

	config.LogDebug(ctx,
		fmt.Sprintf("Using backend %s", config.StringValue("DOOCOT_HOST")))

	var requestData RequestData

	// create is mutually exclusive with a supplied data value.
	if *create > 0 {
		requestData.Create = *create
	} else {
		// get the supplied data value. the data value is the rest of the
		// command line after the flags args.
		if len(fs.Args()) == 0 {
			fs.Usage()
			os.Exit(1)
		}

		value := strings.Join(fs.Args(), " ")
		requestData.DataValue = value
	}

	// we only need to send the bool values if they are true
	if *once {
		requestData.Once = once
	}
	if *words {
		requestData.Words = words
	}
	requestData.ClientVersion = Version

	switch *lang {
	case "de":
		requestData.Lang = "de"
	case "es":
		requestData.Lang = "es"
	case "fr":
		requestData.Lang = "fr"
	default:
		requestData.Lang = "en"
	}

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

	config.LogDebug(ctx,
		fmt.Sprintf("Making request %s://%s%s (body size %d bytes)",
			req.URL.Scheme, req.URL.Host, req.URL.Path, len(requestDataJsonBytes)))

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

	config.LogDebug(ctx, fmt.Sprintf("Received response (body size %d bytes)", len(bodyBytes)))

	type RespData struct {
		Id         string `json:"id"`
		Url        string `json:"url"`
		Expires    string `json:"expires"`
		Encryption struct {
			Algorithm string `json:"algorithm"`
			Mode      string `json:"mode"`
		} `json:"encryption"`
		Scrypt struct {
			Salt      []byte `json:"salt"`       // for passphrase use
			N         int    `json:"n"`          // for passphrase use
			R         int    `json:"r"`          // for passphrase use
			P         int    `json:"p"`          // for passphrase use
			KeyLength int    `json:"key_length"` // for passphrase use
		} `json:"scrypt"`
	}

	respData := &RespData{}
	err = json.Unmarshal(bodyBytes, respData)
	if err != nil {
		config.LogError(ctx, fmt.Sprintf("Failed to read api response: %s", err.Error()))
		os.Exit(1)
	}

	config.LogDebug(ctx, fmt.Sprintf("Data passphrase: %s", respData.Id))
	config.LogDebug(ctx, fmt.Sprintf("Data url: %s", respData.Url))
	config.LogDebug(ctx, fmt.Sprintf("Data expires: %s", respData.Expires))
	config.LogDebug(ctx, fmt.Sprintf("Encryption: %+v", respData.Encryption))
	config.LogDebug(ctx, fmt.Sprintf("Scrypt: %+v", respData.Scrypt))

	if *jsonOut {
		fmt.Print(string(bodyBytes))
	} else {
		// print in a way that makes copy/paste easy for the user
		fmt.Printf("doocot get %s\n", respData.Id)
		fmt.Printf("curl %v\n", respData.Url)
	}
}
