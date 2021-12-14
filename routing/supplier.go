package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var supplierService *services.SupplierService

var scooterModelKeyRoutes = []Route{
	{
		Uri:     `/getScooterModels`,
		Method:  http.MethodGet,
		Handler: getScooterModels,
	},
	{
		Uri:     `/createScooterModel`,
		Method:  http.MethodPost,
		Handler: createScooterModel,
	},
	{
		Uri:     `/getScooterModel/{id}`,
		Method:  http.MethodGet,
		Handler: getScooterModelById,
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
		Uri:     `/editScooter/{id}`,
		Method:  http.MethodPut,
		Handler: editScooter,
	},
	{
		Uri:     `/deleteScooter/{id}`,
		Method:  http.MethodDelete,
		Handler: deleteScooter,
	},
}

type scooterWithModelList struct {
	models.Scooter
}

func AddSupplierHandler(router *mux.Router, service *services.SupplierService){
	supplierService = service
	for _, rt := range scooterModelKeyRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func createScooterModel(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	model := &models.ScooterModel{}
	DecodeRequest(FormatJSON, w, r, model, nil)

	if err := supplierService.CreateScooterModel(model); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, model)
}

func getScooterModels(w http.ResponseWriter, r *http.Request) {
	var modelList = &models.ScooterModelList{}
	var err error
	format := GetFormatFromRequest(r)

	r.ParseForm()

	modelList, err = supplierService.GetScooterModels()
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, modelList, HTMLPath+"scooter-model-list.html")
}

func getScooterModelById(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	fmt.Println(mux.Vars(r))
	modelId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	scooterModel, err := supplierService.GetScooterModelById(modelId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, &scooterModel, HTMLPath+"scooter-model-edit.html")
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

func editScooter(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	scooterId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	scooterData, err := supplierService.GetScooterByID(scooterId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	DecodeRequest(format, w, r, &scooterData, DecodeUserUpdateRequest)
	scooterData, err = supplierService.UpdateScooter(scooterId, scooterData)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, &scooterData, HTMLPath+"scooter-edit.html")
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
