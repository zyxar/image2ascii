package ico

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	files, _ := filepath.Glob("testdata/favicons/*.ico")
	for i, f := range files {
		fmt.Println(i, "WORKING WITH", f)
		icoData, err := ioutil.ReadFile(f)
		assert.NoError(err, f)

		r := bytes.NewReader(icoData)
		images, err := DecodeAll(r)
		assert.NoError(err, f)
		if err != nil {
			continue
		}

		for i, im := range images {
			var jpgName string
			if len(images) == 1 {
				jpgName = f + ".jpg"
			} else {
				jpgName = f + fmt.Sprintf("-%d.jpg", i)
			}
			jpgData, err := ioutil.ReadFile(jpgName)
			assert.NoError(err, jpgName)

			r = bytes.NewReader(jpgData)
			jpgImage, err := jpeg.Decode(r)
			assert.NoError(err, jpgName)
			if err != nil {
				continue
			}

			assert.Equal(im.Bounds(), jpgImage.Bounds())
		}
	}
}
