package internal

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// The main function that checks before deleting a bucket and deletes it also sends a response.
func DeleteHandler(w http.ResponseWriter, r *http.Request, deleteBucket string) {
	results, check := EditFile(w, r, deleteBucket)
	if results == nil {
		return
	}
	WriteObjectsorBucket("", "metadata.csv", nil, results, w, r)
	deletePath := fmt.Sprintf("%s/%s", *dir, deleteBucket)
	if check {
		Printxml(w, r, "Conflict for a non-empty bucket", 409, false)
		return
	} else if !check {
		Printxml(w, r, "The deletion was successful.", 204, true)
	}
	err := os.RemoveAll(deletePath)
	if err != nil {
		Printxml(w, r, "Error bucket not deleted", 404, false)
		return
	}
}

// Checks if there are files in the bucket
func CheckEmptyFile(w http.ResponseWriter, r *http.Request, deleteBucket string) bool {
	path := fmt.Sprintf("%s/%s", *dir, deleteBucket)
	entries, err := os.ReadDir(path)
	if err != nil {
		Printxml(w, r, "Error read directory", 409, false)
	}
	if len(entries) == 1 {
		return true
	}
	return false
}

// Changes data in csv file
func EditFile(w http.ResponseWriter, r *http.Request, deleteBucket string) ([][]string, bool) {
	csvList := ReadCSVFile("", "metadata.csv", w, r)
	check := false
	checkDeleteName := false
	for i, v := range csvList {
		if len(v) < 4 {
			Printxml(w, r, "Not enough information", 409, false)
			return nil, false
		}
		if v[0] == deleteBucket {
			checkDeleteName = true
			if CheckEmptyFile(w, r, deleteBucket) {
				csvList = append(csvList[:i], csvList[i+1:]...)
				continue
			}
			check = true
			nowTime := time.Now()
			timeString := nowTime.Format("2006-01-02 15:04:05")
			v[3] = "marked for deletion"
			v[2] = timeString
		}
	}
	if !checkDeleteName {
		Printxml(w, r, "No delete name", 404, false)
		return nil, false
	}
	return csvList, check
}
