package internal

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// The main function that receives requests and responds to them also parses them into methods
func MeHandler(w http.ResponseWriter, r *http.Request) {
	structsBucket := Buckets{}
	structsObject := Objects{}
	urlPath := strings.Trim(r.URL.Path, "/")
	var newPath []string
	if urlPath != "" {
		newPath = strings.Split(urlPath, "/")
	}
	err := CheckPath(w, r, newPath)
	if err != nil {
		return
	}
	switch r.Method {
	case http.MethodGet:
		if len(newPath) == 0 {
			structsBucket.GetHandler(w, r)
		} else if len(newPath) == 2 {
			structsObject.GetObjects(newPath, w, r)
		} else {
			Printxml(w, r, "Wrong size path", 409, false)
		}
	case http.MethodDelete:
		if len(newPath) == 2 {
			BasicDeleteObjects(newPath, w, r)
		} else {
			DeleteHandler(w, r, newPath[0])
		}
	case http.MethodPut:

		if len(newPath) == 2 {
			CreateObject(newPath, w, r)
		} else {
			CreateBucket(w, r, newPath[0])
		}
	default:
		Printxml(w, r, "Method error or wrong size path", 404, false)
	}
}

// Changed help flag
func init() {
	flag.Usage = func() {
		fmt.Println(
			`Simple Storage Service.

**Usage:**
triple-s [-port <N>] [-dir <S>]  
triple-s --help
		
**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`)
	}
}

// When the server starts, it creates a directory,also creates a csv file inside it
func CreateDirectory() {
	os.Mkdir(*dir, 0o755)
	option := os.O_CREATE | os.O_EXCL | os.O_WRONLY

	path1 := fmt.Sprintf("%s/metadata.csv", *dir)
	file, err1 := os.OpenFile(path1, option, 0o644)
	defer file.Close()
	if err1 != nil {
		if os.IsExist(err1) {
			log.Println("File already exists")
		} else {
			fmt.Println("Error creating file:", err1)
		}
		return
	}
}

// When creating a bucket, creates a csv file inside the bucket
func CreateObjectsList(bucketName string) {
	option := os.O_CREATE | os.O_EXCL | os.O_WRONLY
	path := fmt.Sprintf("%s/%s/objects.csv", *dir, bucketName)
	filename, err := os.OpenFile(path, option, 0o644)
	defer filename.Close()
	if err != nil {
		if !os.IsExist(err) {
			fmt.Println("Error creating file:", err)
		}
		return
	}
}

func CheckPath(w http.ResponseWriter, r *http.Request, path []string) error {
	lenpath := len(path)
	switch {
	case r.Method != http.MethodDelete && r.Method != http.MethodPut && r.Method != http.MethodGet:
		Printxml(w, r, "Unidentified method", 404, false)
		return errors.New("Unidentified method")
	case lenpath > 2:
		Printxml(w, r, "Incorrect path", 409, false)
		return errors.New("Incorrect path")
	case r.Method == http.MethodGet && (lenpath != 2 && lenpath != 1 && lenpath != 0):
		Printxml(w, r, "Incorrect path want /{bucketName}/{objectName} or /", 409, false)
		return errors.New("Incorrect path want /{bucketName}/{objectName} or /")
	case (r.Method == http.MethodDelete || r.Method == http.MethodPut) && (lenpath != 1 && lenpath != 2):
		Printxml(w, r, "Incorrect path want /{bucketName}/{objectName} or /{bucketName}", 409, false)
		return errors.New("Incorrect path want /{bucketName}/{objectName} or /{bucketName}")
	}

	return nil
}

// Checks the address and methods
func Server() {
	flag.Parse()
	log.Println("Start server")
	log.Printf("Port:%s", *port)
	CreateDirectory()
	http.HandleFunc("/", MeHandler)
	port := fmt.Sprintf(":%s", *port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
