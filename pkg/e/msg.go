package e

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	ERROR:                          "fail",
	INVALID_PARAMS:                 "请求参数错误",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",

	ERROR_USER_NOT_EXIST: "用户不存在",
	ERROR_USER_EXIST:     "用户名已存在",
	ERROR_USER_WRONG_PWD: "密码错误",
	ERROR_INVALID_CAPTCHA: "验证码错误",
	ERROR_SEND_EMAIL: "发送邮件失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}