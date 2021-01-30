package main

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	//"io"
	"io/ioutil"
	"os"

	// depends on libvips-dev
	"github.com/daddye/vips"
	"github.com/julienschmidt/httprouter"
)

var breakpoints = map[string]breakpoint{
	"xs": {200, 75},
	"sm": {600, 75},
	"md": {960, 80},
	"lg": {1280, 90},
	"xl": {1920, 95},
}

type breakpoint struct{
	size int
	quality int
}


func uploadImage(db map[string]string, dir string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		//r.ParseMultipartForm(32 << 20)

		file, handler, err := r.FormFile("nt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		path := WriteFileToUUIDPath(file, dir)
		db[handler.Filename] = path
	}
}

func handleImage(db map[string]string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		img := params.ByName("img")

		sizeIndex := strings.LastIndex(img, "-")
		ext := path.Ext(img)

		size := strings.Trim(img[sizeIndex+1:], ext)
		if _, ok := breakpoints[size]; ok {
			img = img[:sizeIndex] + ext
		}

		p := db[img]


		f, err := os.Open(p)
		if err != nil {
			fmt.Println(err)
		}

		options := vips.Options{
			Width:        breakpoints[size].size,
			Height:       breakpoints[size].size,
			Crop:         false,
			Extend:       vips.EXTEND_WHITE,
			Interpolator: vips.BILINEAR,
			Gravity:      vips.CENTRE,
			Quality:      breakpoints[size].quality,
		}
		inBuf, _ := ioutil.ReadAll(f)
		buf, err := vips.Resize(inBuf, options)
		if err != nil {
			fmt.Println(err)
			return
		}

		contentType := http.DetectContentType(buf)
		fmt.Println(contentType)

		w.Header().Set("Content-Type", contentType)
		w.Write(buf)
	}
}