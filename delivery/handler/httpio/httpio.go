package httpio

import (
	"central-dummy-master/entities"
	"central-dummy-master/entities/params"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type RequestIO interface {
	Recv()
	Bind(obj interface{})
	BindJSON(obj interface{})
	Response(statusCode int, responseBody interface{})
	ResponseWithAbort(statuscode int, responseBody interface{})
	ResponseString(statusCode int, response string)
	ResponseStringWithAbort(statusCode int, response string)
}

func NewRequestIO(ctx *gin.Context) RequestIO {
	_ = ctx.Request.ParseForm()
	return &formio{
		context: ctx,
		request: ctx.Request,
	}
}

type formio struct {
	context *gin.Context
	request *http.Request
}

// No Content Receiver
func (f *formio) Recv() {
	header := params.Header{}
	path := fmt.Sprintf("%s %s", f.request.Method, f.request.URL.Path)
	_ = f.context.ShouldBindHeader(&header)
	go receiveForm("RECV", header, "", f.request.RemoteAddr, path)
}

// These methods use MustBindWith under the hood. If there is a binding error, the request is aborted with c.AbortWithError(400, err).SetType(ErrorTypeBind).
// This sets the response status code to 400 and the Content-Type header is set to text/plain; charset=utf-8
func (f *formio) Bind(body interface{}) {
	header := params.Header{}
	path := fmt.Sprintf("%s %s", f.request.Method, f.request.URL.Path)
	_ = f.context.ShouldBindHeader(&header)
	_ = f.context.Bind(body)
	go receiveForm("RECV", header, body, f.request.RemoteAddr, path)
}

//These methods use MustBindWith under the hood. If there is a binding error, the request is aborted with c.AbortWithError(400, err).SetType(ErrorTypeBind).
//This sets the response status code to 400 and the Content-Type header is set to text/plain; charset=utf-8

func (f *formio) BindJSON(body interface{}) {
	path := fmt.Sprintf("%s %s", f.request.Method, f.request.URL.Path)
	header := params.Header{}
	_ = f.context.ShouldBindHeader(&header)
	_ = f.context.BindJSON(body)
	go receiveJSON("RECV", header, body, f.request.RemoteAddr, path)

}

func (f *formio) Response(StatusCode int, responseBody interface{}) {
	f.sendResponse(StatusCode, responseBody, false)
}

func (f *formio) ResponseWithAbort(statuscode int, responseBody interface{}) {
	f.sendResponse(statuscode, responseBody, true)
}

func (f *formio) ResponseString(statusCode int, response string) {
	f.sendResponseString(statusCode, response, false)
}

func (f *formio) ResponseStringWithAbort(statusCode int, response string) {
	f.sendResponseString(statusCode, response, true)
}

func (f *formio) sendResponseString(statusode int, responseBody string, abort bool) {
	out("SEND", responseBody, f.request.RemoteAddr, f.request.URL.Path)
	f.context.String(statusode, responseBody)
	if abort {
		f.context.Abort()
	}
}

func (f *formio) sendResponse(statusode int, responseBody interface{}, abort bool) {
	outJson("SEND", responseBody, f.request.RemoteAddr, f.request.URL.Path)
	f.context.Header("Content-Type", "application/json")
	f.context.JSON(statusode, responseBody)
	if abort {
		f.context.Abort()
	}
}

func outJson(cmd string, values interface{}, addr string, path string) {
	var rsBuff []string
	title := ""
	if cmd == "RECV" {
		title = "\nPAYLOAD :\n"
	} else {
		title = "\nRESPONSE : \n"
	}
	rsBuff = append(rsBuff, title)
	j, _ := json.MarshalIndent(values, "", "    ")
	rsBuff = append(rsBuff, string(j))
	out(cmd, strings.Join(rsBuff, ""), addr, path)
}

func receiveJSON(cmd string, header interface{}, body interface{}, addr string, path string) {
	var rsBuff []string
	rsBuff = append(rsBuff, parseContent(header, "header"))
	j, _ := json.MarshalIndent(body, "", "    ")
	rsBuff = append(rsBuff, string(j))
	out(cmd, strings.Join(rsBuff, ""), addr, path)
}

func receiveForm(cmd string, header interface{}, body interface{}, addr string, path string) {
	msgHeader := parseContent(header, "header")
	msgBody := parseContent(body, "form")
	msg := fmt.Sprint(msgHeader, msgBody)
	out(cmd, msg, addr, path)
}

func parseContent(values interface{}, tag string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in parseContent", r)
		}
	}()

	v := reflect.ValueOf(values)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return "empty body"
	}

	valType := v.Type()
	var tagItem string
	msg := fmt.Sprintf("\n%s :\n", strings.ToUpper(tag))
	for i := 0; i < v.NumField(); i++ {
		field := valType.Field(i)
		tagItem = field.Tag.Get(tag)
		if tagItem != "" {
			value := fmt.Sprint(v.Field(i).Interface())
			if value != "" {
				msg = fmt.Sprintf("%s[%v]=%v\n", msg, tagItem, value)
			}
		}
	}
	return msg
}

func out(cmd string, values string, addr string, path string) {
	entities.PrintLog(fmt.Sprintf("%s, %s, %s, %s", cmd, addr, path, values))
}
