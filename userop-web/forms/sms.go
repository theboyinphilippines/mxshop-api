package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //自定义validator 验证手机号
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`
}
