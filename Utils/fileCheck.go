package utils

import (
	"log/slog"
	"os"
)

func CheckIfExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		slog.Info(path + " exists")
		return true
	} else {
		slog.Info(path + " does ont exist")
		return false
	}
}
