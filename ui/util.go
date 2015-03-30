package ui

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
)

func HashFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(data)), nil
}
