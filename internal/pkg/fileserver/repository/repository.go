package repository

import (
	"io/ioutil"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/fileserver"
)

type FileserverRepository struct{}

func New() fileserver.Repository {
	return FileserverRepository{}
}

func (f FileserverRepository) Save(imageName string, content []byte) error {
	err := ioutil.WriteFile(configs.ImageDir+imageName, content, 0666)
	if err != nil {
		return err
	}

	return nil
}
