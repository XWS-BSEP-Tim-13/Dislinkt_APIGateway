package domain

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/domain/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	CompanyName string      `json:"company_name"`
	Description string      `json:"description"`
	Location    string      `json:"location"`
	Website     string      `json:"website"`
	CompanySize string      `json:"company_size"`
	Industry    string      `json:"industry"`
}

type Company struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	CompanyName string             `bson:"company_name" validate:"required,companyName"`
	Username    string             `bson:"username" validate:"required,username"`
	Email       string             `bson:"email" validate:"required,email"`
	PhoneNumber string             `bson:"phone_number" validate:"required,numeric,min=9,max=10"`
	Description string             `bson:"description"`
	Location    string             `bson:"location" validate:"required,max=256"`
	Website     string             `bson:"website" validate:"required,website"`
	CompanySize string             `bson:"company_size" validate:"required,companyName"`
	Industry    string             `bson:"industry" validate:"required,max=256"`
	IsActive    bool               `bson:"is_active"`
}

type JobOffer struct {
	Id             primitive.ObjectID  `bson:"_id"`
	Position       string              `bson:"position"`
	JobDescription string              `bson:"job_description"`
	Prerequisites  string              `bson:"prerequisites"`
	Company        Company             `bson:"company"`
	EmploymentType enum.EmploymentType `bson:"employment_type"`
	Published      time.Time           `bson:"published" validate:"required"`
}
