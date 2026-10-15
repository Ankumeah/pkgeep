package modes

import (
	"strings"
)

func isFlag(arg string, flag flag) (string, bool) {
	if strings.HasPrefix(arg, flag.name) {
		return strings.TrimPrefix(arg, flag.name+"="), true
	} else if flag.shortName != "" && strings.HasPrefix(arg, flag.shortName) {
		return strings.TrimPrefix(arg, flag.shortName+"="), true
	} else {
		return "", false
	}
}
