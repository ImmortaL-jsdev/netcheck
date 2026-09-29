package zapret

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Zapret struct {
	Path    string
	Version string
}

func Detect() (*Zapret, error) {
	var path string

	switch runtime.GOOS {
	case "linux":
		path = "/opt/zapret"
	case "windows":
		path = `C:\zapret`
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("zapret not found at %s", path)
	}
	version := "unknown"
	versionData, err := os.ReadFile(filepath.Join(filepath.Dir(path), "zapret-ver"))
	if err == nil {
		version = strings.TrimSpace(string(versionData))
	}
	return &Zapret{Path: path, Version: version}, nil
}

func (z *Zapret) IsRunning() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	cmd := exec.Command("systemctl", "is-active", "zapret")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "active"
}
func (z *Zapret) HasBinary(name string) bool {
	paths := []string{
		filepath.Join(z.Path, "nfq", name),
		filepath.Join(z.Path, "binaries", name),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return true
		}

	}
	return false
}
