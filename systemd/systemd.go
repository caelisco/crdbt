package systemd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/caelisco/crdbt/exec"
)

// sudo systemctl start cockroach
func Start() (string, error) {
	return exec.RunOutput("sudo", "systemctl", "start", "cockroach")
}

// sudo systemctl stop cockroach
func Stop() (string, error) {
	return exec.RunOutput("sudo", "systemctl", "stop", "cockroach")
}

// sudo systemctl status cockroach
func Status() (string, error) {
	return exec.RunOutput("sudo", "systemctl", "status", "cockroach")
}

// sudo systemctl restart cockroach
func Restart() (string, error) {
	return exec.RunOutput("sudo", "systemctl", "restart", "cockroach")
}

// sudo systemctl reload cockroach
func Reload() (string, error) {
	return exec.RunOutput("sudo", "systemctl", "reload", "cockroach")
}

// sudo systemctl enable cockroach
func Enable() (string, error) {
	return exec.RunOutput("sudo", "systemctl", "enable", "cockroach")
}

// sudo systemctl disable cockroach
func Disable() (string, error) {
	return exec.RunOutput("sudo", "systemctl", "disable", "cockroach")
}

// sudo systemctl reload-daemon
func DaemonReload() (string, error) {
	return exec.RunOutput("sudo", "systemctl", "daemon-reload")
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
	buff := `[Unit]
Description=Cockroach Database cluster node
Requires=network.target
[Service]
Type=notify
WorkingDirectory=/var/lib/cockroach
ExecStart=/usr/local/bin/cockroach start-single-node --certs-dir=/var/lib/cockroach/certs --advertise-addr=localhost --cache=.25 --max-sql-memory=.25
TimeoutStopSec=60
Restart=always
RestartSec=10
SyslogIdentifier=cockroach
User=cockroach
[Install]
WantedBy=default.target`
	err := os.WriteFile("cockroach.service", []byte(buff), 0755)
	if err != nil {
		log.Println(err)
	}
}

// install the cockroach.service file to /etc/systemd/system/
func InstallService() {
	exec.RunOutput("sudo", "mv", "cockroach.service", "/etc/systemd/system/")
}

// remove the service by copying the file out
func UninstallService() {
	exec.RunOutput("sudo", "mv", "/etc/systemd/system/cockroach.service", ".")
}
