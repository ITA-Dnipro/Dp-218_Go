package models

import "time"

type ScooterModel struct {
	ID          		int         `json:"id"`
	PaymentType 		PaymentType `json:"payment_type"`
	ModelName         	string 		`json:"model_name"`
	MaxWeight         	int    		`json:"max_weight"`
	Speed 			  	int    		`json:"speed"`
}

type ScooterModelList struct {
	ScooterModels []ScooterModel `json:"scooter_models"`
}

type Scooter struct {
	ID      		int 			`json:"id"`
	ScooterModel 	ScooterModel 	`json:"scooter_model"`
	User 			User			`json:"user"`
	SerialNumber 	string 			`json:"serial_number"`
}

type ScooterList struct {
	Scooters []Scooter `json:"scooters"`
}

type Location struct {
	ID      	int 	`json:"id"`
	Latitude    int		`json:"latitude"`
	Longitude   int		`json:"longitude"`
	Label       string  `json:"label"`
}

type LocationList struct {
	Locations []Location `json:"locations"`
}

type ScooterStation struct {
	ID      	int 	 `json:"id"`
	Location    Location `json:"location"`
	Name 		string   `json:"name"`
	IsActive    bool     `json:"is_active"`
}

type ScooterStationList struct {
	ScooterStations []ScooterStation `json:"scooter_stations"`
}

type ScooterStatuses struct {
	ID      	 int 		`json:"id"`
	User       	 User    	`json:"user"`
	Scooter      Scooter 	`json:"scooter"`
	Station      ScooterStation    `json:"station"`
	DateTime 	 time.Time	`json:"date_time"`
	Location     Location	`json:"location"`
	IsReturned   bool		`json:"is_returned"`
}

type ScooterStatusesInRent struct {
	ID      	 int 		`json:"id"`
	User       	 User    	`json:"user"`
	Scooter      Scooter 	`json:"scooter"`
	Station      ScooterStation    `json:"station"`
	DateTime 	 time.Time	`json:"date_time"`
	Location     Location	`json:"location"`
	IsReturned   bool		`json:"is_returned"`
}

type ScooterStatusesInRentList struct {
	ScootersStatusesInRent []ScooterStatusesInRent `json:"scooters_statuses_in_rent"`
}

type SupplierPrices struct {
	ID  	int `json:"id"`
	Price 	int `json:"price"`
	PaymentType PaymentType `json:"payment_type"`
	User    User  `json:"user"`
}

type SupplierPricesList struct {
	SupplierPrices []SupplierPrices `json:"supplier_prices_list"`
}
