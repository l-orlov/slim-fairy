package model

type (
	GetDietParams struct {
		Age              int
		Weight           int
		Height           int
		Gender           Gender
		PhysicalActivity PhysicalActivityLevel
		MealTimes        int
		SnackTimes       int
	}
)
