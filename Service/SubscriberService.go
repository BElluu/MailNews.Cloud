package Service

import (
	"MailNews.Subscriber/Common/Validator"
	"errors"
)

func Subscribe(email string) (bool, error) {
	if EmailValidator.Validate(email) {
		return true, nil
	}
	return false, errors.New("email address is not valid")
}

func UnSubscribe() {

}
