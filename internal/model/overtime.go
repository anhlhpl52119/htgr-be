package model

import "time"

type Overtime struct {
	ID             string    `json:"id" db:"id"`
	EmployeeID     string    `json:"employee_id" db:"employee_id"`
	RequestedHours float64   `json:"requested_hours" db:"requested_hours"`
	ActualHours    float64   `json:"actual_hours" db:"actual_hours"`
	Status         string    `json:"status" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateOvertimeRequest struct {
	EmployeeID     string  `json:"employee_id"`
	RequestedHours float64 `json:"requested_hours"`
}

type OvertimeResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
