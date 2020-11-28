package usecase

import (
	"github.com/friends/internal/pkg/chat"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
)

type ChatUsecase struct {
	chatRepository    chat.Repository
	profileRepository profile.Repository
}

func New(chatRepository chat.Repository, profileRepository profile.Repository) chat.Usecase {
	return ChatUsecase{
		chatRepository:    chatRepository,
		profileRepository: profileRepository,
	}
}

func (c ChatUsecase) Save(msg models.Message) error {
	return c.chatRepository.Save(msg)
}

func (c ChatUsecase) GetChat(orderID int) ([]models.Message, error) {
	return c.chatRepository.GetChat(orderID)
}

func (c ChatUsecase) GetUserChats(userID string) ([]models.Chat, error) {
	chats, err := c.chatRepository.GetUserChats(userID)
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
