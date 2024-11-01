deploy-localstack:
	make build-lambdas
	make deploy-lambdas
	make build-rds

deploy-lambdas:
	awslocal lambda create-function \
		--function-name GetUsers \
		--runtime provided.al2 \
		--zip-file fileb://cmd/users/function.zip \
		--handler users \
		--role arn:aws:iam::000000000000:role/lambda-role
	awslocal lambda create-function \
		--function-name GetMovies \
		--runtime provided.al2 \
		--zip-file fileb://cmd/movies/function.zip \
		--handler movies \
		--role arn:aws:iam::000000000000:role/lambda-role

build-rds:
	awslocal rds create-db-cluster \
		--db-cluster-identifier foundflix \
		--engine aurora-postgresql \
		--database-name movies \
		--master-username localhost \
		--master-user-password localhost
	awslocal rds create-db-instance \
		--db-instance-identifier db1-instance \
		--db-cluster-identifier foundflix \
		--engine aurora-postgresql \
		--db-instance-class db.t3.large
		
build-lambdas:
	cd cmd/users && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go
	cd cmd/users && zip -j function.zip ./bootstrap
	cd cmd/movies && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go
	cd cmd/movies && zip -j function.zip ./bootstrap

update-lambdas:
	make build-lambdas
	awslocal lambda update-function-code --function-name GetMovies --zip-file fileb://cmd/movies/function.zip
	rm -f cmd/users/bootstrap cmd/users/function.zip
	rm -f cmd/movies/bootstrap cmd/movies/function.zip

rebuild-lambdas:
	rm -f cmd/users/bootstrap cmd/users/function.zip
	rm -f cmd/movies/bootstrap cmd/movies/function.zip
	awslocal lambda delete-function --function-name GetUsers
	awslocal lambda delete-function --function-name GetMovies
	make deploy-lambdas

clean:
	rm -f cmd/users/bootstrap cmd/users/function.zip
	rm -f cmd/movies/bootstrap cmd/movies/function.zip
	awslocal lambda delete-function --function-name GetUsers
	awslocal lambda delete-function --function-name GetMovies
	awslocal rds delete-db-instance --db-instance-identifier db1-instance
	awslocal rds delete-db-cluster --db-cluster-identifier foundflix

test:
	awslocal lambda invoke \
		--function-name GetUsers \
		--payload "$(echo '{"httpMethod": "GET", "path": "/users"}' | base64)" \
		./tmp/get_users_response.json
	awslocal lambda invoke \
		--function-name GetMovies \
		--payload "$(echo '{"httpMethod": "GET", "path": "/movies"}' | base64)" \
		./tmp/get_movies_response.json

.PHONY: build-lambdas build-rds deploy-localstack deploy-lambdas clean test
