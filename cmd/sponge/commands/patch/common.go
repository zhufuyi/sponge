package patch

import (
	"os"
	"strings"
)

// get moduleName and serverName from directory
func getNamesFromOutDir(dir string) (moduleName string, serverName string) {
	if dir == "" {
		return "", ""
	}
	data, err := os.ReadFile(dir + "/docs/gen.info")
	if err != nil {
		return "", ""
	}

	ms := strings.Split(string(data), ",")
	if len(ms) != 2 {
		return "", ""
	}

	return ms[0], ms[1]
}
