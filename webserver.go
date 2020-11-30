// main.go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const basePath = "/munki_repo"

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

	//
	// Grab values from the form
	//
	file_type_path := r.FormValue("upload_type")
	subdir := r.FormValue("subdir")
	full_folder_path := filepath.Join(basePath, file_type_path, subdir)
	//
	// Check if path exists, if it doesnt, create it and all children
	//
	if _, err := os.Stat(full_folder_path); os.IsNotExist(err) {
		os.MkdirAll(full_folder_path, 0744)
	}

	//
	// Now that we know the path exists, create the full file path and create the file
	//
	full_file_path := filepath.Join(full_folder_path, handler.Filename)
	newFile, err := os.Create(full_file_path)
	if err != nil {
		return
	}
	defer newFile.Close()

	//
	// Write our byte array to the file we created
	//
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		return
	}
}

func main() {

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.Handle("/repo/", http.StripPrefix("/repo/", http.FileServer(http.Dir("/munki_repo"))))
	http.HandleFunc("/upload", uploadFile)

	http.ListenAndServe(":8080", nil)
}
