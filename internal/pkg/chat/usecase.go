package chat

import "github.com/friends/internal/pkg/models"

//go:generate mockgen -destination=./usecase_mock.go -package=chat github.com/friends/internal/pkg/chat Usecase
type Usecase interface {
	Save(models.Message) error
	GetChat(orderID int, userID string) ([]models.Message, error)
	GetVendorChats(vendorID string) (models.VendorChatsWithInfo, error)
}
