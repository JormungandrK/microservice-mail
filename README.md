# microservice-email

[![Build](https://travis-ci.com/Microkubes/microservice-mail.svg?token=CxRhMM58BLTPkR2pgxqw&branch=master)](https://travis-ci.com/Microkubes/microservice-mail)
[![Test Coverage](https://api.codeclimate.com/v1/badges/9cee97cb8b85f66c7185/test_coverage)](https://codeclimate.com/repos/5a01ff2cf1c05e02dd000014/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/9cee97cb8b85f66c7185/maintainability)](https://codeclimate.com/repos/5a01ff2cf1c05e02dd000014/maintainability)

Sending of activation mail to the user

## Prerequisite
Create a project directory. Set GOPATH enviroment variable to that project. Add $GOPATH/bin to the $PATH
```
export GOPATH=/path/to/project
export PATH=$GOPATH/bin:$PATH
```

## Compile and run the service:
Clone the repo:
```
cd $GOPATH/src
git clone git@github.com:Microkubes/microservice-mail.git
```

Install the dependencies:
```bash
dep ensure -v
```

Then compile and run:
```
cd microservice-mail
go build -o microservice-mail
./microservice-mail
```

## Tests
From root of the project run:
```
go test -v $(go list ./... | grep -v vendor)
```

## Docker Image
To build the docker image run:
```bash
docker build -t microkubes/microservice-mail .
```

## Running the microservice
To run the service type:
```bash
docker run -it -e SERVICE_CONFIG_FILE=config.json microkubes/microservice-mail
```

## Service configuration

The service loads the  configuration from a JSON file /run/secrets/microservice_mail_config.json. To change the path set the
**SERVICE_CONFIG_FILE** env var.

For testing purposes you may want to enable sending the email notifications to a local SMTP server over unencrypted connection.
To allow this, you must set:

```bash
export ALLOW_UNENCRYPTED_CONNECTION=true
```

**Note**: Be careful not to allow sending email over unencrypted connection in production mode!

Here's an example of a JSON configuration file:

```json
{
	"templatesBaseLocation": "./public/template/",
	"templates": {
		"templateName": {
			"filename" : "template-filename.html",
			"subject": "Subject of the mail",
			"data": {
				"example": "value"
			}
		}
	},
	"mail": {
		"host": "fakesmtp",
		"port": "25",
		"user": "fake@email.com",
		"password": "password",
		"email": "dev@microkubes.org"
	},
	"rabbitmq": {
		"username": "guest",
		"password": "guest",
		"host": "rabbitmq",
		"port": "5672"
	}
}
```
## How to use
Send message to AMQP topic "email-queue" formatted in the following structure
```
type AMQPMessage struct {
	Email        string            `json:"email,omitempty"`
	Template     string            `json:"template,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
}
```

Email json["email"] - "Send To" Email

Template json["template"] - name of the template defined in config.json

Data json["data"] - map with all properties that are required by the template

## Contributing

For contributing to this repository or its documentation, see the [Contributing guidelines](CONTRIBUTING.md).

