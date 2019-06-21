package parse

import (
	"errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
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
	var bodyBytes []byte
	args = make(map[string]interface{})
	if res.Body != nil{
		bodyBytes, err = ioutil.ReadAll(res.Body)
	}
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
			switch v.argType {
			case Int:
				num,err := strconv.Atoi(value)
				args[v.name] = num
				if err != nil{
					args[v.name] = nil
				}
			case Bool:
				parseboolstr, err := strconv.ParseBool(value)
				args[v.name] = parseboolstr
				if err != nil{
					args[v.name] = nil
				}
			case Str:
				if value == ""{
					args[v.name] = nil
				}else {
					args[v.name] = value
				}
			}
		} else if v.location == Post {
			value := gjson.GetBytes(bodyBytes, v.name)
			switch v.argType {
			case Bool:
				args[v.name] = value.Bool()
			case Str:
				args[v.name] = value.String()
			case Json:
				args[v.name] = value.String()
			case Array:
				args[v.name] = value.String()
			case Int:
				args[v.name] = value.Int()
			}
			if value.String() == "" && v.required == true {
				return args, errors.New(v.help)
			}
			if v.required == true{
				if v.argType == Json && !value.IsObject() {
					return args, errors.New(v.name + " type error,real type is " + value.Type.String())
				}
				if v.argType == Array && !value.IsArray() {
					return args, errors.New(v.name + " type error,real type is " + value.Type.String())
				}
				if v.argType == Bool && !(value.Type == gjson.True || value.Type == gjson.False)  {
					return args, errors.New(v.name + " type error,real type is " + value.Type.String())
				}
				if v.argType == Int && !(value.Type == gjson.Number)  {
					return args, errors.New(v.name + " type error,real type is " + value.Type.String())
				}
			} else {
				if v.argType == Json && !(value.IsObject() || (value.Type == gjson.Null)) {
					return args, errors.New(v.name + " type error,real type is " + value.Type.String())
				}
				if v.argType == Array && !(value.IsArray() || (value.Type == gjson.Null)){
					return args, errors.New(v.name + " type error,real type is " + value.Type.String())
				}
				if v.argType == Bool && !((value.Type == gjson.True || value.Type == gjson.False) || (value.Type == gjson.Null)) {
					return args, errors.New(v.name + " type error,real type is " + value.Type.String())
				}
				if v.argType == Int && !((value.Type == gjson.Number) || (value.Type == gjson.Null)){
					return args, errors.New(v.name + " type error,real type is " + value.Type.String())
				}
			}

		} else if v.location == Header {
			value := res.Header.Get(v.name)
			if value == "" && v.required == true {
				return args, errors.New(v.help)
			}
			switch v.argType {
			case Int:
				num,err := strconv.Atoi(value)
				args[v.name] = num
				if err != nil{
					args[v.name] = nil
				}
			case Bool:
				parseboolstr, err := strconv.ParseBool(value)
				args[v.name] = parseboolstr
				if err != nil{
					args[v.name] = nil
				}
			case Str:
				args[v.name] = value
			}
		}
	}
	return args, nil
}
