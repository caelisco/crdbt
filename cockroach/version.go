package cockroach

import (
	"fmt"
	"strings"

	"github.com/caelisco/crdbt/exec"
)

func Version() (string, error) {
	return exec.RunCombined("cockroach", "version")
}

func GetVersion() (string, error) {
	out, err := Version()
	return fmt.Sprintf("v%s", strings.Split(strings.Split(out, "\n")[0], "v")[1]), err
}
