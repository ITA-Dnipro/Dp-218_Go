package repositories

import (
	"Dp218Go/models"
)

type SupplierRepoI interface {
	GetModels()(*models.ModelODTList, error)
	SelectModel(id int)(*models.ModelODT, error)
	AddModel(modelData *models.ModelODT)error
	EditPrice(modelData *models.ModelODT) error
	GetPrices()(*models.SupplierPricesODTList, error)
	GetPrice(paymentTypeId, userId int) (int, error)
	AddPaymentTypeId(modelName string)(int,error)
	GetPaymentTypeByModelName(modelName string) (int, error)


	GetAllScooters() (*models.ScooterList, error)
	GetScooterByID(id int) (models.Scooter, error)
	AddScooter(scooter *models.Scooter) error
	EditScooter(scooterId int,scooterData models.Scooter) (models.Scooter, error)
	DeleteScooter(id int) error
	ConvertToStruct(path string)[]models.Scooter
}

