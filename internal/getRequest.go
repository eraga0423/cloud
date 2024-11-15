package internal

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// Checks everything and outputs the response to the general GET request
func (b Buckets) GetHandler(w http.ResponseWriter, r *http.Request) {
	content := ReadCSVFile("", "metadata.csv", w, r)
	if content == nil {
		return
	}
	for _, record := range content {
		if 4 > len(record) {
			Printxml(w, r, "Not enough information", 409, false)
			return
		}

		b.BucketIn = append(b.BucketIn, Bucket{
			Name:             record[0],
			CreationTime:     record[1],
			LastModifiedTime: record[2],
			Status:           record[3],
		})

	}
	data, err := xml.MarshalIndent(b, " ", "   ")
	if err != nil {
		Printxml(w, r, "Error buckets marshaling xml", 500, false)
		return
	}
	fmt.Fprintln(w, string(data))
	Printxml(w, r, "Get method done successfully", 200, true)
}
