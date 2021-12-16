package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

type SupplierService struct {
	SupplierRepo repositories.SupplierRepoI
}

func NewSupplierService(SupplierRepo repositories.SupplierRepoI) *SupplierService {
	return &SupplierService{
		SupplierRepo : SupplierRepo,
	}
}

func (s *SupplierService)GetAllScooters()(*models.ScooterList, error) {
	return s.SupplierRepo.GetAllScooters()
}

func (s *SupplierService)GetScooterByID(id int) (models.Scooter, error){
	return s.SupplierRepo.GetScooterByID(id)
}

func (s *SupplierService)GetScooterById(id int)(models.Scooter, error){
	return s.SupplierRepo.GetScooterByID(id)
}

func (s *SupplierService)AddScooter(scooter *models.Scooter)error{
	return s.SupplierRepo.AddScooter(scooter)
}

func (s *SupplierService)UpdateScooter(scooterId int,scooterData models.Scooter) (models.Scooter, error){
	return s.SupplierRepo.EditScooter(scooterId,scooterData)
}

func (s *SupplierService)DeleteScooter(id int) error{
	return s.SupplierRepo.DeleteScooter(id)
}

func (s *SupplierService)InsertScootersToDb(path string){
	s.SupplierRepo.ConvertToStruct(path)
}

func (s *SupplierService)GetModels()(*models.ModelODTList, error) {
	return s.SupplierRepo.GetModels()
}

func (s *SupplierService)SelectModel(id int)(*models.ModelODT, error) {
	return s.SupplierRepo.SelectModel(id)
}

func (s *SupplierService)AddModel(modelData *models.ModelODT)error {
	return s.SupplierRepo.AddModel(modelData)
}

func (s *SupplierService)ChangePrice(modelData *models.ModelODT)error {
	return s.SupplierRepo.EditPrice(modelData)
}


