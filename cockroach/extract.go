package cockroach

import (
	"fmt"
	"os"
	"strings"

	"github.com/caelisco/crdbt/exec"
)

func ExtractTGZ(file string) error {
	if !strings.Contains(file, "cockroach-v") {
		file = "cockroach-v" + file + ".linux-amd64.tgz"
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("input file does not exist:", file)
		return err
	}
	fmt.Print("extracting ", file)
	_, err := exec.RunOutput("tar", "-zxvf", file)
	if err != nil {
		return err
	}
	fmt.Println(" complete!")
	return err
}
