ACCOUNT_ID="000000000000"
PATH_PART="movies"
FUNCTION_ARN=$(awslocal lambda list-functions --profile localstack | \
    jq -r '.Functions[] | select(.FunctionName == "GetMovies") | .FunctionArn'
)

API_ID=$(awslocal apigateway create-rest-api --name foundflix --profile localstack | jq -r '.id')
echo "API_ID: $API_ID"

API_ENDPOINT_ID=$(awslocal apigateway get-resources --rest-api-id $API_ID --profile localstack | jq -r '.items[0].id')
echo "API_ENDPOINT_ID: $API_ENDPOINT_ID"

API_RESOURCE_ID=$(awslocal apigateway create-resource \
    --rest-api-id $API_ID \
    --parent-id $API_ENDPOINT_ID \
    --profile localstack \
    --path-part $PATH_PART | jq -r '.id')
echo "API_RESOURCE_ID: $API_RESOURCE_ID"

awslocal apigateway put-method \
    --rest-api-id $API_ID \
    --resource-id $API_RESOURCE_ID \
    --http-method ANY \
    --authorization-type NONE \
    --profile localstack

awslocal apigateway put-integration \
    --rest-api-id $API_ID \
    --resource-id $API_RESOURCE_ID \
    --http-method ANY \
    --type AWS_PROXY \
    --integration-http-method POST \
    --uri arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/$FUNCTION_ARN/invocations \
    --profile localstack

awslocal lambda add-permission \
    --function-name GetMovies \
    --statement-id gguid \
    --action Lambda:InvokeFunction \
    --principal apigateway.amazonaws.com \
    --source-arn arn:aws:execute-api:us-east-1:$ACCOUNT_ID:$API_ID \
    --profile localstack

awslocal apigateway test-invoke-method \
    --rest-api-id $API_ID \
    --resource-id $API_RESOURCE_ID \
    --http-method "GET" \
    --path-with-query-string "/$PATH_PART" \
    --profile localstack


API_DEPLOYMENT_ID=$(awslocal apigateway create-deployment \
    --rest-api-id $API_ID \
    --stage-name staging \
    --profile localstack | \
    jq -r '.id')

curl -X GET http://$API_ID.execute-api.localhost.localstack.cloud:4566/$API_DEPLOYMENT_ID/movies
