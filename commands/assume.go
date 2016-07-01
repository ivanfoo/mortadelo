package commands

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"gopkg.in/ini.v1"
)

type CmdAssume struct {
	Alias   string `short:"a" long:"alias" description:"alias to use from roles file"`
	Role    string `short:"r" long:"role" description:"arn role to use"`
	Session string `short:"s" long:"session" description:"session name to use"`
	MFA     bool   `long:"mfa" description:"use mfa token"`

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
		err = c.setAwsCredentialsWithMFA()
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

	fmt.Println("Temporary credentials saved in ~/.aws/credentials")

	return nil
}

func (c *CmdAssume) validate() error {
	if c.Alias != "" && c.Role == "" && c.Session == "" {
		c.Role = c.getArnFromAliasFile(aliasFile)
		c.Session = c.Alias
		return nil
	} else if c.Alias == "" && c.Role != "" && c.Session != "" {
		return nil
	}

	return fmt.Errorf("You must use either -a or -r flag")
}

func (c *CmdAssume) setAwsCredentialsWithMFA() error {
	svc := sts.New(session.New())

	identity, _ := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	mfaSerialNumber := c.getMfaSerialNumber(*identity.Arn)

	var mfaTokenCode string

	fmt.Print("Insert MFA code: ")
	fmt.Scanf("%s", &mfaTokenCode)

	resp, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(c.Role),
		RoleSessionName: aws.String(c.Session),
		SerialNumber:    aws.String(mfaSerialNumber),
		TokenCode:       aws.String(mfaTokenCode),
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

func (c *CmdAssume) setAWSCredentials() error {
	svc := sts.New(session.New())

	resp, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(c.Role),
		RoleSessionName: aws.String(c.Session),
	})

	if err != nil {
		return fmt.Errorf("Authentication failed. Maybe you should use a MFA token?")
	}

	c.credentials = &awsCredentials{
		AccessKey: *resp.Credentials.AccessKeyId,
		SecretKey: *resp.Credentials.SecretAccessKey,
		Token:     *resp.Credentials.SessionToken,
	}

	return nil
}

func (c *CmdAssume) getArnFromAliasFile(aliasFile string) string {
	file, err := ini.Load(aliasFile)

	if err != nil {
		log.Fatal(aliasFile + " not found")
	}

	role, _ := file.GetSection(c.Alias)
	arn, err := role.GetKey("arn")

	if err != nil {
		log.Fatal(err)
	}

	return arn.String()
}

func (c *CmdAssume) getMfaSerialNumber(arn string) (string, error) {
	svc := iam.New(session.New())
	re := regexp.MustCompile(`iam\:\:\d+\:\w+\/([\w=\-\,\.\=\@]+)`)
	userName := re.FindStringSubmatch(arn)[1]

	mfaDevice, err := svc.ListMFADevices(&iam.ListMFADevicesInput{
		MaxItems: aws.Int64(1),
		UserName: aws.String(userName),
	})

	if err != nil {
		fmt.Errorf("Error getting mfa serial number")
	}

	return *mfaDevice.MFADevices[0].SerialNumber, nil
}

func (c *CmdAssume) saveAWSCredentials(credentialsFile string) error {
	cfg, err := ini.LooseLoad(awsCredentialsFile)

	if err != nil {
		fmt.Println("Creating new credentials file...")
	}

	cfg.NewSection(c.Session)
	cfg.Section(c.Session).NewKey("aws_access_key_id", c.credentials.AccessKey)
	cfg.Section(c.Session).NewKey("aws_secret_access_key", c.credentials.SecretKey)
	cfg.Section(c.Session).NewKey("aws_security_token", c.credentials.Token)

	err = cfg.SaveTo(awsCredentialsFile)

	if err != nil {
		return err
	}

	return nil
}
