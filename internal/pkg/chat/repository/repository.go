package repository

import (
	"database/sql"
	"fmt"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/chat"
	"github.com/friends/internal/pkg/models"
)

type ChatRepository struct {
	db *sql.DB
}

func New(db *sql.DB) chat.Repository {
	return ChatRepository{
		db: db,
	}
}

func (c ChatRepository) Save(msg models.Message) error {
	_, err := c.db.Exec(
		"INSERT INTO messages (orderID, userID, message_text, sent_at) VALUES ($1, $2, $3, $4)",
		msg.OrderID, msg.UserID, msg.Text, msg.SentAt,
	)

	if err != nil {
		return fmt.Errorf("couldn't insert message on order %v from user with id %v. Error: %w", msg.OrderID, msg.UserID, err)
	}

	return nil
}

func (c ChatRepository) GetChat(orderID int) ([]models.Message, error) {
	rows, err := c.db.Query(
		"SELECT userID, message_text, sent_at FROM messages WHERE orderID = $1",
		orderID,
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get messages for order id %v. Error: %w", orderID, err)
	}
	defer rows.Close()

	msgs := make([]models.Message, 0)
	var msg models.Message
	for rows.Next() {
		err = rows.Scan(&msg.UserID, &msg.Text, &msg.SentAt)
		if err != nil {
			return nil, fmt.Errorf("couldn't get msg for order id %v. Error: %w", orderID, err)
		}
		msg.SentAtStr = msg.SentAt.Format(configs.TimeFormat)

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (c ChatRepository) GetUserChats(userID string) ([]models.Chat, error) {
	rows, err := c.db.Query(
		`SELECT orderID, userID, message_text FROM messages
		WHERE userID = $1 AND sent_at IN (SELECT MAX(sent_at) FROM messages GROUP BY orderID, userID)`,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get chats for userID %s. Error: %w", userID, err)
	}
	defer rows.Close()

	chats := make([]models.Chat, 0)
	chat := models.Chat{}
	for rows.Next() {
		err = rows.Scan(&chat.OrderID, &chat.InterlocutorID, &chat.LastMsg)
		if err != nil {
			return nil, fmt.Errorf("couldn't get chat for userID %s. Error: %w", userID, err)
		}

		chats = append(chats, chat)
	}

	return chats, nil
}
