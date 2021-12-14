package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

type SupplierService struct {
	ScooterModelRepo repositories.ScooterModelRepoI
	PaymentTypeRepo repositories.PaymentTypeRepoI
	ScooterRepo repositories.ScooterRepoI
}

func NewSupplierService(ScooterModelRepo repositories.ScooterModelRepoI, PaymentTypeRepo repositories.PaymentTypeRepoI, ScooterRepo repositories.ScooterRepoI) *SupplierService {
	return &SupplierService{
		ScooterModelRepo : ScooterModelRepo,
		PaymentTypeRepo : PaymentTypeRepo,
		ScooterRepo : ScooterRepo,
	}
}

func (sup *SupplierService)CreateScooterModel(scooterModel *models.ScooterModel)error{
	return sup.ScooterModelRepo.CreateScooterModel(scooterModel)
}

func (sup *SupplierService)GetScooterModels()(*models.ScooterModelList, error) {
	return sup.ScooterModelRepo.GetScooterModels()
}

func (sup *SupplierService)GetScooterModelById(modelId int) (models.ScooterModel, error){
	return sup.ScooterModelRepo.GetScooterModelById(modelId)
}

func (sup *SupplierService)GetAllScooters()(*models.ScooterList, error) {
	return sup.ScooterRepo.GetAllScooters()
}

func (sup *SupplierService)FindScooterList(scooterModel *models.ScooterModelList, scooterModelId int)(models.ScooterModel, error){
	return sup.ScooterRepo.FindScooterList(scooterModel,scooterModelId)
}

func (sup *SupplierService)GetScooterByID(id int) (models.Scooter, error){
	return sup.ScooterRepo.GetScooterByID(id)
}

func (sup *SupplierService)GetScooterById(id int)(models.Scooter, error){
	return sup.ScooterRepo.GetScooterByID(id)
}

func (sup *SupplierService)AddScooter(scooter *models.Scooter)error{
	return sup.ScooterRepo.AddScooter(scooter)
}

func (sup *SupplierService)UpdateScooter(scooterId int,scooterData models.Scooter) (models.Scooter, error){
	return sup.ScooterRepo.EditScooter(scooterId,scooterData)
}

func (sup *SupplierService)DeleteScooter(id int) error{
	return sup.ScooterRepo.DeleteScooter(id)
}

func (sup *SupplierService)GetPaymentTypes() (*models.PaymentTypeList, error){
	return sup.PaymentTypeRepo.GetPaymentTypes()
}

func (sup *SupplierService)GetPaymentTypeById(paymentTypeId int) (models.PaymentType, error) {
	return  sup.PaymentTypeRepo.GetPaymentTypeById(paymentTypeId)
}

func (sup *SupplierService)FindPaymentTypeList(paymentType *models.PaymentTypeList, paymentTypeId int)(models.PaymentType, error){
	return  sup.PaymentTypeRepo.FindPaymentTypeList(paymentType,paymentTypeId)
}

