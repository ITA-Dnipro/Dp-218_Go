package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var scooterInitService *services.ScooterInitService
var scooterInitRoutes = []Route{
	{
		Uri:     `/init`,
		Method:  http.MethodGet,
		Handler: getAllocationData,
	},
	{
		Uri:     `/transfer`,
		Method:  http.MethodPost,
		Handler: addStatusesToScooters,
	},
}

func AddScooterInitHandler(router *mux.Router, service *services.ScooterInitService) {
	scooterInitService = service
	for _, rt := range scooterInitRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getAllocationData(w http.ResponseWriter, r *http.Request){
	var dataAllocation = &models.ScootersStationsAllocation{}
	var err error
	format := GetFormatFromRequest(r)
	err = r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	dataAllocation = scooterInitService.ConvertForTemplateStruct()

	EncodeAnswer(format, w, dataAllocation, HTMLPath+"scooters-init.html")
}

func addStatusesToScooters(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	scooterIds := r.Form["new_data"]
	stationId := r.Form["station_data"]
	intStationId, err := strconv.Atoi(stationId[0])
	if err != nil {
		log.Println(err)
	}

	var intScooterIds []int
	for _, i := range scooterIds {
		intId, err := strconv.Atoi(i)
		if err != nil {
			log.Println(err)
		}
		intScooterIds = append(intScooterIds, intId)
	}

	stationData, err := stationService.GetStationById(intStationId)

	err = scooterInitService.AddStatusesToScooters(intScooterIds, stationData)
	if err != nil {
		return
	}

	http.Redirect(w,r,"/init",  http.StatusFound)
}

