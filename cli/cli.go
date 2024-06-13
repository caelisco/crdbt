package cli

import (
	"fmt"
	"log"
	"os"
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

	case "interactive":
		result, _ := action.Interactive()
		fmt.Printf("%v", result)
		res, err := action.GetReleases()
		fmt.Printf("%v", res)
		if err != nil {
			fmt.Println(err)
			return
		}

		ver, err := action.InteractiveVersion(res)
		if err != nil {
			fmt.Println(err)
			return
		}
		r, err := action.InteractiveRelease(ver, res)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%v\n", r)

	case "version":
		if !cockroach.Installed() {
			fmt.Println("Cockroach does not appear to be installed")
			fmt.Println("Use: crdbt install")
			return
		}
		out, err := cockroach.Version()
		if err != nil {
			color.Println("<fg=white;bg=red;>Error:</>", err)
		}
		fmt.Println("crdbt version:    0.2.0")
		fmt.Println("Runtime:         ", action.GetOS())
		fmt.Println(out)

	case "download":
		if ok := argCountCheck(args, 2); ok {
			_, err := action.Download(args[1])
			if err != nil {
				color.Println("<fg=white;bg=red>Error:</>", err)
				color.Println("<cyan>Potential fix:</> Use crdbt download <yellow>latest</>")
				color.Println("<cyan>Potential fix:</> Use a pattern like <yellow>23.1.1</>")
				fmt.Println("crdbt supports downloading alpha, beta and rc versions")
			}
			return
		}
		usage("not enough command line arguments to download.\n<cyan>Try:  </> crdbt download latest")

	case "update":
		// crdbt update
		// Works similar to apt update where it provides a list of updates that are available
		// It gets the version of CockroachDB installed, then gets the full list of releases from the
		// Cockroach Labs website.

		if !cockroach.Installed() {
			fmt.Println("Cockroach does not appear to be installed")
			fmt.Println("Use: crdbt install")
			return
		}

		installed, err := cockroach.GetVersion()
		if err != nil {
			color.Println("<fg=white;bg=red;>Error:</>", err)
			color.Println("<cyan>Potential fix:</> Use <yellow>crdbt install latest</>")
			return
		}

		releases, err := action.GetReleases()
		if err != nil {
			color.Println("<fg=white;bg=red;>Error:</>", err)
			return
		}

		// only interested in the major/minor so we trim from the pos of the last fullstop
		index := strings.LastIndex(installed, ".")
		installedv := installed[:index]
		var latestCurrent string
		var latestGlobal string

		// the last item in the array containes the very latest available release
		if len(releases) > 0 {
			latestGlobal = releases[len(releases)-1].Releases[0].Version
		}

		for _, version := range releases {
			if version.VersionPrefix == installedv {
				latestCurrent = version.Releases[0].Version
				break
			}
		}

		fmt.Printf("| Installed           | %-20s | %-25s |\n", installed, " ")
		fmt.Printf("| Latest in installed | %-20s | ", latestCurrent)
		if latestCurrent > installed {
			color.Printf("%-25s |\n", "<yellow>crdbt upgrade current</>")
		} else {
			fmt.Printf("%25s |\n", " ")
		}
		fmt.Printf("| Latest build        | %-20s | ", latestGlobal)
		if latestGlobal > latestCurrent {
			color.Printf("%-25s |\n", "<yellow>crdbt upgrade latest</>")
		} else {
			fmt.Printf("%25s |\n", " ")
		}
		// provide a warning
		words := []string{"alpha", "beta", "rc"}
		for _, word := range words {
			if strings.Contains(latestGlobal, word) {
				color.Printf("<fg=white;bg=red;>Warning:</> <yellow>this is a %s version</>\n", word)
			}
		}
	case "upgrade":
		if !cockroach.Installed() {
			fmt.Println("Cockroach does not appear to be installed")
			fmt.Println("Use: crdbt install")
			return
		}
		if ok := argCountCheck(args, 2); ok {
			//cockroach.Upgrade(args[1])
			return
		}
		usage("Not enough command line arguments to upgrade.\n\t Try: crdbt upgrade latest")

	case "list":
		if len(args) == 1 {
			action.PrintReleases()
		}
		if len(args) == 2 {
			action.PrintReleases(args[1])
		}

	case "extract":
		if ok := argCountCheck(args, 2); ok {
			err := action.ExtractTGZ(args[1])
			if err != nil {
				color.Println("<fg=white;bg=red>Error:</>", err)
				return
			}
			return
		}
		color.Println("<fg=white;bg=red>ERROR:</> Provide a filename to extract")
		color.Println("Usage: crdbt extract <yellow>[filename]</>")

	case "install":
		if ok := argCountCheck(args, 2); ok {
			err := action.Install(args[1])
			if err == nil {
				fmt.Println(err)
				return
			}
		}
		usage("Not enough command line arguments to install.\n\tTry: crdbt install latest")
	case "tidy":
		action.Tidy()

	case "certs-dir":
		fmt.Print(cockroach.GetCertsDir())

	case "systemd":
		switch args[1] {
		case "start":
			systemd.Start()

		case "stop":
			systemd.Stop()

		case "status":
			systemd.Status()

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
			out, err := action.CreateUser()
			if err != nil {
				log.Fatal(err)
			}
			// using strings.TrimSpace to remove newline character
			fmt.Println(strings.TrimSpace(out))

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
		systemd.Status()

	case "restart":
		systemd.Restart()

	case "reload":
		systemd.Reload()

	default:
		usage("No such argument exists!")
	}
}

