package parse

import (
	"errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

type ArgType int
type Location int

const (
	Json ArgType = iota
	Array
	Bool
	Int
	Str
	Header Location = iota
	Post
	Get
)

type RequestParser struct {
	args []requestParserArgs
}

type requestParserArgs struct {
	name     string
	argType  ArgType
	required bool
	help     string
	location Location
}

func (r *RequestParser) AddArgument(name string, required bool, help string, argType ArgType, location Location) {
	r.args = append(r.args, requestParserArgs{name: name, argType: argType, required: required, help: help, location: location})
}

func NewRequestParser() *RequestParser {
	return &RequestParser{}
}



func (r *RequestParser) Parse(res *http.Request) (args map[string]interface{}, err error) {
	args = make(map[string]interface{})
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = res.ParseForm()
	if err != nil {
		return nil, err
	}
	for _, v := range r.args {
		if v.location == Get {
			value := res.Form.Get(v.name)
			if value == "" && v.required == true {
				return args, errors.New(v.help)
			}
			args[v.name] = value
		} else if v.location == Post {
			value := gjson.GetBytes(bodyBytes, v.name)
			args[v.name] = value.String()
			if value.String() == "" && v.required == true {
				return args, errors.New(v.help)
			}
			if v.argType == Json && !value.IsObject() {
				return args, errors.New(v.name + " type error,real type is " + value.Type.String())
			}
			if v.argType == Array && !value.IsArray() {
				return args, errors.New(v.name + " type error,real type is " + value.Type.String())
			}
			if v.argType == Bool && !(value.Type == gjson.True || value.Type == gjson.False) {
				return args, errors.New(v.name + " type error,real type is " + value.Type.String())
			}
			if v.argType == Int && !(value.Type == gjson.Number) {
				return args, errors.New(v.name + " type error,real type is " + value.Type.String())
			}
		} else if v.location == Header {
			value := res.Header.Get(v.name)
			if value == "" && v.required == true {
				return args, errors.New(v.help)
			}
			args[v.name] = value
		}
	}
	return args, nil
}
