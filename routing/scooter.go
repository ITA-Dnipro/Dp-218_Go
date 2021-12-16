package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
)

var scooterKeyRoutes = []Route{
	{
		Uri:     `/upload`,
		Method:  http.MethodPost,
		Handler: uploadFile,
	},
	{
		Uri:     `/getScooters`,
		Method:  http.MethodGet,
		Handler: getAllScooters,
	},
	{
		Uri:     `/getScooterById/{id}`,
		Method:  http.MethodGet,
		Handler: getScooterById,
	},
	{
		Uri:     `/createScooter`,
		Method:  http.MethodPost,
		Handler: createScooter,
	},
	{
		Uri:     `/deleteScooter/{id}`,
		Method:  http.MethodDelete,
		Handler: deleteScooter,
	},
}

type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
}

type scooterWithModelList struct {
	models.Scooter
}

func AddScooterHandler(router *mux.Router, service *services.SupplierService){
	supplierService = service
	for _, rt := range scooterModelKeyRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getAllScooters(w http.ResponseWriter, r *http.Request) {
	var scooters = &models.ScooterList{}
	var err error
	format := GetFormatFromRequest(r)

	r.ParseForm()
	searchData := r.FormValue("SearchData")
	if len(searchData)==0 {
		scooters, err = supplierService.GetAllScooters()
	}
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, scooters, HTMLPath+"scooters-list.html")
}

func getScooterById(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	scooterId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	scooter, err := supplierService.GetScooterById(scooterId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, &scooterWithModelList{scooter}, HTMLPath+"user-edit.html")
}

func createScooter(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	scooter := &models.Scooter{}
	DecodeRequest(FormatJSON, w, r, scooter, nil)

	if err := supplierService.AddScooter(scooter); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, scooter)
}

func deleteScooter(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	userId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	err = supplierService.DeleteScooter(userId)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}
	EncodeAnswer(format, w, ErrorRenderer(fmt.Errorf(""), "success", http.StatusOK))
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	filepath := "./internal/"+handler.Filename
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	supplierService.InsertScootersToDb(filepath)
}
