package request

type Admin struct {
	UserName string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}
