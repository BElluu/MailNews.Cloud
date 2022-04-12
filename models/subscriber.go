package models

type Subscriber struct {
	UUID           string `json:"uuid"`
	Email          string `json:"email"`
	ActivateURL    string `json:"activate-url"`
	UnSubscribeURL string `json:"un-subscribe-url"`
	IsActive       bool   `json:"is-active"`
}
