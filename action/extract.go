package action

import (
	"fmt"
	"os"
	"strings"

	"github.com/caelisco/crdbt/exec"
)

func ExtractTGZ(file string) error {
	fmt.Println("file:", file)
	if !strings.Contains(file, ".tgz") {
		return fmt.Errorf("unknown extension, cannot extract from %s", file)
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", file)
	}
	fmt.Print("extracting ", file)
	_, err := exec.RunOutput("tar", "-zxvf", file)
	if err != nil {
		return err
	}
	fmt.Println(" complete!")
	return err
}
