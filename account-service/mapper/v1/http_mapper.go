package mapper_v1

import (
	entityDbV1Path "anti-fraud/account-service/entity/db/v1"
	entityHttpV1Path "anti-fraud/account-service/entity/http/v1"
	"strconv"
)

func AccountDetailsResponseMapper(account *entityDbV1Path.Account) *entityHttpV1Path.CreateAccountResponse {
	return &entityHttpV1Path.CreateAccountResponse{
		AccountID:      strconv.FormatUint(uint64(account.ID), 10),
		DocumentNumber: account.DocumentNumber,
	}
}
