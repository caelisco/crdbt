package cockroach

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/caelisco/crdbt/exec"
	"github.com/caelisco/crdbt/systemd"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html"
)

func Version() (string, error) {
	return exec.RunOutput("cockroach", "version")
}

func GetVersion() string {
	out, _ := Version()
	return strings.Split(strings.Split(out, "\n")[0], "v")[1]
}

func GetCertsDir() string {
	out, _ := systemd.Status()
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
	fmt.Println("Installed version of cockroach DB:", instver)
	instveri, _ := strconv.Atoi(strings.ReplaceAll(instver, ".", ""))
	relver, _ := GetReleases(false)
	relveri, _ := strconv.Atoi(strings.ReplaceAll(relver, ".", ""))
	fmt.Println("Available version of cockroach DB:", relver)

	if relveri == instveri {
		fmt.Println("You have the latest version of Cockroach DB installed")
	} else {
		if relveri > instveri {
			fmt.Println("You can upgrade to version", relver, "of Cockroach DB!")
			fmt.Println("Options:")
			fmt.Println("\tcrdbt upgrade latest")
			fmt.Println("\tcrdbt upgrade <version>")
			fmt.Println("\tcrdbt list")
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

	filename, err := Download(ver, uri)
	if err != nil {
		fmt.Println(err)
	}
	ExtractTGZ(filename)
	log.Println("Stopping cockroach...")
	systemd.Stop()
	log.Println("Stopped!")
	log.Println("copying over cockroach into /usr/local/bin/")

	dir := "./cockroach-v" + version + ".linux-amd64/"

	exec.RunOutput("sudo", "cp", dir+"cockroach", "/usr/local/bin")
	// Install extra files - will assume that the mkdir command will fail
	exec.RunOutput("sudo", "mkdir", "-p", "/usr/local/lib/cockroach")
	exec.RunOutput("sudo", "cp", "-i", dir+"lib/libgeos.so", "/usr/local/lib/cockroach/")
	exec.RunOutput("sudo", "cp", "-i", dir+"lib/libgeos_c.so", "/usr/local/lib/cockroach/")

	log.Println("Starting cockroach...")
	systemd.Start()
}

func Download(version string, uri string) (string, error) {
	if strings.EqualFold(version, "latest") {
		version, uri = GetReleases(false)
	}
	fmt.Println(version, uri)
	file := "cockroach-v" + version + ".linux-amd64.tgz"
	if uri == "" {
		uri = "https://binaries.cockroachdb.com/" + file
	}

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		fmt.Print("The file already exists. Extract (e) / Overwrite (o):")
		var option string
		fmt.Scanf("%s", &option)
		if strings.EqualFold(option, "e") {
			err := ExtractTGZ(file)
			if err != nil {
				fmt.Println(" ----- error ----- ")
				fmt.Println(err)
			}
			return version, nil
		}
	}

	req, _ := http.NewRequest("GET", uri, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return version, err
	}
	defer resp.Body.Close()

	f, _ := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading Cockroach v"+version,
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
	fmt.Println("Download complete!")
	if err != nil {
		return version, err
	}
	return version, nil
}

func ExtractTGZ(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("input file does not exist")
		return err
	}
	fmt.Print("extracting ", file)
	_, err := exec.RunOutput("tar", "-zxvf", file)
	if err != nil {
		return err
	}
	fmt.Println(" complete!")
	return err
}

type Releases struct {
	List   []Latest
	Latest Latest
}

type Latest struct {
	Version string
	URI     string
}

func GetReleases(echo bool) (string, string) {
	r := Releases{}
	return r.Releases(echo)
}

func (r *Releases) Releases(echo bool) (string, string) {
	resp, err := http.Get("https://www.cockroachlabs.com/docs/releases/")
	if err != nil {
		fmt.Println("error getting a response from the release page")
		fmt.Println(err)
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
	}
	r.parse_html(doc)

	if echo {
		fmt.Printf("%25s | %20s \n", "VERSION", "DOWNLOAD URL")
		for i := len(r.List) - 1; i >= 0; i-- {
			fmt.Printf("%25s | %20s \n", r.List[i].Version, r.List[i].URI)
		}
		fmt.Printf("\n%25s \n", "Latest:")
		fmt.Printf("%25s | %20s \n", r.Latest.Version, r.Latest.URI)
	}
	return r.Latest.Version, r.Latest.URI
}

func (r *Releases) parse_html(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, e := range n.Attr {
			if e.Key == "href" {
				// only interested in Linux without verification checksums
				if strings.Contains(e.Val, "linux-amd64.tgz") && strings.Contains(e.Val, "v") && !strings.Contains(e.Val, "sha256sum") && !strings.Contains(e.Val, "sql") {
					version := strings.TrimSuffix(strings.Split(e.Val, "v")[1], ".linux-amd64.tgz")
					if !strings.Contains(e.Val, "alpha") && !strings.Contains(e.Val, "beta") && r.Latest == (Latest{}) {
						r.Latest.Version = version
						r.Latest.URI = e.Val
					}
					r.List = append(r.List, Latest{Version: version, URI: e.Val})
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.parse_html(c)
	}
}
