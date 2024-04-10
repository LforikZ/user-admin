package controllers

import (
	"context"
	"encoding/json"
	"food-service/global/resp"
	"food-service/infrastructure/common"
	"github.com/astaxie/beego/validation"
	"net/http"

	"github.com/astaxie/beego"
)

// 公共controller
type BaseController struct {
	beego.Controller
	AppCtx context.Context
}

// 操作成功
func (c *BaseController) jsonEchoOk() {
	c.jsonEcho(map[string]interface{}{
		"F_responseNo": common.RESP_OK,
	})
}

// json echo
func (c *BaseController) jsonEcho(datas map[string]interface{}) {
	responseMsg, ok := datas[common.F_RES_MES]
	if !ok || (ok && len(responseMsg.(string)) <= 0) {
		datas[common.F_RES_MES] = ""
		msg, ok := common.MyConfig.ConfigMyResponse[datas[common.F_RES_NO].(int)]
		if ok {
			datas[common.F_RES_MES] = msg
		}
	}
	if datas[common.F_RES_NO].(int) == common.RESP_PARAM_ERR { //参数错误
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
	}
	if datas[common.F_RES_NO].(int) == common.RESP_TOKEN_ERR { //token(access token , refresh access token) 错误
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
	}
	//
	c.Data["json"] = datas
	//跨域支持
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	//输出
	c.ServeJSON()
}

// 验证表单参数
func (c *BaseController) paramsValid(obj interface{}) bool {
	if err := c.ParseForm(obj); err == nil {
		valid := validation.Validation{}
		b, err := valid.Valid(obj)
		if err == nil && b {
			return true
		}
	}
	return false
}

// 验证Json参数
func (c *BaseController) jsonValid(jsonData []byte, obj interface{}) bool {
	if err := json.Unmarshal(jsonData, obj); err == nil {
		valid := validation.Validation{}
		b, err := valid.Valid(obj)
		if err == nil && b {
			return true
		}
	}
	return false
}

type HttpResponse1 struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Code       int         `json:"code"`
	Data       interface{} `json:"data"`
}

func (c *BaseController) renderJSON(code int, msg string, data interface{}) {
	c.Ctx.Output.Header("Content-Type", "application/json")
	switch code {
	case resp.RESP_PARAM_ERR:
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
	case resp.RESP_TOKEN_ERR:
		c.Ctx.Output.SetStatus(http.StatusForbidden)
	case resp.RESP_NO_ACCESS:
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
	case resp.RESP_ERR:
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
	default:
		c.Ctx.Output.SetStatus(http.StatusOK)
	}

	res := HttpResponse1{
		StatusCode: c.Ctx.Output.Status,
		Code:       code,
		Message:    msg,
	}
	if data != nil {
		res.Data = data
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// 成功渲染函数
func (c *BaseController) renderSuccessJSON(data interface{}) {
	c.renderJSON(resp.RESP_OK, "ok", data)
}

// 成功渲染函数
func (c *BaseController) renderSuccessJSON1(data interface{}, msg string) {
	c.renderJSON(resp.RESP_OK, msg, data)
}

// 参数错误
func (c *BaseController) renderParamsError(data interface{}) {
	c.renderJSON(resp.RESP_PARAM_ERR, "参数错误", data)
}

// 参数错误 (额外描述)
func (c *BaseController) renderParamsError1(msg string, data interface{}) {
	c.renderJSON(resp.RESP_PARAM_ERR, msg, data)
}

// 未知错误
func (c *BaseController) renderUnknownError(msg string, data interface{}) {
	c.renderJSON(resp.RESP_ERR, msg, data)
}
