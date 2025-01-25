package account_entity_http_v1

type CreateAccountResponse struct {
	AccountID      string `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
