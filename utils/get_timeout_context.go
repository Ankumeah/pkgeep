package utils

import (
	"context"
	"time"
)

func GetTimeoutContext(config map[string]any) (context.Context, func()) {
	cancel := func() {}
	ctx := context.Background()

	timeout := GetFromConfig[float64]("timeout_seconds", config)
	if timeout < 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
	}

	return ctx, cancel
}
