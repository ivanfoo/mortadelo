package commands

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"gopkg.in/ini.v1"
)

var (
	err           error
	maxMfaDevices int64 = 1
)

type CmdAssume struct {
	Alias    string `short:"a" long:"alias" description:"alias to use from roles file"`
	Role     string `short:"r" long:"role" description:"arn role to use"`
	Session  string `short:"s" long:"session" description:"session name to use"`
	MFA      bool   `long:"mfa" description:"ask for a mfa token"`
	Duration int64  `short:"d" long:"duration" description:"duration of the session in seconds" default:"900"`

	credentials *awsCredentials
}

type awsCredentials struct {
	AccessKey string
	SecretKey string
	Token     string
}

func (c *CmdAssume) Execute(args []string) error {
	err := c.validate()

	if err != nil {
		return err
	}

	if c.MFA {
		err = c.setAWSCredentialsUsingMFA()
	} else {
		err = c.setAWSCredentials()
	}

	if err != nil {
		return err
	}

	err = c.saveAWSCredentials(awsCredentialsFile)

	if err != nil {
		return err
	}

	fmt.Println("temporary credentials saved in ~/.aws/credentials")

	return nil
}

func (c *CmdAssume) validate() error {
	if c.Alias != "" && c.Role == "" && c.Session == "" {
		c.Role, err = c.getArnFromAliasFile(aliasFile)
		if err != nil {
			return err
		}
		c.Session = c.Alias
		return nil
	} else if c.Alias == "" && c.Role != "" && c.Session != "" {
		return nil
	}

	return fmt.Errorf("use either -a or -r flag")
}

func (c *CmdAssume) setAWSCredentials() error {
	session := sts.New(session.New())

	var (
		mfaTokenCode    string
		mfaSerialNumber string
	)

	resp, err := session.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         &c.Role,
		DurationSeconds: &c.Duration,
		RoleSessionName: &c.Session,
		SerialNumber:    &mfaSerialNumber,
		TokenCode:       &mfaTokenCode,
	})

	if err != nil {
		return fmt.Errorf("unable to get valid aws credentials")
	}

	c.credentials = &awsCredentials{
		AccessKey: *resp.Credentials.AccessKeyId,
		SecretKey: *resp.Credentials.SecretAccessKey,
		Token:     *resp.Credentials.SessionToken,
	}

	return nil
}

func (c *CmdAssume) setAWSCredentialsUsingMFA() error {
	session := sts.New(session.New())
	identity, _ := session.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	mfaSerialNumber, err := c.getMfaSerialNumber(*identity.Arn)

	if err != nil {
		return fmt.Errorf("unable to get a valid mfa serial number")
	}

	var mfaTokenCode string

	fmt.Print("Insert MFA token: ")
	fmt.Scanf("%s", &mfaTokenCode)

	resp, err := session.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         &c.Role,
		RoleSessionName: &c.Session,
		DurationSeconds: &c.Duration,
		SerialNumber:    &mfaSerialNumber,
		TokenCode:       &mfaTokenCode,
	})

	if err != nil {
		return fmt.Errorf("unable to get valid aws credentials")
	}

	c.credentials = &awsCredentials{
		AccessKey: *resp.Credentials.AccessKeyId,
		SecretKey: *resp.Credentials.SecretAccessKey,
		Token:     *resp.Credentials.SessionToken,
	}

	return nil
}

func (c *CmdAssume) getArnFromAliasFile(aliasFile string) (string, error) {
	file, err := ini.Load(aliasFile)

	if err != nil {
		return "", fmt.Errorf("unable to open " + aliasFile)
	}

	alias, err := file.GetSection(c.Alias)

	if err != nil {
		return "", fmt.Errorf("alias not found in file")
	}

	arn, err := alias.GetKey("arn")

	if err != nil {
		return "", fmt.Errorf("malformed alias file: missing arn key")
	}

	return arn.String(), nil
}

func (c *CmdAssume) getMfaSerialNumber(arn string) (string, error) {
	session := iam.New(session.New())
	re := regexp.MustCompile(`iam\:\:\d+\:\w+\/([\w=\-\,\.\=\@]+)`)
	userName := re.FindStringSubmatch(arn)[1]
	mfaDevice, err := session.ListMFADevices(&iam.ListMFADevicesInput{
		MaxItems: &maxMfaDevices,
		UserName: &userName,
	})

	if err != nil {
		return "", fmt.Errorf("unable to get a valid mfa serial number")
	}

	return *mfaDevice.MFADevices[0].SerialNumber, nil
}

func (c *CmdAssume) saveAWSCredentials(credentialsFile string) error {
	cfg, err := ini.LooseLoad(awsCredentialsFile)

	if err != nil {
		fmt.Println("creating new credentials file...")
	}

	cfg.NewSection(c.Session)
	cfg.Section(c.Session).NewKey("aws_access_key_id", c.credentials.AccessKey)
	cfg.Section(c.Session).NewKey("aws_secret_access_key", c.credentials.SecretKey)
	cfg.Section(c.Session).NewKey("aws_security_token", c.credentials.Token)
	cfg.Section(c.Session).NewKey("aws_session_token", c.credentials.Token)

	err = cfg.SaveTo(awsCredentialsFile)

	if err != nil {
		return fmt.Errorf("saving credentials file failed")
	}

	return nil
}
