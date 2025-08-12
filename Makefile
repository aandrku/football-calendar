include .env

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## build: build binary
.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap ./cmd/function

## zip: create a zip file ready for deployment
.PHONY: zip
zip:
	zip footballFixtures.zip bootstrap
	
## create: create lambda function on AWS
.PHONY: create
create: build zip
	aws lambda create-function --function-name footballFixtures \
		--runtime provided.al2023 --handler bootstrap \
		--role arn:aws:iam::318112817111:role/footballFixtures \
		--zip-file fileb://footballFixtures.zip

## update: update lambda function code on AWS
.PHONY: update
update: build zip
	aws lambda update-function-code --function-name footballFixtures \
		--zip-file fileb://footballFixtures.zip

## update-config: update environment variable configuration
.PHONY: update-config
update-config: 
	aws lambda update-function-configuration \
		--function-name footballFixtures \
		--environment file://.env.json	

## invoke: invoke funcion on the cloud
.PHONY: invoke
invoke:
	aws lambda invoke --function-name footballFixtures response.json


## delete: delete function from the cloud
.PHONY: delete
delete: 
	aws lambda delete-function \
		--function-name footballFixtures

## clean: remove build files produced by build and zip
.PHONY: clean
clean:
	rm bootstrap footballFixtures.zip


