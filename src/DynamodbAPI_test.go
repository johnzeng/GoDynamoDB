package GoDynamoDB

import "testing"
import "github.com/aws/aws-sdk-go/service/dynamodb"
import "github.com/aws/aws-sdk-go/aws/session"
import "github.com/aws/aws-sdk-go/aws"

func TestDynamoScanAPI(t *testing.T) {

	db := dynamodb.New(session.New(
		&aws.Config{
			Region:   aws.String("us-east-1"),
			Endpoint: aws.String("http://127.0.0.1:8000")}))

	params := &dynamodb.ScanInput{
		TableName: aws.String("Test"), // Required
		AttributesToGet: []*string{
			aws.String("id"), // Required
		},
		Limit: aws.Int64(100),
	}
	resp, err := db.Scan(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		t.Error(resp.String())
	}
}

func TestDynamoPutAPI(t *testing.T) {

	db := dynamodb.New(session.New(
		&aws.Config{
			Region:   aws.String("us-east-1"),
			Endpoint: aws.String("http://127.0.0.1:8000")}))

	params := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{ // Required
			"id": { // Required
				S: aws.String("thisisId"),
			},
			"tet": {
				S: aws.String("value"),
			},
			// More values...
		},
		TableName:    aws.String("Test"), // Required
		ReturnValues: aws.String("NONE"),
	}
	resp, err := db.PutItem(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		t.Error(resp.String())
	}

}
