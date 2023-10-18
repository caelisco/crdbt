package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/caelisco/crdbt/action"
	"github.com/caelisco/crdbt/cockroach"
	"github.com/caelisco/crdbt/exec"
	"github.com/caelisco/crdbt/systemd"
	"github.com/gookit/color"
)

// This package contains all of the command line functions/commands

func argCountCheck(arg []string, min int) bool {
	return len(arg) >= min
}

func ParseArgs() {
	// Handle any special commands.
	// Special commands are automatically stripped from command line options
	var args []string
	for _, v := range os.Args[1:] {
		if strings.EqualFold(v, "version") {
			fmt.Println("crdbt - a command line utility for working with CockroachDB")
			fmt.Println("version: ")
			fmt.Println("Runtime: ", action.GetOS())
		}
		if strings.EqualFold(v, "-v") || strings.EqualFold(v, "--verbose") {
			log.Println("Running crdbt in verbose mode")
			exec.Verbose = true
		} else {
			args = append(args, v)
		}
	}
	// check that there are enough commands to be useful
	if len(args) < 1 {
		usage("Not enough command line arguments for crdbt to be useful!")
		return
	}
	// handle commands
	switch args[0] {
	case "help":
		usage("")
	case "version":
		out, err := cockroach.Version()
		if err != nil {
			color.Println("<fg=white;bg=red;>Error:</>", err)
			return
		}
		fmt.Println(out)
	case "download":
		if ok := argCountCheck(args, 2); ok {
			action.Download(args[1], "")
			return
		}
		usage("not enough command line arguments to download.\n\t Try crdbt download latest")
	case "update":
		cockroach.Update()
	case "upgrade":
		if ok := argCountCheck(args, 2); ok {
			cockroach.Upgrade(args[1])
			return
		}
		usage("Not enough command line arguments to upgrade.\n\t Try: crdbt upgrade latest")
	case "list":
		cockroach.GetReleases(true)
	case "extract":
		if ok := argCountCheck(args, 2); ok {
			action.ExtractTGZ(args[1])
			return
		}
		color.Println("<fg=white;bg=red>ERROR:</> Provide a filename to extract")
		color.Println("Usage: crdbt extract <yellow>[filename]</>")
	case "tidy":
		files, err := filepath.Glob("cockroach-v*.linux-amd64*")
		if err != nil {
			fmt.Println(err)
		}
		if len(files) == 0 {
			fmt.Println("No files to tidy")
		}
		for _, v := range files {
			err := os.RemoveAll(v)
			if err != nil {
				fmt.Println("failed to delete:", v)
			}
			fmt.Println("Deleted:", v)
		}
		if len(files) > 0 {
			color.Println("<green>SUCCESS!</> Folder is now tidy")
		}
	case "certs-dir":
		fmt.Print(cockroach.GetCertsDir())
	case "systemd":
		switch args[1] {
		case "start":
			systemd.Start()
		case "stop":
			systemd.Stop()
		case "status":
			out, _ := systemd.Status()
			fmt.Println(out)
		case "restart":
			systemd.Restart()
		case "reload":
			systemd.Reload()
		case "enable":
			systemd.Enable()
		case "disable":
			systemd.Disable()
		case "daemon-reload":
			systemd.DaemonReload()
		case "create-user":
			systemd.CreateUser()
		case "create":
			systemd.CreateServiceFile()
		case "install":
			systemd.InstallService()
		case "uninstall":
			systemd.UninstallService()
		default:
			usage("No such argument exists!")
		}
	case "start":
		systemd.Start()
	case "stop":
		systemd.Stop()
	case "status":
		out, _ := systemd.Status()
		fmt.Println(out)
	case "restart":
		systemd.Restart()
	case "reload":
		systemd.Reload()
	default:
		usage("No such argument exists!")
	}
}

func usage(err string) {
	fmt.Println("crdbt: a command line utility to help work with CockroachDB")
	if err != "" {
		color.Println("\n<fg=white;bg=red;>ERROR:</>", err)
		fmt.Println()
	}
	fmt.Println("crdbt command [options] <version>")
	fmt.Println()
	txtformat("version", "output the version of crdbt and CockroachDB")
	txtformat("update", "check to see if there are any updates available for CockroachDB")
	txtformat("upgrade <version>", "upgrade CockroachDB to the specified version")
	txtformat("upgrade latest", "upgrade CockroachDB to the latest version based on the releases page")
	txtformat("list", "list all releases of CockroachDB")
	txtformat("download <version>", "download the specified version of CockroachDB")
	txtformat("download latest", "download the latest version of CockroachDB based on the releases page")
	txtformat("extract <file>", "extract the contents from the specified .tgz file")
	txtformat("tidy", "Clean up old downloads and extracted files")
	fmt.Println("\nsystemd commands")
	txtformat("status", "alias of systemd status")
	txtformat("start", "alias of systemd start")
	txtformat("stop", "alias of systemd stop")
	txtformat("restart", "alias of systemd restart")
	txtformat("reload", "alias of systemd reload")
	txtformat("systemd enable", "enable the CockroachDB service to run at boot")
	txtformat("systemd disable", "disable the CockroachDB service running at boot")
	txtformat("systemd daemon-reload", "run systemctl daemon-reload")
	txtformat("systemd create-user", "create a cockroach user")
	txtformat("systemd create", "create a cockroach.service file in the current directory")
	txtformat("systemd install", "install the cockroach.service file to /etc/systemd/system/")
	txtformat("systemd uninstall", "uninstall the cockroach.service file from /etc/systemd/system/")

}

func txtformat(cmd string, desc string) {
	// TODO: can use text formatting to add offsets to ensure this always aligns properly
	// just need to find that fmt.fprintln
	fmt.Printf("%10s %-30s %s \n", "crdbt", cmd, desc)
}
