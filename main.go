package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/thisdougb/doocot/config"
)

var (
	Version = "dev"     // compiler injects git tag version
	Commit  = "none"    // short commit id from git version
	Date    = "unknown" // date binary was built
)

const doocotUsage = `Usage:
  doocot put [options] <data>
  doocot get [options] <id>`

func main() {

	ctx := config.SetContextCorrelationId(
		context.Background(),
		fmt.Sprintf("%d", time.Now().Unix()))

	// high level usage if no subcommand is given
	if len(os.Args) < 2 {
		fmt.Println(doocotUsage)
		os.Exit(1)
	}

	switch os.Args[1] {

	case "put":
		put(ctx, os.Args[2:])
		os.Exit(0)

	case "get":
		get(ctx, os.Args[2:])
		os.Exit(0)

	default:
		fmt.Println(doocotUsage)
	}

}
