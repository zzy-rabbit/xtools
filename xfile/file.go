package xfile

import (
	"context"
	"errors"
	"os"
)

func IsExist(ctx context.Context, path string) bool {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		return true
	}
	return true
}

func GetFileSize(ctx context.Context, file string) (int64, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}
