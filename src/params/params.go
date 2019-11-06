package params

import "log"

type _params struct {
	Args   *_args
	Config *_config
	token  string
}

func (p *_params) Token() string {
	return p.token
}

func New() (*_params, error) {
	var params _params
	/*
		go won't let me assign values to a struct attributes, and an
		undeclared var in the same line without declaring the err
		variable here.
	*/
	var err error

	params.Args = getArgs()
	params.Config, err = readConfig()
	if err != nil {
		return &_params{}, err
	}
	params.token, err = readToken()
	if err != nil {
		return &_params{}, err
	}

	if params.Args.Verbose() {
		log.Printf("Args: %+v\n", params.Args)
		log.Printf("Config: %+v\n", params.Config)
		log.Printf("Token: %s\n", params.Token())
	}

	return &params, nil
}
