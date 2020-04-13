package service

import (
	"feelingliu/modles"
	"feelingliu/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte(modles.AppInfo.JwtSecret)


type CustomClaims struct {
	User
	jwt.StandardClaims
}

type User struct {
	ID           int    `json:"id" db:"id" form:"id"`
	Username     string `json:"username" db:"username" form:"username"`
	Password     string `json:"password" db:"password" form:"password"`
	Introduction string `json:"introduction" db:"introduction" form:"introduction"`
	Avatar       string `json:"avatar" db:"avatar" form:"avatar"`
	Nickname     string `json:"nickname" db:"nickname" form:"nickname"`
	About        string `json:"about" db:"about" form:"about"`
}


func (u *User) CheckAuth() bool {
	var user User
	db := modles.DB.Where("username = ? AND password = ?", u.Username, u.Password).Find(&user)
	if db.Error != nil {
		utils.WriteErrorLog(db.Error)
	}
	return user.ID > 0
}

func (u User) GenToken() (string, error) {
	claims := CustomClaims{u, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * time.Duration(modles.AppInfo.TokenTimeout)).Unix(),
		Id:        fmt.Sprintf("%v", time.Now().UnixNano()),
	}}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, e := tokenClaims.SignedString(jwtSecret)
	return token, e
}
