package domain

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/domain/enum"
	"time"
)

type RegisterRequest struct {
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	Role        string      `json:"role"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Email       string      `json:"email"`
	PhoneNumber string      `json:"phone_number"`
	Gender      enum.Gender `json:"gender"`
	DateOfBirth time.Time   `json:"date_of_birth"`
	Biography   string      `json:"biography"`
	IsPrivate   bool        `json:"is_private"`
	IsCompany   bool        `json:"is_company"`
}
