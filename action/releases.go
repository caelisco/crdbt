package action

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
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
			if strings.Contains(version, "\n") {
				version = strings.Split(version, "\n")[0]
			}
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

	// change the order to reverse items

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
			fmt.Printf("\n| Version Prefix: %-38s |\n", versionGroup.VersionPrefix)
			fmt.Printf("| %s|%s|%s|\n", strings.Repeat("-", 30), strings.Repeat("-", 12), strings.Repeat("-", 11))
			fmt.Printf("| Version%22s | Date%6s | Available |\n", "", "")
			fmt.Printf("| %s|%s|%s\n", strings.Repeat("-", 30), strings.Repeat("-", 12), strings.Repeat("-", 12))
			for _, release := range versionGroup.Releases {
				available := "YES"
				if release.DownloadURI == "WITHDRAWN" {
					available = release.DownloadURI
				}
				fmt.Printf("| %-29s | %s | %-9s |\n", release.Version, release.Date, available)
			}
			fmt.Printf("| %s|%s|%s|\n", strings.Repeat("-", 30), strings.Repeat("-", 12), strings.Repeat("-", 11))
		}
	} else {
		found := false
		for _, versionGroup := range versionGroups {
			if versionGroup.VersionPrefix == version {
				found = true
				color.Printf("\n| Version Prefix: <green>%-38s</> |\n", versionGroup.VersionPrefix)
				fmt.Printf("| %s|%s|%s|\n", strings.Repeat("-", 30), strings.Repeat("-", 12), strings.Repeat("-", 11))
				fmt.Printf("| Version%22s | Date%6s | Available |\n", "", "")
				fmt.Printf("| %s|%s|%s\n", strings.Repeat("-", 30), strings.Repeat("-", 12), strings.Repeat("-", 12))
				for _, release := range versionGroup.Releases {
					available := "YES"
					if release.DownloadURI == "WITHDRAWN" {
						available = release.DownloadURI
					}
					fmt.Printf("| %-29s | %s | %-9s |\n", release.Version, release.Date, available)
				}
				fmt.Printf("| %s|%s|%s|\n", strings.Repeat("-", 30), strings.Repeat("-", 12), strings.Repeat("-", 11))
			}
		}
		if !found {
			fmt.Println("There do not appear to be any downloads for this version")
		}
	}
	return nil
}

func Interactive() (Release, error) {
	var release Release
	// ask the user to select a version
	releases, err := GetReleases()
	if err != nil {
		return release, err
	}
	// build a list of versions
	var versions []string
	for _, v := range releases {
		versions = append(versions, v.VersionPrefix)
	}
	prompt := promptui.Select{
		Label: "Select a version of Cockroach to use",
		Items: versions,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return release, err
	}
	var details []string
	for _, v := range releases {
		if v.VersionPrefix == result {
			for _, r := range v.Releases {
				details = append(details, fmt.Sprintf("%s released on %s", r.Version, r.Date))
			}
		}
	}
	fmt.Println(result)
	prompt = promptui.Select{
		Label: "Select a release",
		Items: details,
	}

	_, result, err = prompt.Run()
	if err != nil {
		return release, err
	}

	for _, r := range releases {
		for _, v := range r.Releases {
			if v.Version == result {
				release = v
			}
		}
	}

	return release, nil
}

func InteractiveVersion(versionGroup []VersionGroup) (string, error) {
	// build a list of versions
	var versions []string
	for _, v := range versionGroup {
		versions = append(versions, v.VersionPrefix)
	}
	prompt := promptui.Select{
		Label: "Select a version of Cockroach to use",
		Items: versions,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func InteractiveRelease(version string, versionGroup []VersionGroup) (Release, error) {
	var release Release
	var details []string
	for _, v := range versionGroup {
		if v.VersionPrefix == version {
			for _, r := range v.Releases {
				details = append(details, fmt.Sprintf("%s released on %s", r.Version, r.Date))
			}
		}
	}

	prompt := promptui.Select{
		Label: "Select a release",
		Items: details,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return release, err
	}
	// scan through the list to get the details
	fmt.Println("Searching for:", strings.Split(result, " ")[0])
	for _, v := range versionGroup {
		for _, r := range v.Releases {
			if r.Version == strings.Split(result, " ")[0] {
				release = r
			}
		}
	}
	return release, nil
}
