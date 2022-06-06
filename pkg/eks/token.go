package eks

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"time"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/hupe1980/gotoaws/pkg/config"
	"github.com/hupe1980/gotoaws/pkg/iam"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/apis/clientauthentication"
	clientauthv1beta1 "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

const (
	// The actual token expiration (presigned STS urls are valid for 15 minutes after timestamp in x-amz-date).
	presignedURLExpiration = 15 * time.Minute
	clusterNameHeader      = "x-k8s-aws-id"
	v1Prefix               = "k8s-aws-v1."
	kindExecCredential     = "ExecCredential" // nolint: gosec // no hardcoded credential
	execInfoEnvKey         = "KUBERNETES_EXEC_INFO"
)

// Token is generated and used by Kubernetes client-go to authenticate with a Kubernetes cluster.
type Token struct {
	Token      string
	Expiration time.Time
}

type TokenGen interface {
	// Get a token using credentials in the default credentials chain.
	Get(clusterName string) (*Token, error)

	// GetWithRole creates a token by assuming the provided role, using the credentials in the default chain.
	GetWithRole(clusterName, role string) (*Token, error)

	// FormatJSON returns the client auth formatted json for the ExecCredential auth.
	FormatJSON(Token) string
}

type tokenGen struct {
	client  *sts.Client
	region  string
	account string
}

// NewTokenGen creates a TokenGen and returns it.
func NewTokenGen(cfg *config.Config) TokenGen {
	return &tokenGen{
		client:  sts.NewFromConfig(cfg.AWSConfig),
		region:  cfg.Region,
		account: cfg.Account,
	}
}

// Get uses the directly available AWS credentials to return a token valid for the clusterName.
func (t *tokenGen) Get(clusterName string) (*Token, error) {
	return t.get(t.client, clusterName)
}

// GetWithRole assumes the given AWS IAM role and returns a token valid for the clusterName.
func (t *tokenGen) GetWithRole(clusterName, role string) (*Token, error) {
	prov := stscreds.NewAssumeRoleProvider(t.client, iam.RoleARN(t.account, role), func(aro *stscreds.AssumeRoleOptions) {
		aro.RoleSessionName = "EKSGetTokenAuth"
	})

	config, err := aws_config.LoadDefaultConfig(
		context.TODO(),
		aws_config.WithRegion(t.region),
		aws_config.WithCredentialsProvider(prov),
	)
	if err != nil {
		return nil, err
	}

	return t.get(sts.NewFromConfig(config), clusterName)
}

func (t *tokenGen) get(client *sts.Client, clusterName string) (*Token, error) {
	presignClient := sts.NewPresignClient(client, func(po *sts.PresignOptions) {
		po.ClientOptions = []func(*sts.Options){
			sts.WithAPIOptions(smithyhttp.SetHeaderValue("X-Amz-Expires", "60")),
			sts.WithAPIOptions(smithyhttp.SetHeaderValue(clusterNameHeader, clusterName)),
		}
	})

	getCallerIdentity, err := presignClient.PresignGetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	// Set token expiration to 1 minute before the presigned URL expires for some cushion
	tokenExpiration := time.Now().Local().Add(presignedURLExpiration - 1*time.Minute)

	return &Token{
		Token:      v1Prefix + base64.RawURLEncoding.EncodeToString([]byte(getCallerIdentity.URL)),
		Expiration: tokenExpiration,
	}, nil
}

// FormatJSON formats the json to support ExecCredential authentication
func (t *tokenGen) FormatJSON(token Token) string {
	apiVersion := clientauthv1beta1.SchemeGroupVersion.String()

	env := os.Getenv(execInfoEnvKey)
	if env != "" {
		cred := &clientauthentication.ExecCredential{}
		if err := json.Unmarshal([]byte(env), cred); err == nil {
			apiVersion = cred.APIVersion
		}
	}

	expirationTimestamp := metav1.NewTime(token.Expiration)
	execInput := &clientauthv1beta1.ExecCredential{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiVersion,
			Kind:       kindExecCredential,
		},
		Status: &clientauthv1beta1.ExecCredentialStatus{
			ExpirationTimestamp: &expirationTimestamp,
			Token:               token.Token,
		},
	}
	enc, _ := json.Marshal(execInput)

	return string(enc)
}

func getToken(cfg *config.Config, clusterName, role string) (*Token, error) {
	gen := NewTokenGen(cfg)

	if role != "" {
		return gen.GetWithRole(clusterName, role)
	}

	return gen.Get(clusterName)
}
