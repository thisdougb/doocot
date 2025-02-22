package config

import (
	"context"
	"time"

	"math/rand"
)

type (
	CorrelationContextKey string
	DebugContextKey       string
	TimeCreatedContextKey string
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func SetContextCorrelationId(ctx context.Context, value string) context.Context {

	// create a unique-ish correlation identifier
	id := make([]byte, 8)
	for idx := range len(id) {
		n := rand.Intn(len(chars))
		id[idx] = chars[n]
	}

	newctx := context.WithValue(ctx, CorrelationContextKey("cid"), string(id))

	// if the created time is unset then set it. test for -1 as 0 could be
	// a symptom of a default unset value
	t := GetContextTimeCreated(ctx)
	if t == -1 {
		newctx = context.WithValue(
			newctx,
			TimeCreatedContextKey("timeCreated"),
			time.Now().Unix())
	}

	return newctx
}

func GetContextTimeCreated(ctx context.Context) int64 {

	key := TimeCreatedContextKey("timeCreated")

	if v := ctx.Value(key); v != nil {
		return v.(int64)
	}
	return -1
}

func AppendToContextCorrelationId(ctx context.Context, value string) context.Context {
	key := CorrelationContextKey("cid")
	id := GetContextCorrelationId(ctx)
	newctx := context.WithValue(ctx, key, id+"-"+value)
	return newctx
}

func GetContextCorrelationId(ctx context.Context) string {

	key := CorrelationContextKey("cid")

	if v := ctx.Value(key); v != nil {
		return v.(string)
	}

	return "no-id"
}

func GetContextDebug(ctx context.Context) bool {

	key := DebugContextKey("debug")

	if v := ctx.Value(key); v != nil {
		return v.(bool)
	}

	return false
}

func SetContextDebug(ctx context.Context, debug bool) context.Context {
	return context.WithValue(ctx, DebugContextKey("debug"), debug)
}
