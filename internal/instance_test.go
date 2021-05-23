package internal

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/require"
)

func TestParseIdentifierNameTag(t *testing.T) {
	identifier := "myserver"
	expected := types.Filter{
		Name:   aws.String("tag:Name"),
		Values: []string{identifier},
	}
	actual := parseIdentifier(identifier)
	require.Equal(t, expected, actual)
}

func TestParseIdentifierPublicIP(t *testing.T) {
	identifier := "80.1.2.3"
	expected := types.Filter{
		Name:   aws.String("ip-address"),
		Values: []string{identifier},
	}
	actual := parseIdentifier(identifier)
	require.Equal(t, expected, actual)
}

func TestParseIdentifierPrivateIP(t *testing.T) {
	identifier := "10.1.0.12"
	expected := types.Filter{
		Name:   aws.String("private-ip-address"),
		Values: []string{identifier},
	}
	actual := parseIdentifier(identifier)
	require.Equal(t, expected, actual)
}
