package internal

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Reads all csv files
func ReadCSVFile(bucketName string, csvFileName string, w http.ResponseWriter, r *http.Request) [][]string {
	path := ""
	if bucketName == "" {
		path = fmt.Sprintf("%s/%s", *dir, csvFileName)
	} else {
		path = fmt.Sprintf("%s/%s/%s", *dir, bucketName, csvFileName)
	}

	file, err := os.Open(path)
	if err != nil {
		Printxml(w, r, "Error open file", 404, false)
		return nil
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		Printxml(w, r, "Error read meta data", 404, false)
		return nil
	}
	return records
}

// Writes data to a csv file
func WriteObjectsorBucket(bucketname string, nameCsvFile string, data []string, data2 [][]string, w http.ResponseWriter, r *http.Request) {
	path := ""
	if bucketname == "" {
		path = fmt.Sprintf("%s/%s", *dir, nameCsvFile)
	} else {
		path = fmt.Sprintf("%s/%s/%s", *dir, bucketname, nameCsvFile)
	}
	if data2 == nil {

		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			Printxml(w, r, "Error open objects", 409, false)
			return
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		err = writer.Write(data)
		if err != nil {
			Printxml(w, r, "Error append file", 409, false)
			return
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			Printxml(w, r, "Error append file", 409, false)
			return
		}
	} else {
		newfile, err := os.Create(path)
		if err != nil {
			Printxml(w, r, "Error delete bucket", 409, false)
		}
		writer := csv.NewWriter(newfile)
		err = writer.WriteAll(data2)
		if err != nil {
			Printxml(w, r, "Error write bucket", 409, false)
		}
		writer.Flush()
	}
}

// Accepts the text of an error or a successful attempt, and outputs it as an xml as a response to the body. Also logs to the terminal
func Printxml(w http.ResponseWriter, r *http.Request, text string, erroring int, check bool) {
	if check {
		var newmessage SuccessResponse
		newmessage.Message = text
		newxml, err := xml.MarshalIndent(newmessage, "   ", "")
		if err != nil {
			log.Println("Error marshaling error XML:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		log.Println(text)
		w.WriteHeader(erroring)
		fmt.Fprintln(w, string(newxml))

	} else {
		var err Errors
		err.Err = text
		newxml, errMarshal := xml.MarshalIndent(err, "   ", "")
		if errMarshal != nil {
			log.Println("Error marshaling error XML:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		log.Println(text)
		w.Header().Set("Content-Type", "application/xml")
		http.Error(w, string(newxml), erroring)
	}
}
