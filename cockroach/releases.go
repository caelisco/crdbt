package cockroach

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Release struct {
	Version      string
	ReleaseNotes string
	DownloadURI  string
	SHA256Sum    string
}

type VersionGroup struct {
	VersionPrefix string
	Releases      []Release
}

func GetReleases(versions ...string) ([]VersionGroup, error) {
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

	fmt.Println("running with scope:", runtime.GOOS)

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

	// this is for testing output
	if len(versions) > 0 {
		if versions[0][:1] != "v" {
			versions[0] = "v" + versions[0]
		}
		for _, versionGroup := range versionGroups {
			if versionGroup.VersionPrefix == versions[0] {
				fmt.Printf("Version Prefix: %s\n", versionGroup.VersionPrefix)
				fmt.Printf("Version%8s | Download%70s | SHA256%75s | Release Notes%56s\n", "", "", "", "")
				fmt.Printf("%s|%s|%s|%s\n", strings.Repeat("-", 16), strings.Repeat("-", 80), strings.Repeat("-", 83), strings.Repeat("-", 70))
				for _, release := range versionGroup.Releases {
					fmt.Printf("%-15s | %-78s | %-81s | %-55s\n", release.Version, release.DownloadURI, release.SHA256Sum, release.ReleaseNotes)
				}
			}
		}
	} else {
		for _, versionGroup := range versionGroups {
			fmt.Printf("Version Prefix: %s\n", versionGroup.VersionPrefix)
			for _, release := range versionGroup.Releases {
				fmt.Printf("\tVersion: %s, Release Notes: %s, Download URI: %s, SHA256Sum: %s\n", release.Version, release.ReleaseNotes, release.DownloadURI, release.SHA256Sum)
			}
		}
	}

	return versionGroups, nil
}
