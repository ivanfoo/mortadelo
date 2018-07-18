package commands

import (
	"fmt"
	"os"
)

var (
	awsCredentialsFile = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), ".aws", "credentials")
	aliasFile          = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), ".mortadelo", "alias")
	mortadeloDir       = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".mortadelo")
)
