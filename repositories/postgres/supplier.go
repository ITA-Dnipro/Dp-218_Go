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

func NewSupplierRepoDB(db repositories.AnyDatabase) *SupplierRepoDB{
	return &SupplierRepoDB{db}
}

func (sup *SupplierRepoDB)CreateScooterModel(model *models.ScooterModel)error{
	var id int
	querySQL := `INSERT INTO scooter_models(payment_type_id, model_name, max_weight, speed) 
		VALUES($1, $2, $3, $4)
		RETURNING id;`
	err := sup.db.QueryResultRow(context.Background(), querySQL, model.PaymentType.ID, model.ModelName, model.MaxWeight, model.Speed).Scan(&id)
	if err != nil {
		return err
	}
	model.ID = id
	return nil
}

func (sup *SupplierRepoDB)GetScooterModels() (*models.ScooterModelList, error) {
	list := &models.ScooterModelList{}

	pTypes, err := sup.GetPaymentTypes()
	if err != nil {
		return list, err
	}

	querySQL := `SELECT id, payment_type_id, model_name, max_weight, speed FROM scooter_models ORDER BY id DESC;`
	rows, err := sup.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var model models.ScooterModel
		var paymentTypeId int
		err := rows.Scan(&model.ID, &paymentTypeId, &model.ModelName, &model.MaxWeight, &model.Speed)
		if err != nil {
			return list, err
		}

		model.PaymentType, err = sup.FindPaymentTypeList(pTypes, paymentTypeId)
		if err != nil {
			return list, err
		}

		list.ScooterModels = append(list.ScooterModels, model)
	}
	return list, nil
}

func (sup *SupplierRepoDB) GetScooterModelById(modelId int) (models.ScooterModel, error) {
	scooterModel := models.ScooterModel{}

	querySQL := `SELECT id, payment_type_id, model_name, max_weight, speed  FROM scooter_models WHERE id = $1;`
	row := sup.db.QueryResultRow(context.Background(), querySQL, modelId)

	var paymentTypeId int
	err := row.Scan(&scooterModel.ID, &paymentTypeId, &scooterModel.ModelName, &scooterModel.MaxWeight, &scooterModel.Speed)
	if err != nil {
		return models.ScooterModel{}, err
	}
	scooterModel.PaymentType, err = sup.GetPaymentTypeById(paymentTypeId)

	return scooterModel, err
}

func (sup *SupplierRepoDB)GetAllScooters() (*models.ScooterList, error) {
	list := &models.ScooterList{}

	scooterModels, err := sup.GetScooterModels()
	if err != nil {
		return list, err
	}

	querySQL := `SELECT * FROM scooters ORDER BY id DESC;`
	rows, err := sup.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var scooter models.Scooter
		var scooterModelId int
		err := rows.Scan(&scooter.ID, &scooterModelId, &userId, &scooter.SerialNumber)
		if err != nil {
			return list, err
		}

		scooter.ScooterModel, err = sup.FindScooterList(scooterModels, scooterModelId)
		if err != nil {
			return list, err
		}

		list.Scooters = append(list.Scooters, scooter)
	}
	return list, nil
}

func (sup *SupplierRepoDB)FindScooterList(scooterModel *models.ScooterModelList, scooterModelId int ) (models.ScooterModel, error) {
	for _, v := range scooterModel.ScooterModels {
		if v.ID == scooterModelId{
			return v, nil
		}
	}
	return models.ScooterModel{}, fmt.Errorf("not found scooterModel id=%d", scooterModelId)
}

func (sup *SupplierRepoDB) GetScooterByID(id int)(models.Scooter, error){
	scooter := models.Scooter{}

	querySQL := `SELECT id, model_id, owner_id, serial_number FROM scooters WHERE id = $1;`
	row := sup.db.QueryResultRow(context.Background(), querySQL, id)

	var modelId int
	err := row.Scan(&scooter.ID, &modelId,  &userId, &scooter.SerialNumber)
	if err != nil {
		return models.Scooter{}, err
	}
	scooter.ScooterModel, err = sup.GetScooterModelById(modelId)

	return scooter, err
}

func (sup *SupplierRepoDB) AddScooter(scooter *models.Scooter) error {
	var id int
	querySQL := `INSERT INTO scooters(model_id, owner_id, serial_number)
	   		VALUES($1, $2, $3)
	   		RETURNING id;`
	err := sup.db.QueryResultRow(context.Background(), querySQL, scooter.ScooterModel.ID, scooter.User.ID, scooter.SerialNumber).Scan(&id)
	if err != nil {
		return err
	}
//	scooter.ID = id
	return nil
}

func (sup *SupplierRepoDB) EditScooter(scooterId int, scooterData models.Scooter)(models.Scooter, error){
	scooter := models.Scooter{}
	querySQL := `UPDATE scooters
	   		SET model_id=$1, owner_id=$2, serial_number=$3
	   		WHERE id=$4
	   		RETURNING id, model_id, owner_id, serial_number;`
	var modelId int
	err := sup.db.QueryResultRow(context.Background(), querySQL, scooterData.ScooterModel.ID, scooterData.User.ID, scooterData.SerialNumber, scooterId).Scan(
		&scooter.ID,  &modelId, &userId,  &scooter.SerialNumber)
	if err != nil {
		return scooter, err
	}

	scooter.ScooterModel, err = sup.GetScooterModelById(modelId)
	if err != nil {
		return scooter, err
	}
	return scooter, nil
}

func (sup *SupplierRepoDB) DeleteScooter(id int) error {
	querySQL := `DELETE FROM scooters WHERE id = $1;`
	_, err := sup.db.QueryExec(context.Background(), querySQL, userId)
	return err
}

func (sup *SupplierRepoDB) GetPaymentTypes() (*models.PaymentTypeList, error) {
	list := &models.PaymentTypeList{}

	querySQL := `SELECT * FROM payment_types ORDER BY id DESC;`
	rows, err := sup.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var paymentType models.PaymentType
		err := rows.Scan(&paymentType.ID, &paymentType.Name)
		if err != nil {
			return list, err
		}
		list.PaymentTypes = append(list.PaymentTypes, paymentType)
	}
	return list, nil
}

func (sup *SupplierRepoDB)FindPaymentTypeList(paymentType *models.PaymentTypeList, paymentTypeId int) (models.PaymentType, error) {
	for _, v := range paymentType.PaymentTypes {
		if v.ID == paymentTypeId{
			return v, nil
		}
	}
	return models.PaymentType{}, fmt.Errorf("not found paymentType id=%d", paymentTypeId)
}

func (sup *SupplierRepoDB) GetPaymentTypeById(paymentTypeId int) (models.PaymentType, error) {
	paymentType := models.PaymentType{}
	querySQL := `SELECT * FROM payment_types WHERE id = $1;`
	row := sup.db.QueryResultRow(context.Background(), querySQL, paymentTypeId)
	err := row.Scan(&paymentType.ID, &paymentType.Name)
	return paymentType, err
}