package postgres

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jszwec/csvutil"
	"io"
	"os"
)


type ODTRepoDB struct {
	db repositories.AnyDatabase
}

func NewODTRepoDB(db repositories.AnyDatabase) *ODTRepoDB{
	return &ODTRepoDB{db}
}

func (o *ODTRepoDB)GetModels()(*models.ModelODTList, error){
	modelsOdtList := &models.ModelODTList{}
	pricesList := &models.SupplierPricesODTList{}

	pricesList, err := o.GetPrices()
	if err != nil {
		return modelsOdtList, err
	}

	querySQL := `SELECT * FROM scooter_models ORDER BY id DESC;`
	rows, err := o.db.QueryResult(context.Background(), querySQL)
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

		model.Price, err = o.findSupplierPricesList(pricesList, paymentTypeID, userId)
		if err != nil {
			return modelsOdtList, err
		}
		modelsOdtList.ModelsODT = append(modelsOdtList.ModelsODT, model)
	}
	return modelsOdtList, nil
}

func (o *ODTRepoDB)SelectModel(id int)(*models.ModelODT, error){
	modelODT := &models.ModelODT{}

	querySQL := `SELECT id, payment_type_id, model_name, max_weight, speed  FROM scooter_models WHERE id = $1;`
	row := o.db.QueryResultRow(context.Background(), querySQL, id)

	var paymentTypeId int
	err := row.Scan(&modelODT.ID, &paymentTypeId, &modelODT.ModelName, &modelODT.MaxWeight, &modelODT.Speed)
	if err != nil {
		return modelODT, err
	}

	modelODT.Price, err = o.GetPrice(paymentTypeId, userId)

	return modelODT, err
}

func (o *ODTRepoDB)AddModel(modelData *models.ModelODT)error{

	var paymentTypeId int
	querySQL := `INSERT INTO payment_types (name) VALUES($1) RETURNING id;`
	err := o.db.QueryResultRow(context.Background(), querySQL, modelData.ModelName).Scan(&paymentTypeId)
	if err != nil {
		return err
	}

	var modelId int
	querySQL = `INSERT INTO scooter_models(payment_type_id, model_name, max_weight, speed)
	   		VALUES($1, $2, $3, $4)
	   		RETURNING id;`
	err = o.db.QueryResultRow(context.Background(), querySQL, paymentTypeId, modelData.ModelName, modelData.MaxWeight, modelData.Speed).Scan(&modelId)
	if err != nil {
		return  err
	}

	var priceId int
	querySQL = `INSERT INTO supplier_prices(price, payment_type_id, user_id)
	   		VALUES($1, $2, $3, $4)
	   		RETURNING id;`
	err = o.db.QueryResultRow(context.Background(), querySQL, modelData.Price, paymentTypeId, userId).Scan(&priceId)
	if err != nil {
		return  err
	}

	return  nil
}

func (o *ODTRepoDB)EditModel(modelData models.ModelODT) error{
	model:= &models.ModelODT{}

	var paymentTypeId int
	querySQL := `SELECT payment_type_id FROM scooter_models WHERE id = $1;`
	row := o.db.QueryResultRow(context.Background(), querySQL, modelData.ID)
	err := row.Scan(&paymentTypeId)
	if err != nil {
		return err
	}

	_, err = o.SetPrice(modelData, paymentTypeId, userId)
	if err != nil {
		return err
	}

	querySQL = `UPDATE scooter_models 
	   		SET  model_name=$2, max_weight=$3, speed=$4
	   		WHERE id=$5
	   		RETURNING id, model_name, max_weight, speed;`
	err = o.db.QueryResultRow(context.Background(), querySQL, modelData.ModelName, modelData.MaxWeight, modelData.Speed, modelData.ID).Scan(
		&model.ID, &model.ModelName, &model.MaxWeight, &model.Speed)
	if err != nil {
		return err
	}

	return nil
}

func (o *ODTRepoDB)GetPrices()(*models.SupplierPricesODTList, error){
	list := &models.SupplierPricesODTList{}

	querySQL := `SELECT * FROM supplier_prices ORDER BY id DESC;`
	rows, err := o.db.QueryResult(context.Background(), querySQL)
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

func (o *ODTRepoDB)SetPrice(modelData models.ModelODT,paymentTypeId, userId int) (*models.ModelODT, error){
	price := &models.ModelODT{}
	querySQL := `UPDATE supplier_prices SET price=$1 WHERE payment_type_id = $2 AND user_id = $3 RETURNING price;`
	err := o.db.QueryResultRow(context.Background(), querySQL, modelData.Price, paymentTypeId, userId).Scan(&price.Price)

	return price, err
}

func (o *ODTRepoDB) GetPrice(paymentTypeId, userId int) (int, error) {
	price := models.ModelODT{}
	querySQL := `SELECT price FROM supplier_prices WHERE payment_type_id = $1 AND user_id = $2;`
	row := o.db.QueryResultRow(context.Background(), querySQL, paymentTypeId, userId)
	err := row.Scan(&price.Price)

	return price.Price, err
}

func (o *ODTRepoDB)findSupplierPricesList(supplierPrice *models.SupplierPricesODTList, paymentTypeId int, userId int )(int, error){
	for _, v := range supplierPrice.SupplierPricesODT {
		if v.PaymentTypeID == paymentTypeId && v.UserId == userId{
			return v.Price, nil
		}
	}
	return 0, fmt.Errorf("not found paymentType id=%d", paymentTypeId)
}

func (o *ODTRepoDB) ConvertToStruct(path string) []models.Scooter{

	csvFile, _ := os.Open(path)
	reader := csv.NewReader(csvFile)
	reader.Comma = ';'

	scooterHeader, _ := csvutil.Header(models.Scooter{}, "csv")
	dec, _ := csvutil.NewDecoder(reader, scooterHeader...)

	var fileData []models.Scooter
	for {
		var s models.Scooter
		if err := dec.Decode(&s.SerialNumber); err == io.EOF {
			break
		}
		fileData = append(fileData, s)
		fmt.Println(fileData)
	}
	return fileData
}

/*
func (o *ODTRepoDB) InsertToDb(modelId int, scooters []models.Scooter) error{

	valueStrings := make([]string, 0, len(scooters))
	valueArgs := make([]interface{}, 0, len(scooters) * 1)
	for i, scooter := range scooters {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i*1+1))
		valueArgs = append(valueArgs, scooter.SerialNumber)
	}

	stmt := fmt.Sprintf("INSERT INTO scooters(scooter_model, user_id, serial_number) VALUES %s", strings.Join(valueStrings, ","))
	if _, err := o.db.QueryExec(context.Background(),stmt, valueArgs...)
		err != nil {
		fmt.Println("Unable to insert due to: ", err)
		return err
	}
	return nil
}
 */