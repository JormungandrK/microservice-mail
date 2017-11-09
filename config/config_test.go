package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := `{
		"verificationURL": "http://kong:8000/users/verify",
		"mail": {
			"host": "fakesmtp",
			"port": "25",
			"user": "fake@email.com",
			"password": "password",
			"email": "dev@jormungandr.org"
		},
		"rabbitmq": {
			"username": "guest",
			"password": "guest",
			"host": "rabbitmq",
			"post": "5672"
		}
	}`

	cnfFile, err := ioutil.TempFile("", "tmp-config")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(cnfFile.Name())

	cnfFile.WriteString(config)

	cnfFile.Sync()

	loadedCnf, err := LoadConfig(cnfFile.Name())

	if err != nil {
		t.Fatal(err)
	}

	if loadedCnf == nil {
		t.Fatal("Configuration was not read")
	}
}
