package model

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *User) Validate() bool {
	if u.Name == "" || len(u.Name) > 50 {
		return false
	}
	if u.Age < 20 || u.Age > 80 {
		return false
	}
	return true
}
