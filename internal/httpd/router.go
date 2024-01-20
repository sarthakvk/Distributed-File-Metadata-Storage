// This file defines the routes in a slice for the applications, the design is inspired from Djnago

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
