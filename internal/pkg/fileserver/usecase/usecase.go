package usecase

import (
	"io/ioutil"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/fileserver"
)

type FileserverUsecase struct{}

func New() fileserver.Usecase {
	return FileserverUsecase{}
}

func (f FileserverUsecase) Save(imageName string, content []byte) error {
	err := ioutil.WriteFile(configs.ImageDir+imageName, content, 0666)
	if err != nil {
		return err
	}

	return nil
}
