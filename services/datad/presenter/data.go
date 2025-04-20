package presenter

type DataRequest struct {
	JWT         string      `json:"jwt"`
	CompanyData interface{} `json:"company_data"`
}

type DataResponse struct {
	DataID      string      `json:"data_id"`
	CompanyData interface{} `json:"company_data"`
}
