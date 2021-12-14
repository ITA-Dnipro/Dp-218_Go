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

var ODTService *services.ODTService

var ODTKeyRoutes = []Route{
	{
		Uri:     `/supplier`,
		Method:  http.MethodGet,
		Handler: getSuppliersPage,
	},
	{
		Uri:     `/upload`,
		Method:  http.MethodPost,
		Handler: uploadFile,
	},
	{
		Uri:     `/supplier`,
		Method:  http.MethodGet,
		Handler: getModels,
	},
	{
		Uri:     `/modelODT/{id}`,
		Method:  http.MethodGet,
		Handler: selectModel,
	},
	{
		Uri:     `/createModel`,
		Method:  http.MethodPost,
		Handler: createModel,
	},
	{
		Uri:     `/editModel/{id}`,
		Method:  http.MethodPut,
		Handler: editModel,
	},
}

type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
}

func AddODTHandler(router *mux.Router, service *services.ODTService){
	ODTService = service
	for _, rt := range ODTKeyRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getModels(w http.ResponseWriter, r *http.Request) {
	var modelList = &models.ModelODTList{}
	var err error
	format := GetFormatFromRequest(r)

	r.ParseForm()

	modelList, err = ODTService.GetModels()
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, modelList, HTMLPath+"supplier.html")
}

func selectModel(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	fmt.Println(mux.Vars(r))
	modelId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	scooterModel, err := ODTService.SelectModel(modelId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, &scooterModel, HTMLPath+"scooter-model-edit.html")
}

func createModel(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	model := &models.ModelODT{}
	DecodeRequest(FormatJSON, w, r, model, nil)

	if err := ODTService.AddModel(model); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, model)
}

func editModel(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	model := models.ModelODT{}

	DecodeRequest(FormatJSON, w, r, model, nil)

	if err :=  ODTService.EditModel(model); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, model)
}

func DecodeModelRequest(r *http.Request, data interface{}) error {
	r.ParseForm()
	modelData := data.(*models.ModelODT)

	if _, ok := r.Form["Price"]; ok {
		modelPrice, err := strconv.Atoi(r.FormValue("Price"))
		if err != nil {
			return err
		}
		modelData.Price = modelPrice
	}
	if _, ok := r.Form["ModelName"]; ok {
		modelData.ModelName= r.FormValue("UserName")
	}
	if _, ok := r.Form["MaxWeight"]; ok {
		modelMaxWeight, err := strconv.Atoi(r.FormValue("MaxWeight"))
		if err != nil {
			return err
		}
		modelData.MaxWeight = modelMaxWeight
	}
	if _, ok := r.Form["Speed"]; ok {
		modelSpeed, err := strconv.Atoi(r.FormValue("Speed"))
		if err != nil {
			return err
		}
		modelData.MaxWeight = modelSpeed
	}

	return nil
}

func getSuppliersPage(w http.ResponseWriter, r *http.Request) {

	EncodeAnswer(FormatHTML, w, nil, HTMLPath+"supplier.html")

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
	ODTService.InsertScootersToDb(filepath)
}


func allSupplierOperation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	r.ParseForm()
	if _, ok := r.Form["ActionType"]; !ok {

		return
	}
	actionType := r.FormValue("ActionType")
	switch actionType {
	case "BlockUser":
		userId, err := strconv.Atoi(r.FormValue("UserID"))
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
		err = userService.ChangeUsersBlockStatus(userId)
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
	default:
		EncodeError(format, w, ErrorRendererDefault(fmt.Errorf("unknown supplier operation")))
	}
	getAllUsers(w, r)
}