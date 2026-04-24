package repository

import (
	"EmqxBackEnd/database"
	"EmqxBackEnd/models"
	"database/sql"
	"log"
)

// SaveMessage 保存消息
func SaveMessage(msg *models.EmpxMessage) error {
	// 注意 postgres sql 插入时，表要标明在哪个域（public）内
	query := `INSERT INTO public.message (node_id, type, message, received_at, user_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := database.DB.Exec(query, msg.NodeID, msg.Type, msg.Value, msg.TS, msg.UserId)
	return err
}

// GetMessages 读取数据库
func GetMessages(messageType, userId int) ([]models.EmpxMessage, error) {
	query := `SELECT node_id, message, received_at FROM public.message WHERE type = $1 AND user_id = $2`
	rows, err := database.DB.Query(query, messageType, userId)
	if err != nil {
		log.Printf("database.DB.Query() failed: %v", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("rows.Close() failed: %v", err)
			return
		}
	}(rows)
	var messages []models.EmpxMessage
	for rows.Next() {
		var msg models.EmpxMessage
		err := rows.Scan(&msg.NodeID, &msg.Value, &msg.TS)
		if err != nil {
			log.Println("rows.Scan() failed: " + err.Error())
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func GetMessagesByDaily(messageType, userId int, startTime, endTime string) ([]models.EmpxMessage, error) {
	query := `SELECT node_id, message, received_at 
				FROM public.message
				WHERE type = $1 AND user_id = $2 AND received_at >= $3 AND received_at <= $4
				ORDER BY received_at DESC 
				LIMIT 20`
	rows, err := database.DB.Query(query, messageType, userId, startTime, endTime)
	if err != nil {
		log.Printf("database.DB.Query() failed: %v", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("rows.Close() failed: %v", err)
			return
		}
	}(rows)
	var messages []models.EmpxMessage
	for rows.Next() {
		var msg models.EmpxMessage
		err := rows.Scan(&msg.NodeID, &msg.Value, &msg.TS)
		if err != nil {
			log.Println("rows.Scan() failed: " + err.Error())
			return nil, err
		}
		messages = append(messages, msg)
	}
	// 按时间排序
	var result []models.EmpxMessage
	for i := len(messages) - 1; i >= 0; i-- {
		result = append(result, messages[i])
		// log.Println(messages[i].TS)
	}

	return result, nil
}
