package verbose

import "fmt"

var _verbose = false

func SetPrinter(verbose bool) {
	_verbose = verbose
}

func Vprintf(msg string, vargs ...interface{}) {
	if _verbose {
		fmt.Printf(fmt.Sprintf("[debug] %s", msg), vargs...)
	}
}
