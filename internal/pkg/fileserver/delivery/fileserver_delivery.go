package delivery

import (
	"fmt"
	"io"

	"github.com/friends/internal/pkg/fileserver"
	"google.golang.org/grpc/metadata"
)

type FileserverDelivery struct {
	fileserverUsecase fileserver.Usecase
}

func New(fileserverUsecase fileserver.Usecase) FileserverDelivery {
	return FileserverDelivery{
		fileserverUsecase: fileserverUsecase,
	}
}

func (f FileserverDelivery) Upload(inStream fileserver.UploadService_UploadServer) error {
	file := make([]byte, 0, 1024)
	for {
		chunk, err := inStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			return err
		}

		file = append(file, chunk.Content...)
	}

	md, _ := metadata.FromIncomingContext(inStream.Context())
	fileName := md.Get("fileName")[0]
	if err := f.fileserverUsecase.Save(fileName, file); err != nil {
		return err
	}

	resp := &fileserver.UploadResponse{}

	if err := inStream.SendAndClose(resp); err != nil {
		return err
	}

	return nil
}
