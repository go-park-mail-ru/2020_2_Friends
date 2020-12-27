package delivery

import (
	"context"

	"github.com/friends/internal/pkg/session"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SessionDelivery struct {
	sessionUsecase session.Usecase
}

func NewSessionDelivery(usecase session.Usecase) session.SessionWorkerServer {
	return SessionDelivery{
		sessionUsecase: usecase,
	}
}

func (s SessionDelivery) Create(ctx context.Context, userID *session.UserID) (*session.SessionName, error) {
	sessionName, err := s.sessionUsecase.Create(userID.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "couldn't create session for user with id = %s. Error: %v", userID.GetId(), err,
		)
	}

	resp := session.SessionName{
		Name: sessionName,
	}

	return &resp, nil
}

func (s SessionDelivery) Check(ctx context.Context, sessionName *session.SessionName) (*session.UserID, error) {
	id, err := s.sessionUsecase.Check(sessionName.GetName())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "couldn't check session with name = %v. Error: %v", sessionName.GetName(), err,
		)
	}

	resp := session.UserID{
		Id: id,
	}

	return &resp, nil
}

func (s SessionDelivery) Delete(
	ctx context.Context, sessionName *session.SessionName,
) (
	*session.DeleteResponse, error,
) {
	err := s.sessionUsecase.Delete(sessionName.GetName())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "couldn't delete session with name = %v. Error: %v", sessionName.GetName(), err,
		)
	}

	return &session.DeleteResponse{}, nil
}

func (s SessionDelivery) SetCSRFToken(
	ctx context.Context, userID *session.UserID,
) (
	*session.Token, error,
) {
	token, err := s.sessionUsecase.SetCSRFToken(userID.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "couldn't set token for user with id = %v. Error: %v", userID.GetId(), err,
		)
	}

	return &session.Token{Value: token}, nil
}

func (s SessionDelivery) GetTokenFromUser(
	ctx context.Context, userID *session.UserID,
) (
	*session.Token, error,
) {
	token, err := s.sessionUsecase.GetTokenFromUser(userID.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "couldn't get token for user with id = %v. Error: %v", userID.GetId(), err,
		)
	}

	return &session.Token{Value: token}, nil
}
