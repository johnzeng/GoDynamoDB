package GoDynamoDB

import "github.com/aws/aws-sdk-go/service/dynamodb"
import "github.com/aws/aws-sdk-go/aws"

type QueryExecutor struct {
	input     *dynamodb.QueryInput
	db        *dynamodb.DynamoDB
	prototype *ReadModel
	ret       []ReadModel
	count     int64
}

func (db GoDynamoDB) GetQueryExecutor(i ReadModel) (*QueryExecutor, error) {
	ret := &QueryExecutor{
		db:        db.db,
		prototype: &i,
		count:     0,
	}
	input := &dynamodb.QueryInput{
		TableName: aws.String(i.GetTableName()),
	}

	ret.input = input
	return ret, nil
}

func (q *QueryExecutor) AddValue(express string, v interface{}) *QueryExecutor {
	if nil == q.input.ExpressionAttributeValues {
		q.input.ExpressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
	}
	att, err := encodeToQueryAtt(v)
	if err != nil {
		return nil
	}
	q.input.ExpressionAttributeValues[express] = att
	return q
}

func (q *QueryExecutor) WithKeyCondition(helper *QueryCondExpressHelper) *QueryExecutor {
	if nil == q.input.ExpressionAttributeNames {
		q.input.ExpressionAttributeNames = make(map[string]*string)
	}
	for key, value := range helper.expressMap {
		q.input.ExpressionAttributeNames[key] = value
	}
	q.input.KeyConditionExpression = aws.String(helper.str)
	return q
}

func (q *QueryExecutor) Exec() error {
	resp, err := q.db.Query(q.input)
	if nil != err {
		return err
	}

	if nil == q.ret {
		q.ret = make([]ReadModel, 0)
	}

	for i := 0; i < len(resp.Items); i++ {
		rspi := resp.Items[i]
		orgPrototype := *q.prototype
		err := decode(rspi, &orgPrototype)
		if err != nil {
			return err
		}
		q.ret = append(q.ret, orgPrototype)
	}
	q.count = *resp.Count
	return nil
}

func (q *QueryExecutor) GetRet() []ReadModel {
	return q.ret
}

func (q *QueryExecutor) GetCount() int64 {
	return q.count
}
