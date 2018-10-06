package commands

import (
	"os/user"
	"path/filepath"
)

func expandPath(path string) string {
	user, _ := user.Current()
	var expandedPath string

	if path[0] == '~' {
		expandedPath = filepath.Join(user.HomeDir, path[1:])
	} else {
		expandedPath, _ = filepath.Abs(path)
	}

	return expandedPath
}
