package internal

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/require"
)

func TestParseIdentifierInstanceID(t *testing.T) {
	identifier := "i-08d0906a5bd77e96f"
	expected := types.Filter{
		Name:   aws.String("instance-id"),
		Values: []string{identifier},
	}
	actual := parseIdentifier(identifier)
	require.Equal(t, expected, actual)
}

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

func TestParseIdentifierPublicDNS(t *testing.T) {
	identifier := "ec2-52-58-189-12.eu-central-1.compute.amazonaws.com"
	expected := types.Filter{
		Name:   aws.String("dns-name"),
		Values: []string{identifier},
	}
	actual := parseIdentifier(identifier)
	require.Equal(t, expected, actual)
}

func TestParseIdentifierPrivateDNS(t *testing.T) {
	identifier := "ip-172-31-33-49.eu-central-1.compute.internal"
	expected := types.Filter{
		Name:   aws.String("private-dns-name"),
		Values: []string{identifier},
	}
	actual := parseIdentifier(identifier)
	require.Equal(t, expected, actual)
}