func usage(err string) {
	fmt.Println("\ncrdbt: a command line utility to work with CockroachDB")
	if err != "" {
		color.Println("\n<fg=white;bg=red;>ERROR:</>", err)
		fmt.Println()
	}
	fmt.Println("crdbt command [options] <version>")
	fmt.Println()
	txtformat("download <version>", "download the specified version of CockroachDB")
	txtformat("download latest", "download the latest version of CockroachDB based on the releases page")
	txtformat("extract <file>", "extract the contents from the specified .tgz file")
	txtformat("install <version>", "download, extract, and install the specified version")
	txtformat("install latest", "download, extract, and install the latest version")
	txtformat("list", "list all releases of CockroachDB")
	txtformat("tidy", "Clean up old downloads and extracted files")
	txtformat("update", "check to see if there are any updates available for CockroachDB")
	txtformat("upgrade <version>", "upgrade CockroachDB to the specified version")
	txtformat("upgrade latest", "upgrade CockroachDB to the latest version based on the releases page")
	txtformat("version", "output the version of crdbt and CockroachDB")

	fmt.Println("\nCockroach commands (through systemd)")
	txtformat("reload", "Reload CochroachDB")
	txtformat("restart", "Restart CockroachDB")
	txtformat("start", "Start CockroachDB")
	txtformat("status", "Get the status of CockroachDB")
	txtformat("stop", "Stop CockroachDB")

	fmt.Println("\nsystemd commands")

	txtformat("systemd create", "create a cockroach.service file in the current directory")
	txtformat("systemd create-user", "create a cockroach user")
	txtformat("systemd daemon-reload", "run systemctl daemon-reload")
	txtformat("systemd disable", "disable the CockroachDB service running at boot")
	txtformat("systemd enable", "enable the CockroachDB service to run at boot")
	txtformat("systemd install", "install the cockroach.service file to /etc/systemd/system/")
	txtformat("systemd uninstall", "uninstall the cockroach.service file from /etc/systemd/system/")

	fmt.Println()
}

func txtformat(cmd string, desc string) {
	// TODO: can use text formatting to add offsets to ensure this always aligns properly
	// just need to find that fmt.fprintln
	fmt.Printf("%10s %-30s %s \n", "crdbt", cmd, desc)
}
