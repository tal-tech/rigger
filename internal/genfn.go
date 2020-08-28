package internal

import (
	"bytes"
	"fmt"
	"strings"
)

func GenXesSDKFunc(tplType, apiName, funcName, path string) (*bytes.Buffer, error) {
	apiName = strings.Title(apiName)
	funcName = strings.Title(funcName)
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf(handleTplType(tplType), apiName, funcName, funcName, funcName, apiName, funcName, apiName, funcName, path, apiName, funcName, apiName, funcName, apiName, funcName))

	return buffer, nil
}

func handleTplType(tplType string) (tpl string) {
	switch tplType {
	case "zhongtai":
		tpl = `func (this *%s) %s(ctx context.Context, req %sRequest) (resp %sResponse, err error) {
	perfutil.CountI("%s.%s")
	defer perfutil.AutoElapsed("%s.%s", time.Now())
	path := "%s"
	reqHeader := header.GenHeader(ctx, header.GenGatewayAuthHeader, header.GenTraceHeader)
	ret, err := request.PostParam(ctx, HOST, path, req, reqHeader, this.client)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,PostParam err:%%v,req:%%+v", "%s", err, req)
		return
	}
	data, err := valid.CheckValid(ctx, ret, valid.CheckStat1)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,CheckValid err:%%v,req:%%+v,ret:%%s", "%s", err, req, ret)
		return
	}
	err = jsutil.Json.Unmarshal(data, &resp)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,json Umarshal err:%%v,req:%%+v,ret:%%s", "%s", err, req, data)
		return
	}
	return
}`
	case "irc":
		tpl = `func (this *%s) %s(ctx context.Context, urlParam url.Values, req %sRequest) (resp %sResponse, err error) {
	perfutil.CountI("%s.%s")
	defer perfutil.AutoElapsed("%s.%s", time.Now())
	path := "%s" + "?" + urlParam.Encode()
	reqHeader := header.GenHeader(ctx, header.GenIRCAuthHeader, header.GenTraceHeader)
	ret, err := request.PostJson(ctx, HOST, path, req, reqHeader, this.client)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,PostJson err:%%v,req:%%+v", "%s", err, req)
		return
	}
	data, err := valid.CheckValid(ctx, ret, valid.CheckCode0)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,CheckValid err:%%v,req:%%+v,ret:%%s", "%s", err, req, ret)
		return
	}
	err = jsutil.Json.Unmarshal(data, &resp)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,json Umarshal err:%%v,req:%%+v,ret:%%s", "%s", err, req, data)
		return
	}
	return
}`
	case "oa":
		tpl = `func (this *%s) %s(ctx context.Context, req %sRequest) (resp %sResponse, err error) {
	perfutil.CountI("%s.%s")
	defer perfutil.AutoElapsed("%s.%s", time.Now())
	path := "%s"
	reqHeader := header.GenHeader(ctx, header.GenGatewayAuthHeader, header.GenTraceHeader)
	ret, err := request.PostParam(ctx, HOST, path, req, reqHeader, this.client)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,PostParam err:%%v,req:%%+v", "%s", err, req)
		return
	}
	data, err := valid.CheckValid(ctx, ret, valid.CheckStat1Msg)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,CheckValid err:%%v,req:%%+v,ret:%%s", "%s", err, req, ret)
		return
	}
	err = jsutil.Json.Unmarshal(data, &resp)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,json Umarshal err:%%v,req:%%+v,ret:%%s", "%s", err, req, data)
		return
	}
	return
}`
	case "tiku":
		tpl = `func (this *%s) %s(ctx context.Context, req %sRequest) (resp %sResponse, err error) {
	perfutil.CountI("%s.%s")
	defer perfutil.AutoElapsed("%s.%s", time.Now())
	path := "%s"
	reqHeader := header.GenHeader(ctx, header.GenGatewayAuthHeader, header.GenTraceHeader)
	ret, err := request.PostParam(ctx, HOST, path, req, reqHeader, this.client)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,PostParam err:%%v,req:%%+v", "%s", err, req)
		return
	}
	data, err := valid.CheckValid(ctx, ret, valid.CheckStatus100)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,CheckValid err:%%v,req:%%+v,ret:%%s", "%s", err, req, ret)
		return
	}
	err = jsutil.Json.Unmarshal(data, &resp)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,json Umarshal err:%%v,req:%%+v,ret:%%s", "%s", err, req, data)
		return
	}
	return
}`
	case "media":
		tpl = `func (this *%s) %s(ctx context.Context, req %sRequest) (resp %sResponse, err error) {
	perfutil.CountI("%s.%s")
	defer perfutil.AutoElapsed("%s.%s", time.Now())
	path := "%s"
	reqHeader := header.GenHeader(ctx, header.GenGatewayAuthHeader, header.GenTraceHeader)
	ret, err := request.GetParam(ctx, HOST, path, req, reqHeader, this.client)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,PostParam err:%%v,req:%%+v", "%s", err, req)
		return
	}
	data, err := valid.CheckValid(ctx, ret, valid.CheckStat0Msg)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,CheckValid err:%%v,req:%%+v,ret:%%s", "%s", err, req, ret)
		return
	}
	err = jsutil.Json.Unmarshal(data, &resp)
	if err != nil {
		logger.Ex(ctx, "%s", "methods:%%s,json Unmarshal err:%%v,req:%%+v,ret:%%s", "%s", err, req, data)
		return
	}
	return
}`
	}
	return
}
