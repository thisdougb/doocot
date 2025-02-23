package config

import (
	"context"
	"fmt"
	"time"
)

func LogInfo(ctx context.Context, msg string) {
	writeToLog(ctx, "INFO", msg)
}

func LogError(ctx context.Context, msg string) {
	writeToLog(ctx, "ERROR", msg)
}

func LogDebug(ctx context.Context, msg string) {
	if GetContextDebug(ctx) {
		writeToLog(ctx, "DEBUG", msg)
	}
}

func writeToLog(ctx context.Context, severity string, msg string) {

	fmt.Printf("%s +%s %s\n",
		severity,
		sinceCreated(ctx),
		msg)
}

func sinceCreated(ctx context.Context) string {

	createdTime := time.Unix(GetContextTimeCreated(ctx), 0)
	t := time.Since(createdTime).Seconds()

	return fmt.Sprintf("%.1fs", t)
}
