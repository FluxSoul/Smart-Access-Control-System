package repository

import (
	"EmqxBackEnd/database"
	"EmqxBackEnd/models"
	"database/sql"
	"errors"
	"log"
	"time"
)

// GetAdminByUser 获取用户
func GetAdminByUser(username string) (*models.EmpxAdmin, error) {
	query := `SELECT id, username, password, status, created_time FROM public.admin WHERE username = $1`
	var admin models.EmpxAdmin
	err := database.DB.QueryRow(query, username).Scan(
		&admin.ID,
		&admin.Username,
		&admin.Password,
		&admin.Status,
		&admin.CreatedTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No admin found with username:", username)
			return nil, nil
		}
		log.Println("Error querying admin:", err)
		return nil, err
	}
	return &admin, nil

}

// CreateAdmin 创建用户
func CreateAdmin(username, password string) (time.Time, error) {
	if username == "" || password == "" {
		return time.Time{}, errors.New("username or password cannot be empty")
	}
	createdTime := time.Now()
	status := 1
	query := `insert into public.admin (username, password, status, created_time) values ($1, $2, $3, $4)`
	_, err := database.DB.Exec(query, username, password, status, createdTime)
	if err != nil {
		return createdTime, err
	}
	return createdTime, nil
}

// SaveToken 保存token
func SaveToken(token string, id int) error {
	query := `update public.admin set token=$1 where id=$2`
	_, err := database.DB.Exec(query, token, id)
	if err != nil {
		log.Println("Error saving token:", err)
		return err
	}
	return err
}

// UpdateExpiresAtTime 更新token过期时间
func UpdateExpiresAtTime(expiresAt time.Time, id int) error {
	query := `update public.admin set token_expires_at=$1 where id=$2`
	_, err := database.DB.Exec(query, expiresAt, id)
	if err != nil {
		log.Println("Error updating expires at time:", err)
		return err
	}

	return err
}

// GetToken 获取token
func GetToken(id int) (string, error) {
	query := `select token from public.admin where id=$1`
	var token string
	err := database.DB.QueryRow(query, id).Scan(&token)
	if err != nil {
		log.Println("Error querying token:", err)
		return "", err
	}
	return token, nil
}

// ChangeUserStatus 切换用户状态
func ChangeUserStatus(id int, status int8) error {
	query := `update public.admin set status=$1 where id=$2`
	_, err := database.DB.Exec(query, status, id)
	if err != nil {
		log.Println("Error changing user status:", err)
	}
	return err
}

func GetAllUsers() ([]models.GetEmpxAdmin, error) {
	query := `select id, username, status from public.admin`
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Println("Error querying all users:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
		}
	}(rows)

	var users []models.GetEmpxAdmin
	for rows.Next() {
		var user models.GetEmpxAdmin
		err := rows.Scan(&user.ID, &user.Username, &user.Status)
		if err != nil {
			log.Println("Error scanning user:", err)
			continue
		}
		if user.ID == 1 {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserIdByToken(token string) (int, error) {
	query := `select id from public.admin where token=$1`
	var id int
	err := database.DB.QueryRow(query, token).Scan(&id)
	if err != nil {
		log.Println("Error querying user id:", err)
		return -1, err
	}
	return id, nil
}
