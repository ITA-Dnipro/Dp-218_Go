package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var scooterService *services.ScooterService
var scooterGrpcService *services.GrpcScooterService
var orderService *services.OrderService
var scooterIDKey = "scooterId"

var chosenScooterID, chosenStationID int
var userFromRequest = models.User{ID: 1, LoginEmail: "guru_admin@guru.com", UserName: "Guru", UserSurname: "Sadh"}

var scooterRoutes = []Route{
	{
		Uri:     `/scooters`,
		Method:  http.MethodGet,
		Handler: getAllScooters,
	},
	{
		Uri:     `/scooter/{` + scooterIDKey + `}`,
		Method:  http.MethodGet,
		Handler: getScooterById,
	},
	{
		Uri:     `/start-trip`,
		Method:  http.MethodGet,
		Handler: showTripPage,
	},
	{
		Uri:     `/run`,
		Method:  http.MethodGet,
		Handler: startScooterTrip,
	},
	{
		Uri:     `/choose-scooter`,
		Method:  http.MethodPost,
		Handler: chooseScooter,
	},
	{
		Uri:     `/choose-station`,
		Method:  http.MethodPost,
		Handler: chooseStation,
	},
}

type combineForTemplate struct {
	models.ScooterListDTO
	models.StationList
}

//AddScooterHandler adds routes to the router from the list of routes.
func AddScooterHandler(router *mux.Router, service *services.ScooterService) {
	scooterService = service
	for _, rt := range scooterRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

//AddGrpcScooterHandler adds routes to the router from the list of routes.
func AddGrpcScooterHandler(router *mux.Router, service *services.GrpcScooterService) {
	scooterGrpcService = service
	for _, rt := range scooterRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getAllScooters(w http.ResponseWriter, r *http.Request) {

	scooters, err := scooterService.GetAllScooters()

	if err != nil {
		ServerErrorRender(FormatJSON, w)
		fmt.Println(err)
		return
	}

	EncodeAnswer(FormatJSON, w, scooters)
}

func getScooterById(w http.ResponseWriter, r *http.Request) {

	scooterID, err := strconv.Atoi(mux.Vars(r)[scooterIDKey])
	if err != nil {
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
		return
	}

	scooter, err := scooterService.GetScooterById(scooterID)
	if err != nil {
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, scooter)
}

func startScooterTrip(w http.ResponseWriter, r *http.Request) {
	station, err := stationService.GetStationById(chosenStationID)
	if err != nil {
		fmt.Println(err)
	}

	chosenWay := models.Coordinate{Latitude: station.Latitude, Longitude: station.Longitude}
	statusStart, err := scooterService.CreateScooterStatusInRent(chosenScooterID)
	if err != nil {
		fmt.Println(err)
	}


	err = scooterGrpcService.InitAndRun(chosenScooterID, chosenWay)
	if err != nil {
		fmt.Println(err)
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
	}

	statusEnd, err := scooterService.CreateScooterStatusInRent(chosenScooterID)

	distance := statusEnd.Location.Distance(statusStart.Location)

	_, err = orderService.CreateOrder(userFromRequest, chosenScooterID, statusStart.ID, statusEnd.ID, distance)
	if err != nil {
		fmt.Println(err)
	}
}

func showTripPage(w http.ResponseWriter, r *http.Request) {
	scooterList, err := scooterService.GetAllScooters()
	if err != nil {
		fmt.Println(err)
	}

	stationList, err := stationService.GetAllStations()
	if err != nil {
		fmt.Println(err)
	}

	EncodeAnswer(FormatHTML, w, &combineForTemplate{*scooterList, *stationList}, HTMLPath+"scooter-run.html")
}

func chooseScooter(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	chosenScooterID, err = strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		fmt.Println(err)
	}
}

func chooseStation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	chosenStationID, err = strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		fmt.Println(err)
	}
}