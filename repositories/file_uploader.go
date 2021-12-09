package repositories

import (
	"Dp218Go/models"
	"mime/multipart"
)

type FileRepoI interface {
	CreateTempFile(file multipart.File)string
	ConvertToStruct(path string)[]models.ScooterUploaded
	InsertScooterModelData(scooters []models.ScooterUploaded)error
	InsertScooterData(scooters []models.ScooterUploaded)error
}