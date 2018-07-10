package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const path = "/var/www/html/img/insp/"
const url = "https://enju.us/img/insp/"

type File struct {
	Name     string  `json:"Name"`
	Url      string  `json:"Url"`
	Path     string  `json:"-"`
	Children []*File `json:"Children"`
}

func iterateJSON(w http.ResponseWriter, r *http.Request) {
	rootOSFile, _ := os.Stat(path)
	rootFile := toFile(rootOSFile, path) //start with root file
	stack := []*File{rootFile}

	for len(stack) > 0 { //until stack is empty,
		file := stack[len(stack)-1] //pop entry from stack
		stack = stack[:len(stack)-1]
		children, _ := ioutil.ReadDir(file.Path) //get the children of entry
		for _, chld := range children {          //for each child
			child := toFile(chld, filepath.Join(file.Path, chld.Name())) //turn it into a File object
			file.Children = append(file.Children, child)                 //append it to the children of the current file popped
			stack = append(stack, child)                                 //append the child to the stack, so the same process can be run again
		}
	}

	output, _ := json.Marshal(rootFile)
	w.Header().Set("Content-Type", "application/json") // set Headers for JSON and Output to ResponseWriter
	w.Write(output)
}

func toFile(file os.FileInfo, path string) *File {
	JSONFile := File{
		Name:     file.Name(),
		Path:     path,
		Url:      url + file.Name(),
		Children: []*File{},
	}
	return &JSONFile
}

func main() {
	http.HandleFunc("/", iterateJSON)
	http.ListenAndServe(":3300", nil)
}
