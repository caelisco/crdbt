package action

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/color"
)

func Tidy() {
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
			fmt.Println("Failed to delete:", v)
		}
		color.Println("<tomato>Deleted</>:", v)
	}
	if len(files) > 0 {
		color.Println("<green>SUCCESS!</> Folder is now tidy")
	}
}
