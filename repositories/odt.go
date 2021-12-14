package repositories

import (
	"Dp218Go/models"
)

type ODTRepoI interface {
	GetModels()(*models.ModelODTList, error)
	SelectModel(id int)(*models.ModelODT, error)
	AddModel(modelData *models.ModelODT)error
	EditModel(modelData models.ModelODT) error
	GetPrices()(*models.SupplierPricesODTList, error)
	GetPrice(paymentTypeId, userId int) (int, error)
	ConvertToStruct(path string)[]models.Scooter
}