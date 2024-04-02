package resp

const (
	//公共响应码
	RESP_OK         = 10000
	RESP_ERR        = 10001
	RESP_PARAM_ERR  = 10002
	RESP_TOKEN_ERR  = 10003
	RESP_NO_ACCESS  = 10004
	RESP_APP_NOT_ON = 10005
	RESP_EXPORT_ERR = 10006
	RESP_Insert_ERR = 10007
)

type MResp struct {
	F_responseNo  int    `required:"true" description:"响应码"`
	F_responseMsg string `description:"响应码描述"`
}

// 获取response no config
func GetResponseConfig() map[int]string {
	rep := map[int]string{}
	rep[RESP_OK] = "成功"
	rep[RESP_ERR] = "失败,未知错误"
	rep[RESP_PARAM_ERR] = "参数错误"
	rep[RESP_TOKEN_ERR] = "token错误"
	rep[RESP_NO_ACCESS] = "没有访问权限"
	rep[RESP_APP_NOT_ON] = "暂时未提供服务"
	rep[RESP_EXPORT_ERR] = "导出数据失败"
	rep[RESP_Insert_ERR] = "导入数据失败"
	return rep
}
