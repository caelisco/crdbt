package action

import (
	"fmt"

	"github.com/gookit/color"
)

func Install(file string) error {
	var err error
	if !FileExists(file) {
		file, err = Download(file)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	if file != "" {
		color.Printf("<cyan>Info:</> File exists <yellow>%s</>\n", file)
		err = ExtractTGZ(file)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	// move files
	err = MoveFiles(file)
	if err != nil {
		return err
	}
	// give 'next step' instructions
	return nil
}
