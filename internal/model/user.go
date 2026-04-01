package model

type User struct {
	ID          string  `json:"id" db:"id"`
	EmployeeID  string  `json:"employee_id" db:"employee_id"`
	Email       float64 `json:"email" db:"email"`
	ActualHours float64 `json:"actual_hours" db:"actual_hours"`
}
