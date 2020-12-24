package fileserver

type Repository interface {
	Save(imageName string, content []byte) error
}
