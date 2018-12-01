package rpc

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
)

type Response struct {
	Params []struct {
		Value Value `xml:"value"`
	} `xml:"params>param"`
	Fault struct {
		Details []Member `xml:"struct>member"`
	} `xml:"fault>value"`
}

type Value struct {
	Array   []Value  `xml:"array>data>value"`
	Struct  []Member `xml:"struct>member"`
	String  string   `xml:"string"`
	Boolean bool     `xml:"boolean"`
	Int     int      `xml:"int"`
	Int4    int32    `xml:"i4"`
	Int8    int64    `xml:"i8"`
	Inner   string   `xml:",innerxml"`
}

type Member struct {
	Name  string `xml:"name"`
	Value Value  `xml:"value"`
}

type Fault struct {
	Code   int32
	String string
}

func (f Fault) Error() string {
	return fmt.Sprintf("%d: %s", f.Code, f.String)
}

func decode(r []byte) (Response, error) {
	var res Response
	err := xml.NewDecoder(bytes.NewReader(r)).Decode(&res)
	if err != nil {
		return res, err
	}
	if len(res.Fault.Details) > 0 {
		return res, getFault(&res)
	}
	return res, nil
}

func DecodeResponseForFault(r []byte) error {
	_, err := decode(r)
	return err
}

func DecodeResponse(r []byte, v interface{}) error {
	res, err := decode(r)
	if err != nil {
		return err
	}
	return convert(res.Params[0].Value.Array, v)
}

func convert(rpcResponse []Value, v interface{}) error {
	rpcResponseLen := len(rpcResponse)
	result := reflect.MakeSlice(reflect.TypeOf(reflect.ValueOf(v).Elem().Interface()), rpcResponseLen, rpcResponseLen)

	if rpcResponseLen > 0 && len(rpcResponse[0].Array) > result.Index(0).NumField() {
		return errors.New("Not enough fields to accommodate the response")
	}

	for i, responseValue := range rpcResponse {
		item := result.Index(i)
		for j, field := range responseValue.Array {
			var val interface{}
			switch item.Field(j).Interface().(type) {
			case string:
				val = field.String
			case bool:
				val = field.Boolean
			case int:
				val = field.Int
			case int32:
				val = field.Int4
			case int64:
				val = field.Int8
			}

			if val != nil {
				item.Field(j).Set(reflect.ValueOf(val))
			}
		}
	}

	reflect.ValueOf(v).Elem().Set(result)
	return nil
}

func getFault(res *Response) Fault {
	var fault Fault

	for _, member := range res.Fault.Details {
		switch member.Name {
		case "faultCode":
			fault.Code = member.Value.Int4
		case "faultString":
			fault.String = member.Value.String
		}
	}

	return fault
}
