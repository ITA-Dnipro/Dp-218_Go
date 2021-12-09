package postgres

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jszwec/csvutil"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

type FileRepoDB struct {
	db repositories.AnyDatabase
}

func NewFileRepoDB(db repositories.AnyDatabase) *FileRepoDB {
	return &FileRepoDB{db}
}

func (f FileRepoDB) CreateTempFile(file multipart.File)string{
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile, err := ioutil.TempFile("./../../internal/temp_files", "upload-*.сsv")
	if err != nil {
		fmt.Println(err)
	}
	//	defer tempFile.Close()
	tempFile.Write(fileBytes)
	return tempFile.Name()
}

func (f FileRepoDB) ConvertToStruct(path string)[]models.ScooterUploaded {
	csvFile, _ := os.Open(path)
	reader := csv.NewReader(csvFile)
	reader.Comma = ';'

	scooterHeader, _ := csvutil.Header(models.ScooterUploaded{}, "csv")
	dec, _ := csvutil.NewDecoder(reader, scooterHeader...)

	var fileData []models.ScooterUploaded
	for {
		var s models.ScooterUploaded
		if err := dec.Decode(&s); err == io.EOF {
			break
		}
		fileData = append(fileData, s)
	}
	return fileData
}

func (f FileRepoDB) InsertScooterData(scooterModelId models.ScooterModel, scooterUploaded []models.ScooterUploaded)error{
	valueStrings := make([]string, 0, len(scooterUploaded))
	model := make([]interface{}, 0, len(scooterUploaded) * 3)

	for i, scooterModel := range scooterUploaded {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*3+1, i*3+2, i*3+3, i*3+4))
		model = append(model, scooterModel.Model.ID)
		model = append(model, scooterModel.Scooter.SerialNumber)
		model = append(model, scooterModel.Owner.ID)
	}

	stmt := fmt.Sprintf("INSERT INTO scooter_models (payment_type_id, model_name, max_weight, speed) VALUES %s ON CONFLICT (model_name) DO NOTHING;", strings.Join(valueStrings, ","))
	if _, err := f.db.QueryExec(context.Background(),stmt, model...)
		err != nil {
		fmt.Println("Unable to insert due to: ", err)
		return err
	}
	return nil
}

