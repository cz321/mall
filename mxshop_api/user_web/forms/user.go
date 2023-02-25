package forms

type PasswordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha   string `from:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `from:"captcha_id" json:"captcha_id" binding:"required"`
}

type RegisterFrom struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Code  string `form:"code" json:"code" binding:"required,min=4,max=4"`
}
