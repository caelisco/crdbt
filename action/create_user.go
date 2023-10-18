package action

import (
	"fmt"
	"strings"

	"github.com/caelisco/crdbt/exec"
)

func CreateUser() (string, error) {
	var input, out string
	var err error
	fmt.Print("You are about to add the user 'cockroach' to the system. Continue? Y/N: ")
	fmt.Scanf("%s", &input)
	if strings.EqualFold(input, "y") {
		out, err = exec.CreateUser()
	} else {
		fmt.Println("Did not add user to system")
	}
	return out, err
}
