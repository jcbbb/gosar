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
