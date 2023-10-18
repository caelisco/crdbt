package systemd

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/caelisco/crdbt/exec"
	"github.com/gookit/color"
)

//go:embed cockroach.service
var cockroachService []byte

// sudo systemctl start cockroach
func Start() (out string, err error) {
	out, err = exec.RunCombined("sudo", "systemctl", "start", "cockroach")
	if err != nil {
		color.Println("<fg=white;bg=red>Warning:</>", strings.TrimSpace(out))
		return
	}
	return
}

// sudo systemctl stop cockroach
func Stop() (out string, err error) {
	out, err = exec.RunCombined("sudo", "systemctl", "stop", "cockroach")
	if err != nil {
		color.Println("<fg=white;bg=red>Warning:</>", strings.TrimSpace(out))
		return
	}
	return
}

// sudo systemctl status cockroach
func Status() (out string, err error) {
	out, err = exec.RunCombined("sudo", "systemctl", "status", "cockroach")
	if err != nil {
		color.Println("<fg=white;bg=red>Warning:</>", strings.TrimSpace(out))
		return
	}
	fmt.Println(out)
	return
}

// sudo systemctl restart cockroach
func Restart() (out string, err error) {
	out, err = exec.RunCombined("sudo", "systemctl", "restart", "cockroach")
	if err != nil {
		color.Println("<fg=white;bg=red>Warning:</>", strings.TrimSpace(out))
		return
	}
	return
}

// sudo systemctl reload cockroach
func Reload() (out string, err error) {
	out, err = exec.RunCombined("sudo", "systemctl", "reload", "cockroach")
	if err != nil {
		color.Println("<fg=white;bg=red>Warning:</>", strings.TrimSpace(out))
		return
	}
	return
}

// sudo systemctl enable cockroach
func Enable() (out string, err error) {
	out, err = exec.RunCombined("sudo", "systemctl", "enable", "cockroach")
	if err != nil {
		color.Println("<fg=white;bg=red>Warning:</>", strings.TrimSpace(out))
		return
	}
	return
}

// sudo systemctl disable cockroach
func Disable() (out string, err error) {
	out, err = exec.RunCombined("sudo", "systemctl", "disable", "cockroach")
	if err != nil {
		color.Println("<fg=white;bg=red>Warning:</>", strings.TrimSpace(out))
		return
	}
	return
}

// sudo systemctl reload-daemon
func DaemonReload() (out string, err error) {
	out, err = exec.RunCombined("sudo", "systemctl", "daemon-reload")
	if err != nil {
		color.Println("<fg=white;bg=red>Warning:</>", strings.TrimSpace(out))
		return
	}
	return
}

// create the cockroach.service file
func CreateServiceFile() {
	err := os.WriteFile("cockroach.service", cockroachService, 0755)
	if err != nil {
		log.Println(err)
	}
}

// install the cockroach.service file to /etc/systemd/system/
func InstallService() {
	out, err := exec.RunCombined("sudo", "mv", "cockroach.service", "/etc/systemd/system/")
	if err != nil {
		color.Println("<fg=white;bg=red>Warning:</>", strings.TrimSpace(out))
		return
	}
}

// remove the service by copying the file out
func UninstallService() {
	out, err := exec.RunCombined("sudo", "mv", "/etc/systemd/system/cockroach.service", ".")
	if err != nil {
		color.Println("<fg=white;bg=red>Warning:</>", strings.TrimSpace(out))
		return
	}
}
