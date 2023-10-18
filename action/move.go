package action

import (
	"fmt"
	"strings"

	"github.com/caelisco/crdbt/exec"
	"github.com/gookit/color"
)

func MoveFiles(file string) {
	dir := fmt.Sprintf("./%s/", strings.TrimSuffix(file, ".tgz"))

	exec.RunCombined("sudo", "-v")
	color.Print("Copying <yellow>cockroach</> into <cyan>/usr/local/bin/</> ")

	out, err := exec.RunOutput("sudo", "cp", fmt.Sprintf("%scockroach", dir), "/usr/local/bin")
	if err != nil {
		color.Println("<tomato>failed!</>")
		fmt.Println(err)
		fmt.Println(out)
		return
	}
	color.Println("<green>complete!</>")
	// Install extra files - will assume that the mkdir command will fail
	out, err = exec.RunCombined("sudo", "mkdir", "-p", "/usr/local/lib/cockroach")
	if err != nil {
		fmt.Println(err)
		fmt.Println(out)
		return
	}

	color.Print("Copying <yellow>libgeos.so</> to <cyan>/user/local/lib/cockroach/</> ")
	out, err = exec.RunCombined("sudo", "cp", "-i", fmt.Sprintf("%slib/libgeos.so", dir), "/usr/local/lib/cockroach/")
	if err != nil {
		color.Println("<tomato>failed!</>")
		fmt.Println(err)
		fmt.Println(out)
		return
	}
	color.Println("<green>complete!</>")

	color.Print("Copying <yellow>libgeos_c.so</> to <cyan>/user/local/lib/cockroach/</> ")
	out, err = exec.RunCombined("sudo", "cp", "-i", fmt.Sprintf("%slib/libgeos_c.so", dir), "/usr/local/lib/cockroach/")
	if err != nil {
		color.Println("<tomato>failed!</>")
		fmt.Println(err)
		fmt.Println(out)
		return
	}
	color.Println("<green>complete!</>")
}

func RollbackMove(file string) {
	// TODO: Implement this
}
