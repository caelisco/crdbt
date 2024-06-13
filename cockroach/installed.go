package cockroach

import "github.com/caelisco/crdbt/action"

func Installed() bool {
	return action.FileExists("/usr/bin/cockroach")
}
