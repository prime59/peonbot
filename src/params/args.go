package params

import "flag"

/* Command line args */

type _args struct {
	verbose bool
}

func (a *_args) Verbose() bool {
	return a.verbose
}

func getArgs() *_args {
	var args _args

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false,
		"Enables additional logging if set to true. Defaults to false.")

	flag.Parse()

	args.verbose = verbose

	return &args
}
