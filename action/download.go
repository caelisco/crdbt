package action

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/caelisco/crdbt/cockroach"
	"github.com/gookit/color"
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

	color.Printf("Attempting to download <yellow>%s</>\n", file)

	if FileExists(file) {
		color.Println("<tomato>Warning:</> The file already exists")
		color.Printf("Use: crdbt extract <yellow>%s</>\n", file)
		color.Printf("Use: crdbt install <yellow>%s</>\n", file)
		return "", nil
	}

	req, _ := http.NewRequest("GET", uri, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return file, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unknown version (unexpected status code: %d)", resp.StatusCode)
	}

	f, _ := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading Cockroach v"+version,
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
	fmt.Println("Download complete!")
	if err != nil {
		return "", err
	}
	return file, nil
}
