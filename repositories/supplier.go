package repositories

import "Dp218Go/models"

type ScooterModelRepoI interface {
	CreateScooterModel(scooterModel *models.ScooterModel)error
	GetScooterModels() (*models.ScooterModelList, error)
	GetScooterModelById(modelId int)(models.ScooterModel, error)
}

type ScooterRepoI interface {
	GetAllScooters() (*models.ScooterList, error)
	GetScooterByID(id int) (models.Scooter, error)
	GetScootersByModelId(id int)(*models.ScooterList, error)
	AddScooter(scooter *models.Scooter) error
	UpdateScooter(scooterId int,scooterData models.Scooter) (models.Scooter, error)
	DeleteScooter(id int) error
	FindScooterList(scooterModel *models.ScooterModelList, scooterModelId int ) (models.ScooterModel, error)
}


type PaymentTypeI interface {
	GetPaymentTypes() (*models.PaymentTypeList, error)
	GetPaymentTypeById(paymentTypeId int) (models.PaymentType, error)
	FindPaymentTypeList(paymentType *models.PaymentTypeList, paymentTypeId int)(models.PaymentType, error)
}

