package service

import (
	"EmqxBackEnd/models"
	"EmqxBackEnd/repository"
	"fmt"
	"log"
)

func ProcessEmpxMessage(msg *models.EmpxMessage) error {
	UserId, err := repository.GetUserIdByNodeId(msg.NodeID)
	if err != nil {
		log.Println("未知的节点，请先注册该节点" + fmt.Sprint(msg.NodeID))
		return err
	}
	// msg.TS = time.Now()
	msg.UserId = UserId
	return repository.SaveMessage(msg)
}
