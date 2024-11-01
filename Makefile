deploy-localstack:
	make build-lambdas
	awslocal lambda create-function \
		--function-name GetUsers \
		--runtime provided.al2 \
		--zip-file fileb://cmd/users/function.zip \
		--handler users \
		--role arn:aws:iam::000000000000:role/lambda-role
	make build-rds

build-rds:
	awslocal rds create-db-cluster \
		--db-cluster-identifier foundflix \
		--engine aurora-postgresql \
		--database-name movies \
		--master-username foundflix \
		--master-user-password foundflix

build-lambdas:
	cd cmd/users && GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	cd cmd/users && zip -j function.zip ./bootstrap

clean:
	rm -f cmd/users/bootstrap cmd/users/function.zip
	awslocal lambda delete-function --function-name GetUsers
	awslocal rds delete-db-cluster --db-cluster-identifier foundflix

test:
	awslocal lambda invoke \
		--function-name GetUsers \
		--payload "$(echo '{"httpMethod": "GET", "path": "/users"}' | base64)" \
		response.json

.PHONY: build-lambdas build-rds deploy-localstack clean test
