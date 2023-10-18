package action

import (
	"fmt"
	"log"
	"strings"

	"github.com/caelisco/crdbt/exec"
)

func MoveFiles(file string) {
	dir := fmt.Sprintf("./%s/", strings.TrimSuffix(file, ".tgz"))
	fmt.Println("file:", file)
	fmt.Println("dir:", dir)
	log.Println("Copying cockroach into /usr/local/bin/ ")

	out, err := exec.RunOutput("sudo", "cp", fmt.Sprintf("%scockroach", dir), "/usr/local/bin")
	if err != nil {
		fmt.Println(err)
		fmt.Println(out)
		return
	}
	// Install extra files - will assume that the mkdir command will fail
	out, err = exec.RunOutput("sudo", "mkdir", "-p", "/usr/local/lib/cockroach")
	if err != nil {
		fmt.Println(err)
		fmt.Println(out)
		return
	}
	log.Print("Copying libgeos.so to /user/local/lib/cockroach/ ")
	out, err = exec.RunOutput("sudo", "cp", "-i", fmt.Sprintf("%slib/libgeos.so", dir), "/usr/local/lib/cockroach/")
	if err != nil {
		fmt.Println(err)
		fmt.Println(out)
		return
	}
	log.Print("Copying libgeos_c.so to /user/local/lib/cockroach/ ")
	out, err = exec.RunOutput("sudo", "cp", "-i", fmt.Sprintf("%slib/libgeos_c.so", dir), "/usr/local/lib/cockroach/")
	if err != nil {
		fmt.Println(err)
		fmt.Println(out)
		return
	}
	fmt.Println("Finished!")
}

func RollbackMove(file string) {
	// TODO: Implement this
}
