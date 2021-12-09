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

func (sup *SupplierRepoDB)CreateScooterModel(model *models.ScooterModel)error{
	var id int
	querySQL := `INSERT INTO scooter_model(payment_type_id, model_name, max_weight, speed) 
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

	querySQL := `SELECT * FROM scooter_model ORDER BY id DESC;`
	rows, err := sup.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var model models.ScooterModel
		var paymentTypeId int
		err := rows.Scan(&model.ID, &model.ModelName, &model.MaxWeight, &model.Speed, &paymentTypeId)
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

	querySQL := `SELECT * FROM scooter_models WHERE id = $1;`
	row := sup.db.QueryResultRow(context.Background(), querySQL, modelId)
	err := row.Scan(&scooterModel.ID, &scooterModel.PaymentType.ID, &scooterModel.ModelName, &scooterModel.MaxWeight, &scooterModel.Speed)
	return scooterModel, err
}

func (sup *SupplierRepoDB)GetAllScooters() (*models.ScooterList, error) {
	list := &models.ScooterList{}

	scooterModels, err := sup.GetScooterModels()
	if err != nil {
		return list, err
	}

	querySQL := `SELECT * FROM scooter ORDER BY id DESC;`
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

	querySQL := `SELECT * FROM scooters WHERE id = $1;`
	row := sup.db.QueryResultRow(context.Background(), querySQL, id)

	err := row.Scan(&scooter.ID, &scooter.ScooterModel.ID, &userId,  &scooter.SerialNumber)
	return scooter, err
}

func (sup *SupplierRepoDB)GetScootersByModelId(id int)(*models.ScooterList, error){
	list := &models.ScooterList{}

	querySQL := `SELECT * FROM scooters WHERE scooter_model = $1;`
	rows, err := sup.db.QueryResult(context.Background(), querySQL, id)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var scooter models.Scooter
		err := rows.Scan(&scooter.ID, &scooter.ScooterModel.ID, &scooter.User,
			&scooter.SerialNumber)
		if err != nil {
			return list, err
		}

		list.Scooters = append(list.Scooters, scooter)
	}
	return list, nil
}


func (sup *SupplierRepoDB) AddScooter(scooter *models.Scooter) error {
	var id int
	querySQL := `INSERT INTO scooters(scooter_model, user, serial_number)
	   		VALUES($1, $2, $3,)
	   		RETURNING id;`
	err := sup.db.QueryResultRow(context.Background(), querySQL, scooter.ScooterModel.ID, scooter.User.ID, scooter.SerialNumber).Scan(&id)
	if err != nil {
		return err
	}
	scooter.ID = id
	return nil
}

func (sup *SupplierRepoDB) UpdateScooter(scooterId int, scooterData models.Scooter)(models.Scooter, error){
	scooter := models.Scooter{}
	querySQL := `UPDATE scooters
	   		SET serial_number=$1,
	   		WHERE id=$2
	   		RETURNING id, serial_number;`
	err := sup.db.QueryResultRow(context.Background(), querySQL, scooterData.SerialNumber).Scan(&scooter.ID,&scooter.SerialNumber)
	if err != nil {
		return scooter, err
	}

	return scooter, nil
}


func (sup *SupplierRepoDB) DeleteScooter(id int) error {
	querySQL := `DELETE FROM scooter WHERE id = $1;`
	_, err := sup.db.QueryExec(context.Background(), querySQL, userId)
	return err
}

func (sup *SupplierRepoDB) GetPaymentTypes() (*models.PaymentTypeList, error) {
	list := &models.PaymentTypeList{}

	querySQL := `SELECT * FROM payment_type ORDER BY id DESC;`
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
	return models.PaymentType{}, fmt.Errorf("not found role id=%d", paymentTypeId)
}

func (sup *SupplierRepoDB) GetPaymentTypeById(paymentTypeId int) (models.PaymentType, error) {
	paymentType := models.PaymentType{}
	querySQL := `SELECT * FROM payment_type WHERE id = $1;`
	row := sup.db.QueryResultRow(context.Background(), querySQL, paymentTypeId)
	err := row.Scan(&paymentType.ID, &paymentType.Name)
	return paymentType, err
}