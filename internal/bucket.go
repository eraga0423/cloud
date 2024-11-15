package internal

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Creates a bucket. This is where the bucket is checked for any errors.
func CreateBucket(w http.ResponseWriter, r *http.Request, bucket string) {
	flag.Parse()
	switch {
	case !DuplicateChecker(w, r, bucket):
		Printxml(w, r, "Error duplicate names", 409, false)
		return
	case !BucketCheck(bucket):
		Printxml(w, r, "Error invalid names", 400, false)
		return
	default:
		temp := fmt.Sprintf("%s/%s", *dir, bucket)
		err := os.MkdirAll(temp, 0o755)
		if err != nil {
			Printxml(w, r, "Error creating bucket", 500, false)
			return
		}
		nowTime := time.Now()
		timeString := nowTime.Format("2006-01-02 15:04:05")
		data := []string{bucket, timeString, timeString, "active"}
		WriteObjectsorBucket("", "metadata.csv", data, nil, w, r)
		CreateObjectsList(bucket)
		Printxml(w, r, "Bucket created succesfully", 200, true)
		return
	}
}

// Checks bucket validity
func BucketCheck(bucket string) bool {
	if len(bucket) > 63 || len(bucket) < 3 {
		return false
	}
	re := regexp.MustCompile(`^[a-z0-9]+([.-][a-z0-9]+)*$`)
	ipre := regexp.MustCompile(`^\d{1,3}(\.\d{1,3}){3}$`)
	return re.MatchString(bucket) && !ipre.MatchString(bucket)
}

func DuplicateChecker(w http.ResponseWriter, r *http.Request, bucket string) bool {
	dir := ReadCSVFile("", "metadata.csv", w, r)
	for _, v := range dir {
		if v[0] == bucket {
			return false
		}
	}
	return true
}
