package flags

import (
	"flag"

	"github.com/File-share/constants"
)

var mode string

func Init() {
	mode = *flag.String("mode", constants.Sender, "flag to set the mode")
}

func Mode() string {
	return mode
}
