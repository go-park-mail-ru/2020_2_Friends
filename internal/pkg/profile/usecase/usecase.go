package usecase

import (
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
)

type ProfileUsecase struct {
	repository profile.Repository
}

func NewProfileUsecase(repo profile.Repository) profile.Usecase {
	return ProfileUsecase{
		repository: repo,
	}
}

func (p ProfileUsecase) Create(userID string) error {
	return p.repository.Create(userID)
}

func (p ProfileUsecase) Get(userID string) (models.Profile, error) {
	return p.repository.Get(userID)
}

func (p ProfileUsecase) Update(profile models.Profile) error {
	return p.repository.Update(profile)
}

func (p ProfileUsecase) Delete(userID string) error {
	return p.repository.Delete(userID)
}
