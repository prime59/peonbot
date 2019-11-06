package params

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

/* config from yaml files in `../config` */

const _FILE_BANLIST = "config/ban_list.yaml"
const _FILE_GREETINGS = "config/greetings.yaml"
const _FILE_PRIVELEGED = "config/priveleged_list.yaml"

type _banlist struct {
	Users []string `yaml:"users"`
}

func readBanlist() (*_banlist, error) {
	var banlist _banlist

	raw, err := ioutil.ReadFile(_FILE_BANLIST)
	if err != nil {
		return &_banlist{}, err
	}

	if err := yaml.Unmarshal(raw, &banlist); err != nil {
		return &_banlist{}, err
	}

	return &banlist, nil
}

type _greetings struct {
	Msg string `yaml:"msg"`
}

func readGreetings() (*_greetings, error) {
	var greetings _greetings

	raw, err := ioutil.ReadFile(_FILE_GREETINGS)
	if err != nil {
		return &_greetings{}, err
	}

	if err := yaml.Unmarshal(raw, &greetings); err != nil {
		return &_greetings{}, err
	}

	return &greetings, nil
}

type _privelegedusers struct {
	Users []string `yaml:"users"`
}

func readPrivelege() (*_privelegedusers, error) {
	var privelege _privelegedusers

	raw, err := ioutil.ReadFile(_FILE_PRIVELEGED)
	if err != nil {
		return &_privelegedusers{}, err
	}

	if err := yaml.Unmarshal(raw, &privelege); err != nil {
		return &_privelegedusers{}, err
	}

	return &privelege, nil
}

type _config struct {
	blist     []string
	greetings string
	pusers    []string
}

func (c *_config) Blist() []string {
	return c.blist
}

func (c *_config) Greetings() string {
	return c.greetings
}

func (c *_config) Pusers() []string {
	return c.pusers
}

func readConfig() (*_config, error) {
	var config _config

	blist, err := readBanlist()
	if err != nil {
		return &_config{}, err
	}

	greetings, err := readGreetings()
	if err != nil {
		return &_config{}, err
	}

	pusers, err := readPrivelege()
	if err != nil {
		return &_config{}, err
	}

	config.blist = blist.Users
	config.greetings = greetings.Msg
	config.pusers = pusers.Users

	return &config, nil
}
