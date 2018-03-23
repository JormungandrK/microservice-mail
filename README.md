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
docker build -t jormungandrk/microservice-mail .
```

## Running the microservice
To run the service type:
```bash
docker run -it -e SERVICE_CONFIG_FILE=config.json jormungandrk/microservice-mail
```

## Service configuration

The service loads the  configuration from a JSON file /run/secrets/microservice_mail_config.json. To change the path set the
**SERVICE_CONFIG_FILE** env var.
Here's an example of a JSON configuration file:

```json
{
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
		"port": "5672"
	}
}
```

