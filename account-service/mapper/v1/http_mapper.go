package account_mapper_v1

import (
	entityDbV1Package "anti-fraud/account-service/entity/db/v1"
	entityHttpV1Package "anti-fraud/account-service/entity/http/v1"
	"strconv"
)

func AccountDetailsResponseMapper(account *entityDbV1Package.Account) *entityHttpV1Package.CreateAccountResponse {
	return &entityHttpV1Package.CreateAccountResponse{
		AccountID:      strconv.FormatUint(uint64(account.ID), 10),
		DocumentNumber: account.DocumentNumber,
	}
}
