package user

type User struct {
	ID       int    `json:"id"`
	Password string `json:"-"`
	Login    string `json:"login"`
	Age      int    `json:"age"`
	Name     string `json:"name"`
}

type Opts struct {
	Login    string
	Age      int
	Name     string
	Password string
}

type Phone struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
	IsFax       bool   `json:"is_fax"`
}

type AddPhoneReq struct {
	phone       string
	description string
	isFax       bool
}

type AddPhoneOpts struct {
	phone       string
	description string
	isFax       bool
	userId      int
}
