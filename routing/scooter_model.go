package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var supplierService *services.SupplierService

var scooterModelKeyRoutes = []Route{
	{
		Uri:     `/models`,
		Method:  http.MethodGet,
		Handler: getModels,
	},
	{
		Uri:     `/createModel`,
		Method:  http.MethodPost,
		Handler: createModel,
	},
	{
		Uri:     `/editPrice`,
		Method:  http.MethodPut,
		Handler: editModel,
	},
}

func AddScooterModelHandler(router *mux.Router, service *services.SupplierService){
	supplierService = service
	for _, rt := range scooterModelKeyRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getModels(w http.ResponseWriter, r *http.Request) {
	var modelList = &models.ModelODTList{}
	var err error
	format := GetFormatFromRequest(r)

	r.ParseForm()

	modelList, err = supplierService.GetModels()
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, modelList, HTMLPath+"supplier.html")
}

func createModel(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	model := &models.ModelODT{}

	DecodeRequest(FormatJSON, w, r, model, decodePriceRequest)

	if err := supplierService.AddModel(model); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, model)
}

func editModel(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	modelId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	modelData, err := supplierService.SelectModel(modelId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	DecodeRequest(FormatJSON, w, r, modelData, decodePriceRequest)
	if err := supplierService.ChangePrice(modelData); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, modelData)
}

func decodePriceRequest(r *http.Request, data interface{}) error {
	r.ParseForm()
	modelData := data.(*models.ModelODT)

	if _, ok := r.Form["price"]; ok {
		modelPrice, err := strconv.Atoi(r.FormValue("price"))
		if err != nil {
			return err
		}
		modelData.Price = modelPrice
	}
	if _, ok := r.Form["modelName"]; ok {
		modelData.ModelName= r.FormValue("modelName")
	}
	if _, ok := r.Form["maxWeight"]; ok {
		modelMaxWeight, err := strconv.Atoi(r.FormValue("maxWeight"))
		if err != nil {
			return err
		}
		modelData.MaxWeight = modelMaxWeight
	}
	if _, ok := r.Form["speed"]; ok {
		modelSpeed, err := strconv.Atoi(r.FormValue("speed"))
		if err != nil {
			return err
		}
		modelData.MaxWeight = modelSpeed
	}

	return nil
}