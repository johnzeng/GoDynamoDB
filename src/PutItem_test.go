package GoDynamoDB

import (
	"testing"
)

func Test_PutItem(t *testing.T) {
	err := GetDBInstance().PutItem(&TestStruct{Name: "John", Id: "123"})
	if err != nil {
		t.Error(err.Error())
	}
}

func init() {
	InitLocalDBInstance("http://127.0.0.1:8000")
}
