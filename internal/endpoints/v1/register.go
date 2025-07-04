package v1

import "github.com/gorilla/mux"

func RegisterAPIRoutes(router *mux.Router) {
	RegisterDirectoryRoutes(router)
	RegisterFileRoutes(router)
	RegisterUploadRoutes(router)
	RegisterDownloadRoutes(router)
}
