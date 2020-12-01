package fileserver

type Usecase interface {
	Save(imageName string, content []byte) error
}
