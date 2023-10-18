package cockroach

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/caelisco/crdbt/exec"
	"github.com/caelisco/crdbt/systemd"
	"github.com/gookit/color"
)

func GetCertsDir() string {
	out, err := systemd.Status()
	if err != nil {
		color.Println("<fg=white;bg=red;>Error:</>", out)
		return ""
	}
	lines := strings.Split(out, "\n")
	for _, v := range lines {
		if strings.Contains(v, "--certs-dir=") {
			dir := strings.Split(strings.SplitAfter(v, "--certs-dir=")[1], " ")
			return dir[0]
		}
	}
	return ""
}

func Update() {
	instver := GetVersion()
	instveri, _ := strconv.Atoi(strings.ReplaceAll(instver, ".", ""))
	relver, _ := GetReleases(false)
	relveri, _ := strconv.Atoi(strings.ReplaceAll(relver, ".", ""))

	if relveri == instveri {
		fmt.Println("You have the latest version of CockroachDB")
		fmt.Println("Installed version:", instver)
		fmt.Println("Latest version:", relver)
	} else {
		if relveri > instveri {
			fmt.Println(" ╭──────────────────────────────────────────────────────────────────────╮")
			fmt.Println(" │\tA new version of CockroachDB is available!                      │")
			fmt.Println(" │\tCurrent version:", instver, "\t\t\t\t\t│")
			fmt.Println(" │\tUpgrade version:", relver, "\t\t\t\t\t│")
			fmt.Println(" │\tTo update to the latest version, use \"crdbt upgrade latest\"\t│")
			fmt.Println(" ╰──────────────────────────────────────────────────────────────────────╯")
			fmt.Println("To view available versions of CockroachDB use \"crdbt list\"")
			fmt.Println("Read the release notes here: https://www.cockroachlabs.com/docs/releases/index.html")
		}
	}
}

func Upgrade(version string) {
	var uri string
	ver := version

	if strings.EqualFold(ver, "latest") {
		version, uri = GetReleases(false)
		fmt.Print("Latest version is ", version, ", proceed with upgrade? Y/N: ")
		var check string
		fmt.Scanf("%s", &check)
		if !strings.EqualFold(check, "Y") {
			fmt.Println("Not proceeding with upgrade")
			return
		}
	}

	filename, err := action.Download(ver, uri)
	if err != nil {
		fmt.Println(err)
	}
	action.ExtractTGZ(filename)
	log.Println("Stopping cockroach...")
	systemd.Stop()
	log.Println("Stopped!")
	log.Println("Copying over cockroach into /usr/local/bin/ ")

	dir := "./cockroach-v" + version + ".linux-amd64/"

	exec.RunOutput("sudo", "cp", dir+"cockroach", "/usr/local/bin")
	// Install extra files - will assume that the mkdir command will fail
	exec.RunOutput("sudo", "mkdir", "-p", "/usr/local/lib/cockroach")
	log.Print("Copying over libgeos.so to /user/local/lib/cockroach/ ")
	exec.RunOutput("sudo", "cp", "-i", dir+"lib/libgeos.so", "/usr/local/lib/cockroach/")
	log.Print("Copying over libgeos_c.so to /user/local/lib/cockroach/ ")
	exec.RunOutput("sudo", "cp", "-i", dir+"lib/libgeos_c.so", "/usr/local/lib/cockroach/")

	log.Println("Starting cockroach service...")
	systemd.Start()
	log.Println("Complete. Use ./crdbt status")
}
