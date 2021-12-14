package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

type ODTService struct {
	ODTRepo repositories.ODTRepoI
}

func NewODTService(ODTRepo repositories.ODTRepoI) *ODTService {
	return &ODTService{
		ODTRepo : ODTRepo,
	}
}

func (odt *ODTService)GetModels()(*models.ModelODTList, error) {
	return odt.ODTRepo.GetModels()
}

func (odt *ODTService)SelectModel(id int)(*models.ModelODT, error) {
	return odt.ODTRepo.SelectModel(id)
}

func (odt *ODTService)AddModel(modelData *models.ModelODT)error {
	return odt.ODTRepo.AddModel(modelData)
}

func (odt *ODTService)EditModel(modelData models.ModelODT)error {
	return odt.ODTRepo.EditModel(modelData)
}

func (odt *ODTService)InsertScootersToDb(path string){
	odt.ODTRepo.ConvertToStruct(path)
}

