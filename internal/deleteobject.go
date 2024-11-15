package internal

import (
	"fmt"
	"net/http"
	"os"
)

// The main function that checks objects before deleting them and deletes them also sends a response.
func BasicDeleteObjects(path []string, w http.ResponseWriter, r *http.Request) {
	newContent := Deleteobject(w, r, path)
	if newContent == nil {
		Printxml(w, r, "This file does not exist", 409, false)
		return
	}
	deletePath := fmt.Sprintf("%s/%s/%s", *dir, path[0], path[1])
	err := os.Remove(deletePath)
	if err != nil {
		Printxml(w, r, "Error deleting file", 409, false)
	}
	WriteObjectsorBucket(path[0], "objects.csv", nil, newContent, w, r)
	content := EditCSVFilePutObject(w, r, path)
	if content == nil {
		Printxml(w, r, "Error read file", 404, false)
		return
	}
	WriteObjectsorBucket("", "metadata.csv", nil, content, w, r)
	Printxml(w, r, "Succes deleting file", 204, true)
}

// Changes data in csv file
func Deleteobject(w http.ResponseWriter, r *http.Request, path []string) [][]string {
	csvList := ReadCSVFile(path[0], "objects.csv", w, r)
	check := false
	for i, v := range csvList {
		if len(v) < 4 {
			Printxml(w, r, "Not enough information", 409, false)
		}
		if v[0] == path[1] {
			csvList = append(csvList[:i], csvList[i+1:]...)
			check = true
		}
	}
	if !check {
		return nil
	}
	return csvList
}
