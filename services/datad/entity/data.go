package entity

import "github.com/google/uuid"

type Data struct {
	DataID      string
	CompanyData interface{}
}

func NewData(companyData interface{}) (*Data) {
	data := &Data{
		DataID:      uuid.New().String(),
		CompanyData: companyData,
	}
	return data
}
