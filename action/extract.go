package action

import (
	"fmt"
	"strings"

	"github.com/caelisco/crdbt/exec"
	"github.com/gookit/color"
)

func ExtractTGZ(file string) error {
	if !strings.Contains(file, ".tgz") {
		return fmt.Errorf("unknown extension, cannot extract from %s", file)
	}
	if !FileExists(file) {
		return fmt.Errorf("input file does not exist: %s", file)
	}
	done := make(chan bool)
	go TaskSpinner(done)
	color.Printf("Extracting <yellow>%s</> ", file)
	_, err := exec.RunOutput("tar", "-zxvf", file)
	done <- true
	if err != nil {
		return err
	}
	color.Println(" <green>complete!</>")
	return err
}
