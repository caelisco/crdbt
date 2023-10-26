package action

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Release struct {
	Version      string
	Date         string
	ReleaseNotes string
	DownloadURI  string
	SHA256Sum    string
}

type VersionGroup struct {
	VersionPrefix string
	Releases      []Release
}

// GetReleases retrieves the release page from Cockroachlabs, and parses it for useful information
func GetReleases() ([]VersionGroup, error) {
	var versionGroups []VersionGroup
	var currentGroup *VersionGroup

	url := "https://cockroachlabs.com/docs/releases/"

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	doc.Find(fmt.Sprintf("section[data-scope='%s'] tr", runtime.GOOS)).Each(func(index int, element *goquery.Selection) {
		if index > 0 { // Skip header row
			version := strings.TrimSpace(element.Find("td").First().Text())
			splitVersion := strings.Split(version, ".")
			if len(splitVersion) >= 2 {
				versionPrefix := strings.TrimSpace(strings.Join(splitVersion[:2], "."))
				if versionPrefix != "Version" {
					if currentGroup == nil || currentGroup.VersionPrefix != versionPrefix {
						if currentGroup != nil {
							versionGroups = append(versionGroups, *currentGroup)
						}
						currentGroup = &VersionGroup{
							VersionPrefix: versionPrefix,
						}
					}
					date := strings.TrimSpace(element.Find("td").Eq(1).Text())
					releaseNotes := fmt.Sprintf("%s%s", "https://www.cockroachlabs.com", strings.TrimSpace(element.Find("td a").First().AttrOr("href", "")))
					downloadURI := strings.TrimSpace(element.Find("td a:contains('Full Binary')").AttrOr("href", ""))
					if downloadURI == "" {
						withdrawnText := element.Find("td span.badge").Text()
						if strings.Contains(withdrawnText, "Withdrawn") {
							downloadURI = "WITHDRAWN"
						}
					}

					sha256SumLink := strings.TrimSpace(element.Find("td a:contains('SHA256')").AttrOr("href", ""))
					release := Release{
						Version:      version,
						Date:         date,
						ReleaseNotes: releaseNotes,
						DownloadURI:  downloadURI,
						SHA256Sum:    sha256SumLink,
					}
					currentGroup.Releases = append(currentGroup.Releases, release)
				}
			}
		}
	})

	if currentGroup != nil {
		versionGroups = append(versionGroups, *currentGroup) // append the last group here
	}

	return versionGroups, nil
}

func PrintReleases(versions ...string) error {
	versionGroups, err := GetReleases()
	if err != nil {
		return err
	}
	version := ""

	if len(versions) > 0 {
		version = versions[0]
		if version[:1] != "v" {
			version = "v" + version
		}
	}

	if version == "" {
		fmt.Printf("Checking releases on %s\n", runtime.GOOS)
	} else {
		fmt.Printf("Getting releases on %s containing version %s\n", runtime.GOOS, version)
	}

	// this is for testing output
	if version == "" {
		for _, versionGroup := range versionGroups {
			fmt.Printf("Version Prefix: %s\n", versionGroup.VersionPrefix)
			for _, release := range versionGroup.Releases {
				fmt.Printf("\tVersion: %s, Release Notes: %s, Download URI: %s, SHA256Sum: %s\n", release.Version, release.ReleaseNotes, release.DownloadURI, release.SHA256Sum)
			}
		}
	} else {
		found := false
		for _, versionGroup := range versionGroups {
			if versionGroup.VersionPrefix == version {
				found = true
				fmt.Printf("| %s|%s|%s|\n", strings.Repeat("-", 16), strings.Repeat("-", 12), strings.Repeat("-", 11))
				fmt.Printf("| Version%8s | Date%6s | Available |\n", "", "")
				fmt.Printf("| %s|%s|%s\n", strings.Repeat("-", 16), strings.Repeat("-", 12), strings.Repeat("-", 12))
				for _, release := range versionGroup.Releases {
					available := "YES"
					if release.DownloadURI == "WITHDRAWN" {
						available = release.DownloadURI
					}
					fmt.Printf("| %-15s | %s | %-9s |\n", release.Version, release.Date, available)
				}
				fmt.Printf("| %s|%s|%s|\n", strings.Repeat("-", 16), strings.Repeat("-", 12), strings.Repeat("-", 11))
			}
		}
		if !found {
			fmt.Println("There do not appear to be any downloads for this version")
		}
	}
	return nil
}
