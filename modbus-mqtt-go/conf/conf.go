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

type Config struct {
	dat map[string]interface{}
}

func NewConfig(path string) (config *Config, err error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	config = &Config{}
	buf, err := ioutil.ReadAll(fh)
	if err = json.Unmarshal(buf, &config.dat); err != nil {
		return nil, err
	}

	return
}

func (c *Config) Value(key string) (result *Config) {
	result = &Config{}
	result.dat = c.dat[key].(map[string]interface{})
	return result
}

func (c *Config) ValueAsString(key string) (result string) {
	return c.dat[key].(string)
}
