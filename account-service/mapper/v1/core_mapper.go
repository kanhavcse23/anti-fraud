package account_mapper_v1

import (
	entityCoreV1Path "anti-fraud/account-service/entity/core/v1"
	entityHttpV1Path "anti-fraud/account-service/entity/http/v1"
)

func CreateAccountPayloadMapper(accountCreationRequest *entityHttpV1Path.CreateAccountRequest) *entityCoreV1Path.CreateAccountPayload {
	return &entityCoreV1Path.CreateAccountPayload{
		DocumentNumber: accountCreationRequest.DocumentNumber,
	}
}
