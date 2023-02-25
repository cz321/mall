package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Type uint32 `form:"type" json:"type" binding:"required,oneof=1 2"`	//1代表注册，2代表验证码登录
}

