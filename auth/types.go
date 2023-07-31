package auth

type Session struct {
	ID     string `json:"id"`
	UserID int    `json:"user_id"`
	Active bool   `json:"active"`
	Token  string `json:"token"`
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
