package mapnik

import (
	"os/exec"
	"strings"
)

const _mapnikConfigCMD = "mapnik-config"

func ConfigFonts() string {
	out, _ := exec.Command(_mapnikConfigCMD, "--fonts").Output()

	return strings.TrimSpace(string(out))
}

func ConfigPlugins() string {
	out, _ := exec.Command(_mapnikConfigCMD, "--input-plugins").Output()

	return strings.TrimSpace(string(out))
}
