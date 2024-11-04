# Golang AWS ApiGateway + Lambda

This is a simple project that attempts to answer the architectural challenge of combining both 
a simple and somewhat cloud agnostic development environment with the popular ApiGateway+Lambda
combination to deploy serverless APIs on AWS.

# Dependencies

- Golang >= 1.23.1
- Localstack = 3.8.1
- awscli = 2.18.15

# Installation

Build the lambdas on Localstack and compile the Go endpoints.

If you use localstack for local development, make sure your have a `localstack` entry in your `$HOME/.aws/credentials` file with the following content:

```bash
[localstack]
aws_access_key_id=test
aws_secret_access_key=test
region=us-east-1
```

Otherwise you might encounter authentication errors with your regular AWS SSO sessions.

```bash
AWS_PROFILE=localstack && make deploy-localstack
```

Test the endpoints

```bash
make test
```

