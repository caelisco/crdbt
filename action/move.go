package action

import (
	"fmt"
	"strings"

	"github.com/caelisco/crdbt/exec"
	"github.com/gookit/color"
)

func MoveFiles(file string) error {
	dir := fmt.Sprintf("./%s/", strings.TrimSuffix(file, ".tgz"))

	// this will force the prompt [sudo] password for <user>:
	exec.RunCombined("sudo", "-v")

	color.Print("Copying <yellow>cockroach</> into <cyan>/usr/local/bin/</> ")
	out, err := exec.RunOutput("sudo", "cp", fmt.Sprintf("%scockroach", dir), "/usr/local/bin")
	if err != nil {
		color.Println("<tomato>failed!</>")
		fmt.Println(err)
		fmt.Println(out)
		return err
	}
	color.Println("<green>complete!</>")

	// Install extra files - will assume that the mkdir command will fail
	out, err = exec.RunCombined("sudo", "mkdir", "-p", "/usr/local/lib/cockroach")
	if err != nil {
		color.Println("<fg=white;bg=red>Fatal:</> failed to create directory <cyan>/usr/local/lib/cockroach</cyan>")
		color.Println("Cockroach is still installed, just without the Geo libraries")
		fmt.Println(err)
		fmt.Println(out)
		return err
	}

	color.Print("Copying <yellow>libgeos.so</> to <cyan>/user/local/lib/cockroach/</> ")
	out, err = exec.RunCombined("sudo", "cp", "-i", fmt.Sprintf("%slib/libgeos.so", dir), "/usr/local/lib/cockroach/")
	if err != nil {
		color.Println("<tomato>failed!</>")
		fmt.Println(err)
		fmt.Println(out)
		return err
	}
	color.Println("<green>complete!</>")

	color.Print("Copying <yellow>libgeos_c.so</> to <cyan>/user/local/lib/cockroach/</> ")
	out, err = exec.RunCombined("sudo", "cp", "-i", fmt.Sprintf("%slib/libgeos_c.so", dir), "/usr/local/lib/cockroach/")
	if err != nil {
		color.Println("<tomato>failed!</>")
		fmt.Println(err)
		fmt.Println(out)
		return err
	}
	color.Println("<green>complete!</>")
	return err
}
