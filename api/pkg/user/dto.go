package user

import "github.com/go-ozzo/ozzo-validation/v4"

type SignUpInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (i SignUpInfo) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Username, validation.Required, validation.Length(2, 50)),
		validation.Field(&i.Password, validation.Required, validation.Length(6, 255)),
	)
}
