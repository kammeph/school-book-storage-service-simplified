package schools

import (
	"time"

	"github.com/google/uuid"
)

type SchoolModel struct {
	ID        string
	Name      string
	Active    bool
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt *time.Time
	UpdatedBy *string
}

func NewSchool(name, createdBy string) SchoolModel {
	return SchoolModel{
		ID:        uuid.NewString(),
		Name:      name,
		Active:    true,
		CreatedAt: time.Now(),
		CreatedBy: createdBy,
	}
}

func NewSchoolWithId(id, name, updatedBy string) SchoolModel {
	return SchoolModel{
		ID:        id,
		Name:      name,
		Active:    true,
		UpdatedBy: &updatedBy,
	}
}

type SchoolDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AddSchoolDto struct {
	Name string `json:"name"`
}

type UpdateSchoolDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
