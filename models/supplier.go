package models

import "time"

type ScooterUploaded struct {
	Model 		 		ScooterModel    `json,scv:"scooter_model"`
	Owner			 	User 			`json,scv:"user"`
	Scooter 	  		Scooter  		`json,scv:"scooter"`
}

type ScooterUploadedList struct {
	ScooterUploaded []ScooterUploaded `json:"roles"`
}

type ScooterModel struct {
	ID          int         `json:"id"`
	PaymentType PaymentType `json:"payment_type"`
	ModelName         	string 		`json:"model_name"`
	MaxWeight         	int    		`json:"max_weight"`
	Speed 			  	int    		`json:"speed"`
}

type ScooterModelList struct {
	ScooterModels []ScooterModel `json:"roles"`
}

type Scooter struct {
	ID      		int 			`json:"id"`
	ScooterModel 	ScooterModel 	`json:"scooter_model"`
	User 			User			`json:"user"`
	SerialNumber 	string 			`json:"serial_number"`
}

type ScooterList struct {
	Scooters []Scooter `json:"roles"`
}

type Location struct {
	ID      	int 	`json:"id"`
	Latitude    int		`json:"latitude"`
	Longitude   int		`json:"longitude"`
	Label       string  `json:"label"`
}

type LocationList struct {
	Locations []Location `json:"roles"`
}

type ScooterStation struct {
	ID      	int 	 `json:"id"`
	Location    Location `json:"location"`
	Name 		string   `json:"name"`
	IsActive    bool     `json:"is_active"`
}

type ScooterStationList struct {
	ScooterStations []ScooterStation `json:"roles"`
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
	ScootersStatusesInRent []ScooterStatusesInRent `json:"roles"`
}
