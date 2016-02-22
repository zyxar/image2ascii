package term_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_GOARCH = "GOARCH=amd64"
)

var OSs = []string{"darwin", "dragonfly", "freebsd", "linux", "netbsd", "openbsd", "windows"}

func build(goos string) error {
	cmd := exec.Command("go", "build")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS="+goos)
	cmd.Env = append(cmd.Env, _GOARCH)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func TestCrossBuild(t *testing.T) {
	t.Parallel()
	ast := assert.New(t)
	for i := range OSs {
		err := build(OSs[i])
		ast.NoError(err, OSs[i])
	}
}
