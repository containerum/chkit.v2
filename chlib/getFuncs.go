package chlib

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadJsonFromFile(path string, b interface{}) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	err = json.NewDecoder(file).Decode(&b)
	return
}

func GetCmdRequestJson(client *Client, kind, name string) (ret []GenericJson, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("can`t extract field: %s", r)
		}
	}()
	apiResult, err := client.Get(kind, name, client.userConfig.Namespace)
	if err != nil {
		return ret, err
	}
	items := apiResult["results"].([]interface{})
	for _, itemI := range items {
		ret = append(ret, itemI.(map[string]interface{}))
	}
	return
}
