package models

type ModelODT struct {
	ID          		int         `json:"id"`
	Price 				int 		`json:"price"`
	ModelName         	string 		`json:"model_name"`
	MaxWeight         	int    		`json:"max_weight"`
	Speed 			  	int    		`json:"speed"`
}

type ModelODTList struct {
	ModelsODT []ModelODT `json:"models_odt"`
}

type SupplierPricesODT struct {
	ID  			int `json:"id"`
	Price 			int `json:"price"`
	PaymentTypeID 	int `json:"payment_type_id"`
	UserId    		int `json:"user_id"`
}

type SupplierPricesODTList struct {
	SupplierPricesODT []SupplierPricesODT `json:"supplier_prices_odt_list"`
}
