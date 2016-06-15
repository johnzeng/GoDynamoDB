package GoDynamoDB

import "github.com/aws/aws-sdk-go/service/dynamodb"
import "github.com/aws/aws-sdk-go/aws"
import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func encodeToAtt(v reflect.Value) (*dynamodb.AttributeValue, error) {
	switch v.Kind() {
	case reflect.Bool:
		b := v.Bool()
		return &dynamodb.AttributeValue{BOOL: aws.Bool(b)}, nil
	case reflect.String:
		s := v.String()
		if len(s) == 0 {
			b := true
			return &dynamodb.AttributeValue{NULL: &b}, nil
		} else {
			return &dynamodb.AttributeValue{S: aws.String(s)}, nil
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := strconv.FormatInt(v.Int(), 10)
		return &dynamodb.AttributeValue{N: &n}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n := strconv.FormatUint(v.Uint(), 10)
		return &dynamodb.AttributeValue{N: &n}, nil

	case reflect.Float32, reflect.Float64:
		f := v.Float()
		if math.IsInf(f, 0) || math.IsNaN(f) {
			panic(fmt.Errorf("aws.dynamodb.convertToNumericString: NaN and infinite floats not supported"))
		}
		fs := strconv.FormatFloat(f, 'g', -1, v.Type().Bits())
		return &dynamodb.AttributeValue{N: &fs}, nil

	case reflect.Struct:
		if v.IsNil() {
			b := true
			return &dynamodb.AttributeValue{NULL: &b}, nil
		} else {
			m, err := encodeStruct(v)
			if err == nil {
				return &dynamodb.AttributeValue{M: m}, nil
			} else {
				return nil, err
			}
		}
	default:
		return nil, nil
	}
	return nil, nil
}

func encodeStruct(v reflect.Value) (map[string]*dynamodb.AttributeValue, error) {
	out := map[string]*dynamodb.AttributeValue{}
	val := v.Elem()
	for i := 0; i < val.NumField(); i++ {
		typeInfo := val.Type()
		fileName := typeInfo.Field(i).Name
		fmt.Printf("field name: %s\n", fileName)
		att, err := encodeToAtt(val.Field(i))
		if err == nil {
			out[fileName] = att
		} else {
			return nil, err
		}
	}
	return out, nil
}

func encode(i interface{}) (map[string]*dynamodb.AttributeValue, error) {
	v := reflect.ValueOf(i)
	return encodeStruct(v)
}
