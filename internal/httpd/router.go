package httpd

import "net/http"

type Url struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
}

var (
	urls = []Url{
		{
			Pattern: "/upload-file",
			Handler: FileUploadHandler,
		},
		{
			Pattern: "/add-replica",
			Handler: AddReplicaHandler,
		},
		{
			Pattern: "/key-store",
			Handler: HandleKeyStoreCommand,
		},
	}
)
