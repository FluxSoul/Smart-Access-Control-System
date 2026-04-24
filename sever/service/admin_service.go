package service

import (
	"EmqxBackEnd/models"
	"EmqxBackEnd/repository"

	"github.com/golang-jwt/jwt/v5"

	"time"
)

var jwtSecret = []byte("cqupt") // Should be stored securely

func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func CreateAdmin(username, password string) (time.Time, error) {
	createdTime, err := repository.CreateAdmin(username, password)
	if err != nil {
		return time.Time{}, err
	}
	return createdTime, nil
}

func CheckLogin(username, password string) (int, bool) {
	var user *models.EmpxAdmin
	user, err := repository.GetAdminByUser(username)
	if err != nil {
		return -1, false
	}
	if user == nil {
		return -1, false
	}
	if user.Status != 1 {
		return -1, false
	}
	if user.Password != password {
		return -1, false
	}
	id := user.ID
	expiresAt := time.Now().Add(time.Hour * 24)
	err = repository.UpdateExpiresAtTime(expiresAt, id)
	if err != nil {
		return -1, false
	}
	return user.ID, true
}

func SaveToken(token string, id int) error {
	return repository.SaveToken(token, id)
}

func ChangeUserStatus(id int, status int8) error {
	return repository.ChangeUserStatus(id, status)
}

func GetAllUsers() ([]models.GetEmpxAdmin, error) {
	return repository.GetAllUsers()
}

func IsAdmin(token string) bool {
	var user *models.EmpxAdmin
	user, err := repository.GetAdminByUser("admin")
	if err != nil {
		return false
	}
	if user == nil {
		return false
	}
	// 2为超级管理员
	tokenAdmin, err := repository.GetToken(user.ID)
	if err != nil {
		return false
	}
	return token == tokenAdmin
}

func GetUserIdByToken(token string) (int, error) {
	id, err := repository.GetUserIdByToken(token)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func GetAllNodeByUserId(id int) ([]models.Node, error) {
	return repository.GetAllNodeByUserId(id)
}
