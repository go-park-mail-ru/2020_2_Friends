package usecase

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
	"github.com/lithammer/shortuuid"
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

func (p ProfileUsecase) UpdateAvatar(userID string, file multipart.File) error {
	img, imgType, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("unsupporter img type: %w", err)
	}

	imgName := shortuuid.New()
	imgFullName := imgName + "." + imgType
	avatarFile, err := os.Create(filepath.Join(configs.FileServerPath+"/img", filepath.Base(imgFullName)))
	if err != nil {
		return fmt.Errorf("couldn't create file: %w", err)
	}

	switch imgType {
	case "png":
		err = png.Encode(avatarFile, img)
	case "jpg":
		fallthrough
	case "jpeg":
		err = jpeg.Encode(avatarFile, img, nil)
	default:
		return fmt.Errorf("unsupporter img type: %w", err)
	}

	if err != nil {
		return fmt.Errorf("couldn't encode: %w", err)
	}

	err = p.repository.UpdateAvatar(userID, imgFullName)
	if err != nil {
		return fmt.Errorf("couldn't save link to avatart: %w", err)
	}

	return nil
}

func (p ProfileUsecase) Delete(userID string) error {
	return p.repository.Delete(userID)
}
