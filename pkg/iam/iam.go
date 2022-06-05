package iam

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

func RoleARN(account, role string) string {
	if arn.IsARN(role) {
		return role
	}

	return fmt.Sprintf("arn:aws:iam::%s:role/%s", account, role)
}
