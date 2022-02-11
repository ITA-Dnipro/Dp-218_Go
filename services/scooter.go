package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

//ScooterService is the service which gives access to the ScooterRepo repository.
type ScooterService struct {
	RepoScooter repositories.ScooterRepo
}

//NewScooterService creates the new ScooterService.
func NewScooterService(repoScooter repositories.ScooterRepo) *ScooterService {
	return &ScooterService{RepoScooter: repoScooter}
}

//GetAllScooters gives the access to the ScooterRepo.GetAllScooters function.
func (ser *ScooterService) GetAllScooters() (*models.ScooterListDTO, error) {
	return ser.RepoScooter.GetAllScooters()
}
func (ser *ScooterService) GetAllScootersByStationID(stationID int) (*models.ScooterListDTO, error) {
	return ser.RepoScooter.GetAllScootersByStationID(stationID)
}

//GetScooterById gives the access to the ScooterRepo.GetScooterById function.
func (ser *ScooterService) GetScooterById(uid int) (models.ScooterDTO, error) {
	return ser.RepoScooter.GetScooterById(uid)
}

//GetScooterStatus gives the access to the ScooterRepo.GetScooterStatus function.
func (ser *ScooterService) GetScooterStatus(scooterID int) (models.ScooterStatus, error) {
	return ser.RepoScooter.GetScooterStatus(scooterID)
}

//SendCurrentStatus gives the access to the ScooterRepo.SendCurrentStatus function.
func (ser *ScooterService) SendCurrentStatus(id, stationID int, lat, lon, battery float64) error {
	return ser.RepoScooter.SendCurrentStatus(id, stationID, lat, lon, battery)
}

//CreateScooterStatusInRent gives the access to the ScooterRepo.CreateScooterStatusInRent function.
func (ser *ScooterService) CreateScooterStatusInRent(scooterID int) (models.ScooterStatusInRent, error) {
	return ser.RepoScooter.CreateScooterStatusInRent(scooterID)
}
