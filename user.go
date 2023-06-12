package cloud_storage

type User struct {
	Id       int    `json:"-" db:"id"`
	Login    string `json:"login" building:"required"`
	Name     string `json:"name" building:"required"`
	Username string `json:"username" building:"required"`
	Password string `json:"password" building:"required"`
}

type Node struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	ServerName string  `json:"server_name"`
	Size       int     `json:"size"`
	CreateDate string  `json:"create_date"`
	Type       string  `json:"type"`
	Children   []*Node `json:"children"`
}

type UserInfo struct {
	Id              int    `json:"id" db:"id"`
	Login           string `json:"login" building:"required"`
	Name            string `json:"name" building:"required"`
	Username        string `json:"username" building:"required"`
	AvailableMemory int    `json:"available_memory" building:"required"`
}
