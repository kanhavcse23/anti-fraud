package account_mapper_v1

import (
	entityCoreV1Path "anti-fraud/account-service/entity/core/v1"
	entityDbV1Path "anti-fraud/account-service/entity/db/v1"
)

func AccountMapper(accountPayload *entityCoreV1Path.CreateAccountPayload) *entityDbV1Path.Account {
	return &entityDbV1Path.Account{
		DocumentNumber: accountPayload.DocumentNumber,
	}
}
