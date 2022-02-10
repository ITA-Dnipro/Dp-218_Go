package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

//OrderService is the service which gives access to the OrderRepo repository.
type OrderService struct {
	RepoOrder repositories.OrderRepo
}

//NewOrderService creates the new OrderService.
func NewOrderService(orderRepo repositories.OrderRepo) *OrderService {
	return &OrderService{RepoOrder: orderRepo}
}

//CreateOrder gives the access to the OrderRepo.CreateOrder function.
func (ors *OrderService) CreateOrder(user models.User, scooterID, startID, endID int,
	distance float64) (models.Order, error) {
	return ors.RepoOrder.CreateOrder(user, scooterID, startID, endID, distance)
}

//GetAllOrders gives the access to the OrderRepo.GetAllOrders function.
func (ors *OrderService) GetAllOrders() (*models.OrderList, error) {
	return ors.RepoOrder.GetAllOrders()
}

//GetOrderByID gives the access to the OrderRepo.GetOrderByID function.
func (ors *OrderService) GetOrderByID(orderID int) (models.Order, error) {
	return ors.RepoOrder.GetOrderByID(orderID)
}

//GetOrdersByUserID gives the access to the OrderRepo.GetOrdersByUserID function.
func (ors *OrderService) GetOrdersByUserID(userID int) (models.OrderList, error) {
	return ors.RepoOrder.GetOrdersByUserID(userID)
}

//GetOrdersByScooterID gives the access to the OrderRepo.GetOrdersByScooterID function.
func (ors *OrderService) GetOrdersByScooterID(scooterID int) (models.OrderList, error) {
	return ors.RepoOrder.GetOrdersByScooterID(scooterID)
}

//GetScooterMileageByID gives the access to the OrderRepo.GetScooterMileageByID function.
func (ors *OrderService) GetScooterMileageByID(scooterID int) (float64, error) {
	return ors.RepoOrder.GetScooterMileageByID(scooterID)
}

//GetUserMileageByID gives the access to the OrderRepo.GetUserMileageByID function.
func (ors *OrderService) GetUserMileageByID(userID int) (float64, error) {
	return ors.RepoOrder.GetUserMileageByID(userID)
}

//UpdateOrder gives the access to the OrderRepo.UpdateOrder function.
func (ors *OrderService) UpdateOrder(orderID int, orderData models.Order) (models.Order, error) {
	return ors.RepoOrder.UpdateOrder(orderID, orderData)
}

//DeleteOrder ives the access to the OrderRepo.DeleteOrder function.
func (ors *OrderService) DeleteOrder(orderID int) error {
	return ors.RepoOrder.DeleteOrder(orderID)
}
