package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadConfigSection(path string, section string) (resultMap map[string]interface{}, err error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var dat map[string]interface{}
	if err := json.Unmarshal(byteValue, &dat); err != nil {
		return nil, err
	}
	resultMap = dat[section].(map[string]interface{})
	return
}
