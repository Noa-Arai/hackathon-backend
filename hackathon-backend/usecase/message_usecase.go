package usecase

import (
	"hackathon-backend/dao"
	"hackathon-backend/model"
)

type MessageUsecase struct {
	Repo dao.MessageRepository
}

func NewMessageUsecase(repo dao.MessageRepository) *MessageUsecase {
	return &MessageUsecase{Repo: repo}
}

func (uc *MessageUsecase) SendMessage(from string, to string, itemID int64, text string) error {
	msg := &model.Message{
		FromUserID: from,
		ToUserID:   to,
		ItemID:     itemID,
		Text:       text,
	}
	return uc.Repo.Insert(msg)
}

func (uc *MessageUsecase) GetMessages(itemID int64) ([]model.Message, error) {
	return uc.Repo.FindByItemID(itemID)
}

func (uc *MessageUsecase) ListUserRooms(userID string) ([]model.MessageRoom, error) {
	return uc.Repo.ListUserRooms(userID)
}

func (uc *MessageUsecase) MarkAsRead(itemID int64, userID string) error {
	return uc.Repo.MarkAsRead(itemID, userID)
}
