package modles

type Auth struct {
	Id       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
