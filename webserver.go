// main.go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const uploadPath = "/munki_repo/pkgs"

func uploadFile(w http.ResponseWriter, r *http.Request) {

	//
	// Get the file from the form or print and error
	//
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	//
	// Read the file into a byte array
	//
	fileBytes, err := ioutil.ReadAll(file)

	subdir := r.FormValue("subdir")

	if _, err := os.Stat(filepath.Join(uploadPath, subdir)); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(uploadPath, subdir), 0744)
	}

	newPath := filepath.Join(uploadPath, subdir, handler.Filename)

	newFile, err := os.Create(newPath)
	if err != nil {
		return
	}
	defer newFile.Close()

	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		return
	}
}

func main() {

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/", fileServer)
	http.Handle("/repo/", http.StripPrefix("/repo/", http.FileServer(http.Dir("/munki_repo"))))
	http.HandleFunc("/upload", uploadFile)

	http.ListenAndServe(":8080", nil)
}
