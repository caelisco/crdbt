package action

import (
	"fmt"
	"os"
	"strings"

	"github.com/caelisco/crdbt/exec"
	"github.com/gookit/color"
)

func ExtractTGZ(file string) error {
	if !strings.Contains(file, ".tgz") {
		return fmt.Errorf("unknown extension, cannot extract from %s", file)
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", file)
	}
	color.Printf("Extracting <yellow>%s</>", file)
	_, err := exec.RunOutput("tar", "-zxvf", file)
	if err != nil {
		return err
	}
	color.Println(" <green>complete!</>")
	return err
}
