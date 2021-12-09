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
		Uri:     `/createScooterModel`,
		Method:  http.MethodPost,
		Handler: createScooterModel,
	},
	{
		Uri:     `/getScooterModels`,
		Method:  http.MethodGet,
		Handler: getScooterModels,
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
		Uri:     `/getScootersByModelId/{id}`,
		Method:  http.MethodGet,
		Handler: getScooterByModelId,
	},
	{
		Uri:     `/createScooter`,
		Method:  http.MethodPost,
		Handler: createScooter,
	},
	{
		Uri:     `/updateScooter`,
		Method:  http.MethodPut,
		Handler: updateScooterSerial,
	},
	{
		Uri:     `/createScooter`,
		Method:  http.MethodDelete,
		Handler: deleteScooter,
	},
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
	searchData := r.FormValue("SearchData")
	if len(searchData)==0 {
		modelList, err = supplierService.GetScooterModels()
	}

	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, modelList, HTMLPath+"user-list.html")
}

func getScooterModelById(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	modelId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	scooterModel, err := supplierService.GetScooterModelById(modelId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, &scooterModel, HTMLPath+"user-edit.html")
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

	EncodeAnswer(format, w, scooters, HTMLPath+"user-list.html")
}

func getScooterByModelId(w http.ResponseWriter, r *http.Request) {
	var modelList = &models.ScooterList{}
	var err error
	format := GetFormatFromRequest(r)

	r.ParseForm()
	modelId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	modelList, err = supplierService.GetScooterByModelId(modelId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, modelList, HTMLPath+"user-list.html")
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


func updateScooterSerial(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	scooterId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
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

	EncodeAnswer(format, w, &scooterData, HTMLPath+"user-edit.html")
}

func deleteScooter(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	userId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
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
