package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// The main function that checks everything and adds a new object.
func CreateObject(path []string, w http.ResponseWriter, r *http.Request) {
	if len(path) != 2 {
		Printxml(w, r, "Not enough information", 409, false)
		return
	}
	if !BucketCheck(path[0]) {
		Printxml(w, r, "Bucket is missing", 409, false)
		return
	}
	newPath := fmt.Sprintf("%s/%s/%s", *dir, path[0], path[1])
	_, err := os.Create(newPath)
	if err != nil {
		Printxml(w, r, "Error create file", 404, false)
		return
	}
	writeBody := r.Body
	file, err := io.ReadAll(writeBody)
	if err != nil {
		Printxml(w, r, "Error read new file", 404, false)
		return
	}
	err1 := os.WriteFile(newPath, file, 0o664)

	if err1 != nil {
		Printxml(w, r, "Error write new file", 404, false)
		return
	}
	objectKey := path[1]
	size := strconv.Itoa(len(file))
	contentType := r.Header.Get("Content-type")
	nowTime := time.Now()
	LastModifie := nowTime.Format("2006-01-02 15:04:05")
	data := []string{objectKey, size, contentType, LastModifie}
	newPath1, check := EditObject(w, r, path)
	if check {
		WriteObjectsorBucket(path[0], "objects.csv", nil, newPath1, w, r)
	}

	WriteObjectsorBucket(path[0], "objects.csv", data, nil, w, r)
	content := EditCSVFilePutObject(w, r, path)
	if content == nil {
		Printxml(w, r, "Error read file", 404, false)
		return
	}
	WriteObjectsorBucket("", "metadata.csv", nil, content, w, r)

	Printxml(w, r, "Put method done successfully", 200, true)
}

// Changes information about an object in a csv file
func EditObject(w http.ResponseWriter, r *http.Request, path []string) ([][]string, bool) {
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

	return csvList, check
}

// Changes bucket information in csv file
func EditCSVFilePutObject(w http.ResponseWriter, r *http.Request, path []string) [][]string {
	content := ReadCSVFile("", "metadata.csv", w, r)
	check := false
	for _, record := range content {
		if record[0] == path[0] {
			check = true
			nowTime := time.Now()
			LastModifie := nowTime.Format("2006-01-02 15:04:05")
			record[2] = LastModifie
		}
	}
	if !check {
		return nil
	}
	return content
}
