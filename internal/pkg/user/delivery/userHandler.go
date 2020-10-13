package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/user"
)

type UserHandler struct {
	UserUsecase user.Usecase
}

func (uh UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		fmt.Println("error")
		return
	}

	err = uh.UserUsecase.Create(*user)
	if err != nil {
		w.Write([]byte(`{"created": false}`))

		return
	}
	w.Write([]byte(`{"created": true}`))
}
