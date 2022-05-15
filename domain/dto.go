package domain

type PostDto struct {
	IdFrom   string `json:"idFrom"`
	IdTo     string `json:"idTo"`
	Username string `json:"username"`
}

type HomepageFeedDto struct {
	Username string `json:"username"`
	Page     int    `json:"page"`
}
