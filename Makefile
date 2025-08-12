include .env

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap ./cmd/function

.PHONY: zip
zip:
	zip footballFixtures.zip bootstrap
	
.PHONY: create
create: build zip
	aws lambda create-function --function-name footballFixtures \
		--runtime provided.al2023 --handler bootstrap \
		--role arn:aws:iam::318112817111:role/footballFixtures \
		--zip-file fileb://footballFixtures.zip

.PHONY: update
update: build zip
	aws lambda update-function-code --function-name footballFixtures \
		--zip-file fileb://footballFixtures.zip

.PHONY: update-config
update-config: 
	aws lambda update-function-configuration \
		--function-name footballFixtures \
		--environment file://.env.json	

.PHONY: invoke
invoke:
	aws lambda invoke --function-name footballFixtures


.PHONY: delete
delete: 
	aws lambda delete-function \
		--function-name footballFixtures

.PHONY: clean
clean:
	rm bootstrap footballFixtures.zip


