package utils

const (
	CODE_SUCCESS = 200
	CODE_ERR_APP = iota + 1000
	CODE_ERR_MSG
	CODE_ERR_PARAM
	CODE_ERR_TOKEN
	CODE_ERR_NO_PRIV
)

var MsgFlags = map[int]string{
	CODE_SUCCESS:     "成功",
	CODE_ERR_MSG:     "未知错误",
	CODE_ERR_PARAM:   "参数错误",
	CODE_ERR_TOKEN:   "没有TOKEN",
	CODE_ERR_NO_PRIV: "沒有权限",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[CODE_ERR_APP]
}
