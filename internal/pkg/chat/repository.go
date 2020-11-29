package chat

import "github.com/friends/internal/pkg/models"

type Repository interface {
	Save(models.Message) error
	GetChat(orderID int) ([]models.Message, error)
	GetVendorChats(orderIDs []int) ([]models.Chat, error)
}
