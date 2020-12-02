package usecase

import (
	"github.com/friends/internal/pkg/fileserver"
)

type FileserverUsecase struct {
	fileserverRepository fileserver.Repository
}

func New(fileserverRepository fileserver.Repository) fileserver.Usecase {
	return FileserverUsecase{
		fileserverRepository: fileserverRepository,
	}
}

func (f FileserverUsecase) Save(imageName string, content []byte) error {
	return f.fileserverRepository.Save(imageName, content)
}
