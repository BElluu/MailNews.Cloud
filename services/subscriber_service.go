package services

import (
	"MailNews.Subscriber/common/validator"
	"errors"
)

func Subscribe(email string) (bool, error) {
	if Email.IsValid(email) {
		return true, nil
	}
	return false, errors.New("email address is not valid")
}

func UnSubscribe() {

}
