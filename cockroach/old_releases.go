package cockroach

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Releases struct {
	List   []Latest
	Latest Latest
}

type Latest struct {
	Version string
	URI     string
}

func OldGetReleases(echo bool) (string, string) {
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
					if !strings.Contains(e.Val, "alpha") && !strings.Contains(e.Val, "beta") && !strings.Contains(e.Val, "-rc") && r.Latest == (Latest{}) {
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
