//go:build dev

package config

import (
	"context"
	"testing"
)

func TestCorrelationIdContext(t *testing.T) {

	var TestCases = []struct {
		description string
		value       string
		key         CorrelationContextKey
	}{
		{
			description: "test set and get id",
			key:         "cid",
			value:       "abc-123456-123456",
		},
	}

	for _, tc := range TestCases {

		ctx := SetContextCorrelationId(context.Background(), tc.value)
		result := GetContextCorrelationId(ctx)

		if tc.value != result {
			t.Error(tc.description)
		}
	}
}

func TestAppendToCid(t *testing.T) {

	ctx := SetContextCorrelationId(context.Background(), "testId")
	if "testId" != GetContextCorrelationId(ctx) {
		t.Error("initial cid")
	}

	ctx = AppendToContextCorrelationId(ctx, "someText")
	if "testId-someText" != GetContextCorrelationId(ctx) {
		t.Error("appended cid")
	}
}
