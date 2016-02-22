package ico_test

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zyxar/image2ascii/ico"
)

func TestDecodeAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	files, _ := filepath.Glob("testdata/favicons/*.ico")
	for _, f := range files {
		rd, err := os.Open(f)
		assert.NoError(err, f)
		images, err := ico.DecodeAll(rd)
		assert.NoError(err, f)
		rd.Close()
		if err != nil {
			continue
		}

		for i := range images {
			var dst string
			if len(images) == 1 {
				dst = f + ".jpg"
			} else {
				dst = f + fmt.Sprintf("-%d.jpg", i)
			}
			rd, err := os.Open(dst)
			assert.NoError(err, dst)
			dstImage, err := jpeg.Decode(rd)
			assert.NoError(err, dst)
			rd.Close()
			if err != nil {
				continue
			}
			assert.Equal(images[i].Bounds(), dstImage.Bounds())
		}
	}
}
