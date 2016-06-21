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
	Alias       string `short:"a" long:"alias" description:"alias to use from roles file"`
	Role        string `short:"r" long:"role" description:"arn role to use"`
	Session     string `short:"s" long:"session" description:"session name to use"`
	credentials *awsCredentials
	awsSession  string
	awsArn      string
}

type awsCredentials struct {
	AccessKey string
	SecretKey string
	Token     string
}

func (c *CmdAssume) Execute(args []string) error {
	if c.Alias != "" && c.Role == "" && c.Session == "" {
		c.awsArn = c.getArnOfAliasFromFile(rolesFile)
		c.awsSession = c.Alias
	} else if c.Alias == "" && c.Role != "" && c.Session != "" {
		c.awsArn = c.Role
		c.awsSession = c.Session
	}

	err := c.askForAWSCredentials()
	if err != nil {
		return err
	}

	err = c.saveAWSCredentials(awsCredentialsFile)
	if err != nil {
		return err
	}

	return nil
}

func (c *CmdAssume) askForAWSCredentials() error {
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

func (c *CmdAssume) getArnOfAliasFromFile(rolesFile string) string {
	file, _ := ini.Load(rolesFile)
	role, _ := file.GetSection(c.Alias)
	arn, err := role.GetKey("arn")
	if err != nil {
		log.Fatal(err)
	}

	return arn.String()
}

func (c *CmdAssume) saveAWSCredentials(credentialsFile string) error {
	cfg := ini.Empty()
	cfg.NewSection(c.awsSession)
	cfg.Section(c.awsSession).NewKey("aws_access_key_id", c.credentials.AccessKey)
	cfg.Section(c.awsSession).NewKey("aws_secret_access_key", c.credentials.SecretKey)
	cfg.Section(c.awsSession).NewKey("aws_security_token", c.credentials.Token)

	err := cfg.SaveTo(awsCredentialsFile)

	if err != nil {
		return err
	}

	return nil
}
