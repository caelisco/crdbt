package cockroach

import (
	"strings"

	"github.com/caelisco/crdbt/exec"
	"github.com/gookit/color"
)

func Version() (string, error) {
	return exec.RunCombined("cockroach", "version")
}

func GetVersion() string {
	out, err := Version()
	if err != nil {
		color.Println("17<fg=white;bg=red;>Error:</>", err)
		return ""
	}
	return strings.Split(strings.Split(out, "\n")[0], "v")[1]
}
