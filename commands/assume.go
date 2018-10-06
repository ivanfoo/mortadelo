package commands

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"gopkg.in/ini.v1"
)

type AssumeCommand struct {
	Alias    string `short:"a" long:"alias" description:"alias to assume (session name if using --role)" required:"true"`
	Role     string `short:"r" long:"role" description:"literal arn role to assume instead of an alias" default:""`
	MFA      bool   `long:"mfa" description:"ask for a mfa token code"`
	Duration int64  `short:"d" long:"duration" description:"duration of the session in seconds" default:"900"`
	File     string `short:"f" long:"file" description:"alias file" default:"~/.mortadelo/alias"`
	Output   string `short:"o" long:"output" description:"file to store the AWS credentials" default:"~/.aws/credentials"`

	session     *session.Session
	credentials *awsCredentials

	roleArn string
}

type awsCredentials struct {
	AccessKey string
	SecretKey string
	Token     string
}

func (c *AssumeCommand) Execute(args []string) error {
	var err error
	var roleArn string

	if c.Role != "" {
		roleArn = c.Role
	} else {
		roleArn, err = c.getRoleArn()
		if err != nil {
			return err
		}
	}

	c.session = session.New()

	var serialNumber string
	var tokenCode string

	if c.MFA {
		serialNumber, err = c.getMFASerialNumber()
		if err != nil {
			return err
		}

		tokenCode = c.askForTokenCode()
	}

	c.credentials, err = c.getCredentials(serialNumber, tokenCode, roleArn)

	if err != nil {
		return err
	}

	err = c.saveCredentials()

	if err != nil {
		return err
	}

	fmt.Printf("temporary credentials saved in %s \n", c.Output)

	return nil
}

func (c *AssumeCommand) getMFASerialNumber() (string, error) {
	service := iam.New(c.session)
	userInfo, err := service.GetUser(&iam.GetUserInput{})

	if err != nil {
		return "", err
	}

	userName := *userInfo.User.UserName

	mfaDevices, err := service.ListMFADevices(&iam.ListMFADevicesInput{
		MaxItems: aws.Int64(1),
		UserName: aws.String(userName),
	})

	if err != nil {
		return "", err
	}

	return *mfaDevices.MFADevices[0].SerialNumber, nil
}

func (c *AssumeCommand) askForTokenCode() string {
	var tokenCode string

	fmt.Print("Insert MFA token: ")
	fmt.Scanf("%s", &tokenCode)

	return tokenCode
}

func (c *AssumeCommand) getCredentials(serialNumber string, tokenCode string, roleArn string) (*awsCredentials, error) {
	service := sts.New(c.session)
	role, err := service.AssumeRole(&sts.AssumeRoleInput{
		RoleSessionName: aws.String(c.Alias),
		DurationSeconds: aws.Int64(c.Duration),
		RoleArn:         aws.String(roleArn),
		SerialNumber:    aws.String(serialNumber),
		TokenCode:       aws.String(tokenCode),
	})

	if err != nil {
		return nil, err
	}

	credentials := &awsCredentials{
		AccessKey: *role.Credentials.AccessKeyId,
		SecretKey: *role.Credentials.SecretAccessKey,
		Token:     *role.Credentials.SessionToken,
	}

	return credentials, nil
}

func (c *AssumeCommand) getRoleArn() (string, error) {
	path := expandPath(c.File)
	file, err := ini.Load(path)

	if err != nil {
		return "", err
	}

	alias, err := file.GetSection(c.Alias)

	if err != nil {
		return "", err
	}

	arn, err := alias.GetKey("arn")

	if err != nil {
		return "", err
	}

	return arn.String(), nil
}

func (c *AssumeCommand) saveCredentials() error {
	path := expandPath(c.Output)
	cfg, err := ini.LooseLoad(path)

	cfg.NewSection(c.Alias)
	cfg.Section(c.Alias).NewKey("aws_access_key_id", c.credentials.AccessKey)
	cfg.Section(c.Alias).NewKey("aws_secret_access_key", c.credentials.SecretKey)
	cfg.Section(c.Alias).NewKey("aws_security_token", c.credentials.Token)
	cfg.Section(c.Alias).NewKey("aws_session_token", c.credentials.Token)

	err = cfg.SaveTo(path)

	if err != nil {
		return err
	}

	return nil
}
