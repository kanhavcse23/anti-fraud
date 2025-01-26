package account_mapper_v1

import (
	entityCoreV1Package "anti-fraud/account-service/entity/core/v1"
	entityDbV1Package "anti-fraud/account-service/entity/db/v1"
)

func AccountMapper(accountPayload *entityCoreV1Package.CreateAccountPayload) *entityDbV1Package.Account {
	return &entityDbV1Package.Account{
		DocumentNumber: accountPayload.DocumentNumber,
	}
}
