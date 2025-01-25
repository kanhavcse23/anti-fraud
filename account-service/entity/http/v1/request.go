package account_entity_http_v1

import "errors"

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

func (createAccountRequest *CreateAccountRequest) Validate() error {
	if createAccountRequest.DocumentNumber == "" {
		return errors.New("document number should not be empty")
	}
	return nil
}
