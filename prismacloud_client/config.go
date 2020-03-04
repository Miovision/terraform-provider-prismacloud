package prismacloud_client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
}

func LoadCredentials() *Credentials {
	creds := new(Credentials)
	homeDir := os.Getenv("HOME")
	credsJsonFilePath := fmt.Sprintf("%s/.prismacloud/credentials", homeDir)
	credsJsonFile, _ := os.Open(credsJsonFilePath)
	json.NewDecoder(bufio.NewReader(credsJsonFile)).Decode(creds)
	return creds
}
