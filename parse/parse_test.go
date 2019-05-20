package parse

import (
	"encoding/json"
	"fmt"
	"github.com/luslin/parse"
	"net/http"
	"testing"
)

var P *parse.RequestParser

func init(){
	P = parse.NewRequestParser()
	P.AddArgument("did",true,"did error",parse.Json,parse.Post)
	P.AddArgument("ssu-uid",true,"ssu-uid error",parse.Str,parse.Header)
}

func TestParse(t *testing.T) {
	http.HandleFunc("/", IndexHandler)
	_ = http.ListenAndServe("127.0.0.1:8001", nil)
}


func IndexHandler(w http.ResponseWriter, r *http.Request) {
	res , err := P.Parse(r)
	if err != nil{
		_ ,_ = fmt.Fprintln(w, err.Error())
		return
	}
	ret ,_ := json.Marshal(res)
	fmt.Println(ret)
	_ ,_ = fmt.Fprintln(w, string(ret))
}


