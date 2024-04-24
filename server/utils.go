package main

import (
	"encoding/json"
)

func PrettyJSON(in interface{}) (string, error) {
	bb, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bb), nil
}
