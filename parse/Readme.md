## 基于gjson的http args 解析方法



快速使用

```
var P *parse.RequestParser

func init(){
	P = parse.NewRequestParser()
	P.AddArgument("did",true,"did error",parse.Json,parse.Post)
	P.AddArgument("ssu-uid",true,"ssu-uid error",parse.Str,parse.Header)
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
```

增加参数方法

```
func (r *RequestParser) AddArgument(name string, required bool, help string, argType ArgType, location Location)()

name : 要查找的字段
required : 是否必须有该字段
help : 解析错误时返回内容
argType : 参数的类别：Json Array Bool Int Str 五种
location : 参数位置 Header  Get  Post(put patch)
```

解析方法

```
func (r *RequestParser) Parse(res *http.Request) (args map[string]interface{}, err error) {}
对上面的参数进行解析，并返回参数map 类型map[string]interface{}
```

