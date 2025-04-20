package data

import "backend/services/datad/entity"

type Repository interface {
	Writer
	Reader
}

type Writer interface {
	CreateData(data *entity.Data) (string, error)
}

type Reader interface {
	GetDataByID(id string) (*entity.Data, error)
}

type Usecase interface {
	CreateData(jwt string, companyData interface{}) (*entity.Data, error)
	GetDataByID(id string) (*entity.Data, error)
}
