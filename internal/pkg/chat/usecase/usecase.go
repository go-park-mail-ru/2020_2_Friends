package usecase

import (
	"github.com/friends/internal/pkg/chat"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"
	"github.com/friends/internal/pkg/profile"
)

type ChatUsecase struct {
	chatRepository    chat.Repository
	profileRepository profile.Repository
	orderRepository   order.Repository
}

func New(chatRepository chat.Repository, profileRepository profile.Repository, orderRepository order.Repository) chat.Usecase {
	return ChatUsecase{
		chatRepository:    chatRepository,
		profileRepository: profileRepository,
		orderRepository:   orderRepository,
	}
}

func (c ChatUsecase) Save(msg models.Message) error {
	return c.chatRepository.Save(msg)
}

func (c ChatUsecase) GetChat(orderID int, userID string) ([]models.Message, error) {
	msgs, err := c.chatRepository.GetChat(orderID)
	if err != nil {
		return nil, err
	}

	for idx := range msgs {
		if msgs[idx].UserID == userID {
			msgs[idx].IsYourMsg = true
		} else {
			msgs[idx].IsYourMsg = false
		}
	}

	return msgs, nil
}

func (c ChatUsecase) GetVendorChats(vendorID string) ([]models.Chat, error) {
	orderIDs, err := c.orderRepository.GetVendorOrdersIDs(vendorID)
	if err != nil {
		return nil, err
	}

	chats, err := c.chatRepository.GetVendorChats(orderIDs)
	if err != nil {
		return nil, err
	}

	for idx := range chats {
		name, err := c.profileRepository.GetUsername(chats[idx].InterlocutorID)
		if err != nil {
			return nil, err
		}

		chats[idx].InterlocutorName = name
	}

	return chats, nil
}
