package utils

import "os"

func FileExist(path string) (bool, error) {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err), err
}
