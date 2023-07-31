package auth

type Session struct {
	ID     string
	UserID int
	Active bool
	Token  string
}

type SignupReq struct {
	Login    string
	Password string
	Name     string
	Age      string
}

type SignupReqOpts struct {
	Login    string
	Password string
	Name     string
	Age      string
}

type LoginReq struct {
	Login    string
	Password string
}
