package postgres

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"context"
	"fmt"
)

var userId = 1

type SupplierRepoDB struct {
	db repositories.AnyDatabase
}

func NewSupplierRepoDB(db repositories.AnyDatabase) *SupplierRepoDB {
	return &SupplierRepoDB{db}
}

func (s *SupplierRepoDB)GetModels()(*models.ModelODTList, error){
	modelsOdtList := &models.ModelODTList{}
	pricesList := &models.SupplierPricesODTList{}

	pricesList, err := s.GetPrices()
	if err != nil {
		return modelsOdtList, err
	}

	querySQL := `SELECT * FROM scooter_models ORDER BY id DESC;`
	rows, err := s.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return modelsOdtList, err
	}

	for rows.Next() {
		var paymentTypeID int
		var model models.ModelODT
		err := rows.Scan(&model.ID, &paymentTypeID, &model.ModelName, &model.MaxWeight, &model.Speed)
		if err != nil {
			return modelsOdtList, err
		}

		model.Price, err = s.findSupplierPricesList(pricesList, paymentTypeID, userId)
		if err != nil {
			return modelsOdtList, err
		}
		modelsOdtList.ModelsODT = append(modelsOdtList.ModelsODT, model)
	}
	return modelsOdtList, nil
}

func (s *SupplierRepoDB)SelectModel(id int)(*models.ModelODT, error){
	modelODT := &models.ModelODT{}

	querySQL := `SELECT id, payment_type_id, model_name, max_weight, speed  FROM scooter_models WHERE id = $1;`
	row := s.db.QueryResultRow(context.Background(), querySQL, id)

	var paymentTypeId int
	err := row.Scan(&modelODT.ID, &paymentTypeId, &modelODT.ModelName, &modelODT.MaxWeight, &modelODT.Speed)
	if err != nil {
		return modelODT, err
	}

	modelODT.Price, err = s.GetPrice(paymentTypeId, userId)

	return modelODT, err
}

func (s *SupplierRepoDB)AddModel(modelData *models.ModelODT)error{

	paymentTypeId, err := s.AddPaymentTypeId(modelData.ModelName)
	if err != nil {
		return err
	}
	var modelId int
	querySQL := `INSERT INTO scooter_models(payment_type_id, model_name, max_weight, speed)
	   		VALUES($1, $2, $3, $4)
	   		RETURNING id;`
	err = s.db.QueryResultRow(context.Background(), querySQL, &paymentTypeId, modelData.ModelName, modelData.MaxWeight, modelData.Speed).Scan(&modelId)
	if err != nil {
		return  err
	}

	var priceId int
	querySQL = `INSERT INTO supplier_prices(price, payment_type_id, user_id)
	   		VALUES($1, $2, $3)
	   		RETURNING id;`
	err = s.db.QueryResultRow(context.Background(), querySQL, modelData.Price, paymentTypeId, userId).Scan(&priceId)
	if err != nil {
		return  err
	}
	return  nil
}

func (s *SupplierRepoDB) EditPrice(modelData *models.ModelODT) error{
	price := &models.ModelODT{}
	paymentTypeId, err := s.GetPaymentTypeByModelName(modelData.ModelName)
	if err != nil {
		return err
	}

	querySQL := `UPDATE supplier_prices SET price=$1 WHERE payment_type_id = $2 AND user_id = $3 RETURNING price;`
	err = s.db.QueryResultRow(context.Background(), querySQL, modelData.Price, paymentTypeId, userId).Scan(&price.Price)
	if err != nil {
		return err
	}

	return nil
}

func (s *SupplierRepoDB)GetPrices()(*models.SupplierPricesODTList, error){
	list := &models.SupplierPricesODTList{}

	querySQL := `SELECT * FROM supplier_prices ORDER BY id DESC;`
	rows, err := s.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var supplierPriceODT models.SupplierPricesODT
		err := rows.Scan(&supplierPriceODT.ID, &supplierPriceODT.Price, &supplierPriceODT.PaymentTypeID, &supplierPriceODT.UserId)

		if err != nil {
			return list, err
		}

		list.SupplierPricesODT = append(list.SupplierPricesODT, supplierPriceODT)
	}
	return list, nil
}

func (s *SupplierRepoDB)findSupplierPricesList(supplierPrice *models.SupplierPricesODTList, paymentTypeId int, userId int )(int, error){
	for _, v := range supplierPrice.SupplierPricesODT {
		if v.PaymentTypeID == paymentTypeId && v.UserId == userId{
			return v.Price, nil
		}
	}
	return 0, fmt.Errorf("not found paymentType id=%d", paymentTypeId)
}

func (s *SupplierRepoDB) GetPrice(paymentTypeId, userId int) (int, error) {
	price := models.ModelODT{}
	querySQL := `SELECT price FROM supplier_prices WHERE payment_type_id = $1 AND user_id = $2;`
	row := s.db.QueryResultRow(context.Background(), querySQL, paymentTypeId, userId)
	err := row.Scan(&price.Price)

	return price.Price, err
}

func (s *SupplierRepoDB) AddPaymentTypeId(modelName string)(int,error){
	var paymentTypeId int
	querySQL := `INSERT INTO payment_types (name) VALUES ($1) RETURNING id;`
	err := s.db.QueryResultRow(context.Background(), querySQL, modelName).Scan(&paymentTypeId)
	if err != nil {
		return 0,err
	}
	return  paymentTypeId, nil
}

func (s *SupplierRepoDB)GetPaymentTypeByModelName(modelName string) (int, error) {
	paymentType := models.PaymentType{}
	querySQL := `SELECT * FROM payment_types WHERE name = $1;`
	row := s.db.QueryResultRow(context.Background(), querySQL, modelName)
	err := row.Scan(&paymentType.ID, &paymentType.Name)
	return paymentType.ID, err
}
