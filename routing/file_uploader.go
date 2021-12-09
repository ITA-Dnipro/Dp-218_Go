package routing

import (
	"Dp218Go/services"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var fileService *services.FileService

var fileKeyRoutes = []Route{
	{
		Uri:     `/uploadFile`,
		Method:  http.MethodPost,
		Handler: uploadFile,
	},
}

func AddFileHandler(router *mux.Router, service *services.FileService) {
	fileService = service
	for _, rt := range fileKeyRoutes{
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request){
	fmt.Println("File Upload Endpoint Hit")

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	filePath := fileService.InsertScootersToDb(file)
	defer os.Remove(filePath)

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

