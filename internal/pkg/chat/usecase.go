package chat

import "github.com/friends/internal/pkg/models"

type Usecase interface {
	Save(models.Message) error
	GetChat(orderID int, userID string) ([]models.Message, error)
	GetUserChats(userID string) ([]models.Chat, error)
}
