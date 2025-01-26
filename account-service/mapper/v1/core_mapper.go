package account_mapper_v1

import (
	entityCoreV1Package "anti-fraud/account-service/entity/core/v1"
	entityHttpV1Package "anti-fraud/account-service/entity/http/v1"
)

func CreateAccountPayloadMapper(accountCreationRequest *entityHttpV1Package.CreateAccountRequest) *entityCoreV1Package.CreateAccountPayload {
	return &entityCoreV1Package.CreateAccountPayload{
		DocumentNumber: accountCreationRequest.DocumentNumber,
	}
}
