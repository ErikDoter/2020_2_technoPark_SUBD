package delivery

import (
	"encoding/json"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/user"
	"github.com/gorilla/mux"
	"net/http"
)

type UserHandler struct {
	UseCase   user.UseCase
}

func (uh *UserHandler) FindByNickname(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	user, err := uh.UseCase.FindByNickname(nickname)
	if err != nil {
		http.Error(w, err.Message, http.StatusBadRequest)
		return
	}
	result, err1 := json.Marshal(user)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}