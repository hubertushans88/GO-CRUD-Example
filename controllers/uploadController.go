package controllers

import (
	"bytes"
	"fmt"
	resizer "github.com/nfnt/resize"
	"github.com/pkg/errors"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	alias := r.FormValue("alias")

	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imageBuf, err := ToPng(buf.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, err := png.Decode(bytes.NewReader(imageBuf))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buff := new(bytes.Buffer)
	err = png.Encode(buff, img)
	if err != nil {
		fmt.Println("failed to create buffer", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := handler.Filename
	if alias != "" {
		//filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
		filename = fmt.Sprintf("%s%s", alias, ".png")
	}

	fileLocation := filepath.Join(dir, "files", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	//if _, err := io.Copy(targetFile, uploadedFile); err != nil {
	if _, err := io.Copy(targetFile, buff); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("done : " + filename))

}

func ViewImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := r.FormValue("id")
	//params := mux.Vars(r)
	//id:= params["id"]

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileLocation := filepath.Join(dir, "files", id+".png")
	ResizeImage(w, fileLocation, 500)

}

func ToPng(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)

	switch contentType {
	case "image/png":
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode jpeg")
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return nil, errors.Wrap(err, "unable to encode png")
		}

		return buf.Bytes(), nil
	}

	return nil, fmt.Errorf("unable to convert %#v to png", contentType)
}

func ResizeImage(w io.Writer, Path string, Width uint) {
	//var ImageExtension = strings.Split(Path, ".png")
	//var ImageNum       = strings.Split(ImageExtension[0], "/")
	//var ImageName      = ImageNum[len(ImageNum)-1]
	//fmt.Println(ImageName)
	file, err := os.Open(Path)
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resizer.Resize(Width, 0, img, resizer.Lanczos3)

	jpeg.Encode(w, m, nil) // Write to the ResponseWriter
}
