package service

import (
	"EmqxBackEnd/models"
	"EmqxBackEnd/repository"
)

func SaveNode(node *models.Node) error {
	isHave, err := repository.CheckNode(node.ID)
	if err != nil {
		return err
	}
	if isHave {
		return repository.UpdateNode(node.ID, node.UserId)
	}
	return repository.SaveNode(node.ID, node.UserId)
}
