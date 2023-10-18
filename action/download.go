package action

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/caelisco/crdbt/cockroach"
	"github.com/schollz/progressbar/v3"
)

func Download(version string) (string, error) {
	var uri string
	if strings.EqualFold("latest", version) {
		version, uri = cockroach.GetReleases(false)
	}

	re := regexp.MustCompile(`^\d+\.\d+.\d+`)
	match := re.MatchString(version)
	if !match {
		return "", fmt.Errorf("incorrect pattern to download Cockroach")
	}

	file := fmt.Sprintf("cockroach-v%s.%s.tgz", version, GetOS())
	if uri == "" {
		uri = fmt.Sprintf("https://binaries.cockroachdb.com/%s", file)
	}

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		fmt.Print("The file already exists. Extract (e) / Overwrite (o): ")
		var option string
		fmt.Scanf("%s", &option)
		if strings.EqualFold("e", option) {
			return file, fmt.Errorf("extract")
		}
	}

	req, _ := http.NewRequest("GET", uri, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return file, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return file, fmt.Errorf("unknown version (unexpected status code: %d)", resp.StatusCode)
	}

	f, _ := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading Cockroach v"+version,
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
	log.Println("Download complete!")
	if err != nil {
		return file, err
	}
	return file, nil
}
