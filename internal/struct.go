package internal

import "flag"

type Bucket struct {
	Name             string `xml:"bucketname"`
	CreationTime     string `xml:"creationtime"`
	LastModifiedTime string `xml:"lastmodifiedtime"`
	Status           string `xml:"status"`
}

type Buckets struct {
	BucketIn []Bucket `xml:"Bucket"`
}

type Object struct {
	NameObject             string `xml:"name"`
	Size                   string `xml:"size"`
	ContenType             string `xml:"content type"`
	LastModifiedTimeObject string `xml:"last modififiedtime"`
}

type Objects struct {
	ObjectIn []Object `xml:"Oject"`
}

type Errors struct {
	Err string `xml:"error"`
}

type SuccessResponse struct {
	Message string `xml:"message"`
}

var (
	dir  = flag.String("dir", "data", "Path to the directory")
	port = flag.String("port", "8080", "Port number")
)
