package model

type (
	MealsAndSnacksNumber struct {
		MealsNumberPerDay  int `json:"meals"`
		SnacksNumberPerDay int `json:"snacks"`
	}
	GetDietParams struct {
		Age              int                   `json:"age"`
		Weight           int                   `json:"weight"`
		Height           int                   `json:"height"`
		Gender           Gender                `json:"gender"`
		PhysicalActivity PhysicalActivityLevel `json:"physical_activity"`
		MealsAndSnacksNumber
	}
)
