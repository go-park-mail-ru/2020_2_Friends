package usecase

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/friends/internal/pkg/fileserver"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
	"github.com/lithammer/shortuuid"
	"google.golang.org/grpc/metadata"
)

type ProfileUsecase struct {
	repository profile.Repository
	fsClient   fileserver.UploadServiceClient
}

func NewProfileUsecase(repo profile.Repository, fileserverClient fileserver.UploadServiceClient) profile.Usecase {
	return ProfileUsecase{
		repository: repo,
		fsClient:   fileserverClient,
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

func (p ProfileUsecase) UpdateAvatar(userID string, file multipart.File, imageType string) (string, error) {
	imgName := shortuuid.New()
	imgFullName := imgName + "." + imageType

	md := metadata.New(map[string]string{"fileName": imgFullName})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := p.fsClient.Upload(ctx)
	if err != nil {
		return "", err
	}

	write := true
	chunk := make([]byte, 1024)

	for write {
		size, err := file.Read(chunk)
		if err != nil {
			if err == io.EOF {
				write = false
				continue
			}
			return "", err
		}
		err = stream.Send(&fileserver.Chunk{Content: chunk[:size]})
		if err != nil {
			return "", err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	err = p.repository.UpdateAvatar(userID, imgFullName)
	if err != nil {
		return "", err
	}

	return imgFullName, nil
}

func (p ProfileUsecase) UpdateAddresses(userID string, addresses []string) error {
	return p.repository.UpdateAddresses(userID, addresses)
}

func (p ProfileUsecase) Delete(userID string) error {
	return p.repository.Delete(userID)
}
