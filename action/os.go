package action

import (
	"fmt"
	"runtime"
)

func GetOS() string {
	return fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
}
