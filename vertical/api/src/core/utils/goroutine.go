package utils

import "go.uber.org/zap"

func GoWithRecovery(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.S().Errorw("Recovered from panic in goroutine", "panic", r)
			}
		}()
		f()
	}()
}
