package controller


type ResCode int64

// 自定义业务的状态码，前后端协商
const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeInvalidToken
	CodeNeedLogin
)

var codeMsgMap = map[ResCode]string {
	CodeSuccess: "success",
	CodeInvalidParam: "请求参数错误",
	CodeUserExist:	"用户已存在",
	CodeUserNotExist: 	"用户不存在",
	CodeInvalidPassword: "密码错误",
	CodeServerBusy:		"服务繁忙",

	CodeInvalidToken:  "无效的Token",
	CodeNeedLogin:	"需要登录",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}