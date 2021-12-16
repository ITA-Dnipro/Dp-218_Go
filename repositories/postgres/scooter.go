package postgres

import (
	"Dp218Go/models"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jszwec/csvutil"
	"io"
	"os"
	"strings"
)


func (s *SupplierRepoDB)GetAllScooters() (*models.ScooterList, error) {
	list := &models.ScooterList{}

/*	scooterModels, err := s.GetModels()
	if err != nil {
		return list, err
	}

 */

	querySQL := `SELECT * FROM scooters ORDER BY id DESC;`
	rows, err := s.db.QueryResult(context.Background(), querySQL)
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

	/*	scooter.ScooterModel, err = s.FindScooterList(scooterModels, scooterModelId)
		if err != nil {
			return list, err
		}

	 */

		list.Scooters = append(list.Scooters, scooter)
	}
	return list, nil
}

func (s *SupplierRepoDB) GetScooterByID(id int)(models.Scooter, error){
	scooter := models.Scooter{}

	querySQL := `SELECT id, model_id, owner_id, serial_number FROM scooters WHERE id = $1;`
	row := s.db.QueryResultRow(context.Background(), querySQL, id)

	var modelId int
	err := row.Scan(&scooter.ID, &modelId,  &userId, &scooter.SerialNumber)
	if err != nil {
		return models.Scooter{}, err
	}
//	scooter.ScooterModel, err = s.GetScooterModelById(modelId)

	return scooter, err
}

func (s *SupplierRepoDB) AddScooter(scooter *models.Scooter) error {
	var id int
	querySQL := `INSERT INTO scooters(model_id, owner_id, serial_number)
	   		VALUES($1, $2, $3)
	   		RETURNING id;`
	err := s.db.QueryResultRow(context.Background(), querySQL, scooter.ScooterModel.ID, scooter.User.ID, scooter.SerialNumber).Scan(&id)
	if err != nil {
		return err
	}
//	scooter.ID = id
	return nil
}

func (s *SupplierRepoDB) EditScooter(scooterId int, scooterData models.Scooter)(models.Scooter, error){
	scooter := models.Scooter{}
	querySQL := `UPDATE scooters
	   		SET model_id=$1, owner_id=$2, serial_number=$3
	   		WHERE id=$4
	   		RETURNING id, model_id, owner_id, serial_number;`
	var modelId int
	err := s.db.QueryResultRow(context.Background(), querySQL, scooterData.ScooterModel.ID, scooterData.User.ID, scooterData.SerialNumber, scooterId).Scan(
		&scooter.ID,  &modelId, &userId,  &scooter.SerialNumber)
	if err != nil {
		return scooter, err
	}

//	scooter.ScooterModel, err = s.SelectModel(modelId)
	if err != nil {
		return scooter, err
	}
	return scooter, nil
}

func (s *SupplierRepoDB) DeleteScooter(id int) error {
	querySQL := `DELETE FROM scooters WHERE id = $1;`
	_, err := s.db.QueryExec(context.Background(), querySQL, userId)
	return err
}

func (s *SupplierRepoDB)FindScooterList(scooterModel *models.ModelODTList, scooterModelId int ) (models.ModelODT, error) {
	for _, v := range scooterModel.ModelsODT {
		if v.ID == scooterModelId{
			return v, nil
		}
	}
	return models.ModelODT{}, fmt.Errorf("not found scooterModel id=%d", scooterModelId)
}

func (s *SupplierRepoDB) ConvertToStruct(path string) []models.Scooter{

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


func (s *SupplierRepoDB) InsertToDb(modelId int, scooters []models.Scooter) error{

	valueStrings := make([]string, 0, len(scooters))
	valueArgs := make([]interface{}, 0, len(scooters) * 1)
	for i, scooter := range scooters {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i*1+1))
		valueArgs = append(valueArgs, scooter.SerialNumber)
	}

	stmt := fmt.Sprintf("INSERT INTO scooters(scooter_model, user_id, serial_number) VALUES %s", strings.Join(valueStrings, ","))
	if _, err := s.db.QueryExec(context.Background(),stmt, valueArgs...)
		err != nil {
		fmt.Println("Unable to insert due to: ", err)
		return err
	}
	return nil
}
