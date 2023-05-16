package cloud_storage

type User struct {
	Id       int    `json:"-" db:"id"`
	Login    string `json:"login" building:"required"`
	Name     string `json:"name" building:"required"`
	Username string `json:"username" building:"required"`
	Password string `json:"password" building:"required"`
}
