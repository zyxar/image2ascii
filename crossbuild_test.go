package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_GOARCH = "GOARCH=amd64"
)

func build(goos string) error {
	cmd := exec.Command("go", "build", "-o", "image2ascii."+goos)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS="+goos)
	cmd.Env = append(cmd.Env, _GOARCH)
	return cmd.Run()
}

func TestCrossBuild(t *testing.T) {
	t.Parallel()
	ast := assert.New(t)
	buf, err := ioutil.ReadFile("ascii/term.go")
	ast.NoError(err)
	n := bytes.IndexByte(buf, '\n')
	if n == -1 || n <= 9 {
		t.Fatal("invalid content of term.go")
	}
	buf = buf[9:n]
	ss := bytes.Split(buf, []byte{' '})
	for i, _ := range ss {
		goos := string(ss[i])
		err = build(goos)
		ast.NoError(err, goos)
		if err != nil {
			println(err.Error())
		}
		os.Remove("image2ascii." + goos) // cleanup
	}
}
