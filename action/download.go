package action

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/caelisco/crdbt/cockroach"
	"github.com/schollz/progressbar/v3"
)

func Download(version string, uri string) (string, error) {
	if strings.EqualFold(version, "latest") {
		version, uri = cockroach.GetReleases(false)
	}
	file := "cockroach-v" + version + ".linux-amd64.tgz"
	if uri == "" {
		uri = "https://binaries.cockroachdb.com/" + file
	}

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		fmt.Print("The file already exists. Extract (e) / Overwrite (o): ")
		var option string
		fmt.Scanf("%s", &option)
		if strings.EqualFold(option, "e") {
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
