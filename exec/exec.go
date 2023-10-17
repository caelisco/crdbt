package exec

import (
	"log"
	"os/exec"
)

var Verbose bool = false

func RunCommand(command string, commands ...string) error {
	if Verbose {
		log.Println("Command:", command, commands)
	}
	err := exec.Command(command, commands...).Run()
	return err
}

func RunOutput(command string, commands ...string) (string, error) {
	if Verbose {
		log.Println("Command:", command, commands)
	}

	cmd := exec.Command(command, commands...)
	out, err := cmd.Output()

	if Verbose {
		log.Println("Output :", string(out))
	}

	return string(out), err
}

func RunCombined(command string, commands ...string) (string, error) {
	if Verbose {
		log.Println("command: ", command, commands)
	}
	cmd := exec.Command(command, commands...)
	out, err := cmd.CombinedOutput()
	if Verbose {
		log.Println(err)
		log.Println("output: ", string(out))
	}
	return string(out), err
}

func CreateUser() (string, error) {
	return RunOutput("sudo", "adduser", "--system", "--disabled-password", "--disabled-login", "--no-create-home", "--shell", "/bin/false", "--group", "cockroach")
}
