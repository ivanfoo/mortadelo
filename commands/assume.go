package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"gopkg.in/ini.v1"
)

var (
	awsCredentialsFile = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), ".aws", "credentials")
	rolesFile          = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), ".joe", "roles")
)

type CmdAssume struct {
	Role        string `short:"r" long:"role" description:"role to use" default:"default"`
	credentials *awsCredentials
	awsSession  string
	awsArn      string
}

type awsCredentials struct {
	AccessKey string
	SecretKey string
	Token     string
}

func NewCmdAssume() *CmdAssume {
	return &CmdAssume{}
}

func (c *CmdAssume) Execute(args []string) error {
	roleInfo, err := loadRolesFile(rolesFile).GetSection(c.Role)
	if err != nil {
		return err
	}

	arn, err := roleInfo.GetKey("arn")
	if err != nil {
		return err
	}

	c.awsSession = c.Role
	c.awsArn = arn.String()

	err = c.generateAWSCredentials()
	if err != nil {
		return err
	}

	c.saveAWSCredentials(awsCredentialsFile)
	if err != nil {
		return err
	}

	return nil
}

func (c *CmdAssume) generateAWSCredentials() error {
	svc := sts.New(session.New())
	resp, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(c.awsArn),
		RoleSessionName: aws.String(c.awsSession),
	})

	if err != nil {
		return err
	}

	c.credentials = &awsCredentials{
		AccessKey: *resp.Credentials.AccessKeyId,
		SecretKey: *resp.Credentials.SecretAccessKey,
		Token:     *resp.Credentials.SessionToken,
	}

	return nil
}

func loadRolesFile(rolesFile string) *ini.File {
	roles, err := ini.Load(rolesFile)

	if err != nil {
		log.Fatal(err)
	}

	return roles
}

func (c *CmdAssume) saveAWSCredentials(credentialsFile string) {
	cfg := ini.Empty()
	cfg.NewSection(c.awsSession)
	cfg.Section(c.awsSession).NewKey("aws_access_key_id", c.credentials.AccessKey)
	cfg.Section(c.awsSession).NewKey("aws_secret_access_key", c.credentials.SecretKey)
	cfg.Section(c.awsSession).NewKey("aws_security_token", c.credentials.Token)

	err := cfg.SaveTo(awsCredentialsFile)

	if err != nil {
		log.Fatal(err)
	}
}
