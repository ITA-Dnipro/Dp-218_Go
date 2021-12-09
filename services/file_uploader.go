package services

import (
	"Dp218Go/repositories"
	"mime/multipart"
)

type FileService struct {
	repoFile repositories.FileRepoI
}

func NewFileService(repoFile repositories.FileRepoI) *FileService {
	return &FileService{repoFile : repoFile }
}

func (f FileService)InsertScootersToDb(file multipart.File)string{
	tempFilePath := f.repoFile.CreateTempFile(file)
	uploadModel := f.repoFile.ConvertToStruct(tempFilePath)

	f.repoFile.InsertScooterModelData(uploadModel)
	f.repoFile.InsertScooterData(uploadModel)

	return tempFilePath
}