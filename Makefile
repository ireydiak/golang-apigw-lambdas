deploy-localstack:
	make build-lambdas
	awslocal lambda create-function \
		--function-name GetUsers \
		--runtime provided.al2 \
		--zip-file fileb://cmd/users/function.zip \
		--handler users \
		--role arn:aws:iam::000000000000:role/lambda-role

build-lambdas:
	cd cmd/users && GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	cd cmd/users && zip -j function.zip ./bootstrap

clean:
	rm cmd/users/bootstrap cmd/users/function.zip
	awslocal lambda delete-function --function-name GetUsers

test:
	awslocal lambda invoke \
		--function-name GetUsers \
		--payload "$(echo '{"httpMethod": "GET", "path": "/users"}' | base64)" \
		response.json

.PHONY: build-lambdas deploy-localstack clean test
