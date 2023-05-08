package model

type (
	GetDietParams struct {
		Age              int                   `json:"age"`
		Weight           int                   `json:"weight"`
		Height           int                   `json:"height"`
		Gender           Gender                `json:"gender"`
		PhysicalActivity PhysicalActivityLevel `json:"physical_activity"`
		MealTimes        int                   `json:"meal_times"`
		SnackTimes       int                   `json:"snack_times"`
	}
)
