package echoapp

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	Password string `json:"password2"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

func (r *User) TableName() string {
	return "users"
}
