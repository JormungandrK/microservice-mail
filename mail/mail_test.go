package mail

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/JormungandrK/microservice-mail/config"
)

func TestParseTemplate(t *testing.T) {
	template := `{
		<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
		<html>

		</head>

		<body>
		<h1>Hello Jormungandr,</h1>

		<p>
			<a href="http://jormungandr/users/verify">Verify your registration</a>
		</p>
		</body>

		</html>
	}`

	templateFile, err := ioutil.TempFile("", "tmp-template.html")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(templateFile.Name())

	templateFile.WriteString(template)

	templateFile.Sync()

	tmpl, err := ParseTemplate(templateFile.Name(), map[string]string{})

	if err != nil {
		t.Fatal(err)
	}

	if tmpl == "" {
		t.Fatal("Failed to parse template file!")
	}
}

func TestSend(t *testing.T) {
	confBytes := []byte(`{
		"verificationURL": "http://kong:8000/users",
		"mail": {
			"host": "smtp.mailtrap.io",
			"port": "2525",
			"user": "7dfa3710bee1c3",
			"password": "68d8ccb96fb52b",
			"email": "e0e3decc9e-f10431@inbox.mailtrap.io"
		}
	}`)

	template := `
		<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
		<html>

		</head>

		<body>
		<h1>Hello Jormungandr,</h1>

		<p>
			<a href="http://jormungandr/users/verify">Verify your registration</a>
		</p>
		</body>

		</html>
	`

	cfg := &config.Config{}
	err := json.Unmarshal(confBytes, cfg)
	if err != nil {
		t.Fatal(err)
	}

	mailInfo := Info{
		"user-id",
		"jormungandr-test",
		"e0e3decc9e-f10431@inbox.mailtrap.io",
		"http://test/verify",
		"some-token",
	}

	err = Send(&mailInfo, cfg, template)
	if err != nil {
		t.Fatal(err)
	}
}
