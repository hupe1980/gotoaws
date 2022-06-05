package iam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleARN(t *testing.T) {
	expected := "arn:aws:iam::1234567890:role/cluster-admin"

	t.Run("only roleName", func(t *testing.T) {
		roleARN := RoleARN("1234567890", "cluster-admin")
		assert.Equal(t, expected, roleARN)
	})

	t.Run("roleArn", func(t *testing.T) {
		roleARN := RoleARN("1234567890", expected)
		assert.Equal(t, expected, roleARN)
	})
}
