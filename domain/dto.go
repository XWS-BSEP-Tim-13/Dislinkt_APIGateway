package domain

import "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/domain/enum"

type PostDto struct {
	IdFrom   string `json:"idFrom"`
	IdTo     string `json:"idTo"`
	Username string `json:"username"`
}

type CreatePostDto struct {
	Content string `json:"content"`
	Image   string `json:"image"`
}

type HomepageFeedDto struct {
	Username string `json:"username"`
	Page     int    `json:"page"`
}

type JobOfferTokenDto struct {
	JobOffer JobOfferDto `json:"jobOffer"`
	Token    string      `json:"token"`
}

type MessageDto struct {
	MessageFrom string `bson:"message_from"`
	MessageTo   string `bson:"message_to"`
	Content     string `bson:"content"`
}

type JobOfferDto struct {
	Position       string              `bson:"position"`
	JobDescription string              `bson:"job_description"`
	Prerequisites  string              `bson:"prerequisites"`
	Company        Company             `bson:"company"`
	EmploymentType enum.EmploymentType `bson:"employment_type"`
}
