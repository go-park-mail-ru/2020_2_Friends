package repository

import (
	"io/ioutil"
	"os"

	"github.com/friends/internal/pkg/fileserver"
)

type FileserverRepository struct{}

func New() fileserver.Repository {
	return FileserverRepository{}
}

func (f FileserverRepository) Save(imageName string, content []byte) error {
	path := os.Getenv("img_path")
	err := ioutil.WriteFile(path+imageName, content, 0666)
	if err != nil {
		return err
	}

	return nil
}
