package action

import (
	"fmt"
	"time"
)

func TaskSpinner(done chan bool) {
	spinnerChars := []rune{'|', '/', '-', '\\'}
	i := 0
	for {
		select {
		case <-done:
			// Erase the last spinner character
			fmt.Printf("\b \b")
			return
		default:
			// Print the next spinner character
			fmt.Printf("%c\b", spinnerChars[i%len(spinnerChars)])
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}
