package internal

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// Checked everything and displays the response to the GET request of the object itself
func (o Objects) GetObjects(path []string, w http.ResponseWriter, r *http.Request) {
	content := ReadCSVFile(path[0], "objects.csv", w, r)
	if content == nil {
		return
	}
	checkExisting := false
	for _, record := range content {
		if len(record) < 4 {
			Printxml(w, r, "Not enough information", 409, false)
			return
		}
		if len(path) == 2 {
			if record[0] == path[1] {
				checkExisting = true
				o.ObjectIn = append(o.ObjectIn, Object{
					NameObject:             record[0],
					Size:                   record[1],
					ContenType:             record[2],
					LastModifiedTimeObject: record[3],
				})
			}
		}

	}
	if !checkExisting {
		Printxml(w, r, "Non-existent file", 409, false)
		return
	}
	data, err := xml.MarshalIndent(o, " ", "   ")
	if err != nil {
		Printxml(w, r, "Error objects marshaling xml", 500, false)
	}
	fmt.Fprintln(w, string(data))
	Printxml(w, r, "Get objects method done successfully", 200, true)
}
