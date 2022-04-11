package models

type Subscriber struct {
	//SubId          int    `json:"sub-id"`
	Email          string `json:"email"`
	ActivateURL    string `json:"activate-url"`
	UnSubscribeURL string `json:"un-subscribe-url"`
	IsActive       bool   `json:"is-active"`
}
