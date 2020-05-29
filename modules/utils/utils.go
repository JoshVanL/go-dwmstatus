package utils

import (
	"bytes"
	"io/ioutil"
	"strconv"
)

func ReadFile(path string) ([]byte, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(f), nil
}

func ReadFileToFloat(path string) (float64, error) {
	b, err := ReadFile(path)
	if err != nil {
		return -1, err
	}

	return strconv.ParseFloat(string(b), 64)
}
