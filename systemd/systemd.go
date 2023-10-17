package systemd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/caelisco/crdbt/exec"
)

//go:embede cockroach.service
var cockroachService []byte

// sudo systemctl start cockroach
func Start() (string, error) {
	return exec.RunCombined("sudo", "systemctl", "start", "cockroach")
}

// sudo systemctl stop cockroach
func Stop() (string, error) {
	return exec.RunCombined("sudo", "systemctl", "stop", "cockroach")
}

// sudo systemctl status cockroach
func Status() (string, error) {
	return exec.RunCombined("sudo", "systemctl", "status", "cockroach")
}

// sudo systemctl restart cockroach
func Restart() (string, error) {
	return exec.RunCombined("sudo", "systemctl", "restart", "cockroach")
}

// sudo systemctl reload cockroach
func Reload() (string, error) {
	return exec.RunCombined("sudo", "systemctl", "reload", "cockroach")
}

// sudo systemctl enable cockroach
func Enable() (string, error) {
	return exec.RunCombined("sudo", "systemctl", "enable", "cockroach")
}

// sudo systemctl disable cockroach
func Disable() (string, error) {
	return exec.RunCombined("sudo", "systemctl", "disable", "cockroach")
}

// sudo systemctl reload-daemon
func DaemonReload() (string, error) {
	return exec.RunCombined("sudo", "systemctl", "daemon-reload")
}

func CreateUser() (string, error) {
	var c string
	fmt.Print("You are about to add the user 'cockroach' to the system. Continue? Y/N: ")
	fmt.Scanf("%s", &c)
	if strings.EqualFold(c, "y") {
		out, err := exec.CreateUser()
		fmt.Println("Added user to system")
		return out, err
	}
	fmt.Println("Did not add user to system")
	return "", nil
}

// create the cockroach.service file
func CreateServiceFile() {
	err := os.WriteFile("cockroach.service", []byte(cockroachService), 0755)
	if err != nil {
		log.Println(err)
	}
}

// install the cockroach.service file to /etc/systemd/system/
func InstallService() {
	exec.RunCombined("sudo", "mv", "cockroach.service", "/etc/systemd/system/")
}

// remove the service by copying the file out
func UninstallService() {
	exec.RunCombined("sudo", "mv", "/etc/systemd/system/cockroach.service", ".")
}
