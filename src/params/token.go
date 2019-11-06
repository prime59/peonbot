package params

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

/* token file that stores api token */

const _FILE_API_TOKEN = "tokens/token.yaml"

type _token struct {
	ApiKey string `yaml:"api_key"`
}

func readToken() (string, error) {
	var token _token

	raw, err := ioutil.ReadFile(_FILE_API_TOKEN)
	if err != nil {
		return "", err
	}

	if err := yaml.Unmarshal(raw, &token); err != nil {
		return "", err
	}

	return token.ApiKey, nil
}
