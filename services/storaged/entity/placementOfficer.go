package entity

type placementOfficer struct {
	ID        uint64
	CompanyID uint64
}

func NewplacementOfficer(id uint64, companyID uint64) (*placementOfficer, error) {
	placementOfficer := &placementOfficer{
		ID:        id,
		CompanyID: companyID,
	}

	return placementOfficer, nil
}
