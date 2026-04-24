package repository

import (
	"EmqxBackEnd/database"
	"EmqxBackEnd/models"
	"database/sql"
	"log"
)

// GetUserIdByNodeId 根据节点id获取用户id
func GetUserIdByNodeId(nodeId int) (int, error) {
	query := `select user_id from public.node where id=$1`
	var userId int
	err := database.DB.QueryRow(query, nodeId).Scan(&userId)
	if err != nil {
		log.Println("GetUserIdByNodeId error:", err)
		return 0, err
	}
	return userId, nil
}

// SaveNode 绑定节点
func SaveNode(nodeId, userId int) error {
	query := `insert into public.node(id, user_id) values($1, $2)`
	_, err := database.DB.Exec(query, nodeId, userId)
	return err
}

func UpdateNode(nodeId, userId int) error {
	query := `update public.node set user_id=$1 where id=$2`
	_, err := database.DB.Exec(query, userId, nodeId)
	return err
}

func CheckNode(nodeId int) (bool, error) {
	query := `select exists(select 1 from public.node where id=$1)`
	var exists bool
	err := database.DB.QueryRow(query, nodeId).Scan(&exists)
	if err != nil {
		log.Println("CheckNode error:", err)
		return false, err
	}
	return exists, nil
}

func GetAllNodeByUserId(userId int) ([]models.Node, error) {
	var query string
	var rows *sql.Rows
	var err error
	if userId == 1 {
		query = `select id from public.node`
		rows, err = database.DB.Query(query)
	} else {
		query = `select id from public.node where user_id=$1`
		rows, err = database.DB.Query(query, userId)
	}

	if err != nil {
		log.Println("GetAllNodeByUserId error:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows in GetAllNodeByUserId:", err)
			return
		}
	}(rows)
	var nodes []models.Node
	for rows.Next() {
		var node models.Node
		err := rows.Scan(&node.ID)
		if err != nil {
			log.Println("GetAllNodeByUserId error:", err)
			return nil, err
		}
		err = rows.Scan(&node.UserId)
		if err != nil {
			log.Println("GetAllNodeByUserId error:", err)
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func GetAllNode() ([]models.Node, error) {
	query := "select * from public.node"
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Println("GetAllNode error:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows in GetAllNode:", err)
			return
		}
	}(rows)
	var nodes []models.Node
	for rows.Next() {
		var node models.Node
		err := rows.Scan(&node.ID, &node.UserId)
		if err != nil {
			log.Println("GetAllNode error:", err)
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
