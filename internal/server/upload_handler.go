package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func UploadHandler(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println("Error parsing form ", err)
		return
	}
	doc, docHeader, err := req.FormFile("document")
	if err != nil {
		http.Error(resp, "message", 400)
		return
	}
	dst, err := os.Create("./uploads/" + docHeader.Filename)
	if err != nil {
		http.Error(resp, "Unable to create the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, doc)
	if err != nil {
		http.Error(resp, "Unable to save the file", http.StatusInternalServerError)
		return

	}

	defer doc.Close()

}
