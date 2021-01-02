package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
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

func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	user := models.User{}
	u := models.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Nickname = nickname
	users, err1 := uh.UseCase.Create(user.Nickname, user.Fullname, user.About, user.Email)
	result, err := json.Marshal(users)
	if err != nil {
		fmt.Println("erik2")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err1 != nil {
		w.WriteHeader(409)
		w.Write(result)
	} else {
		test := *users
		u = test[0]
		res, _ := json.Marshal(u)
		w.WriteHeader(201)
		w.Write(res)
	}
}

func (uh *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	user := models.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Nickname = nickname
	u, err1 := uh.UseCase.Update(user.Nickname, user.Fullname, user.About, user.Email)
	if err1 != nil {
		resultErr, err := json.Marshal(err1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err1.Message == "conflict" {
			w.WriteHeader(409)
			w.Write(resultErr)
		} else {
			w.WriteHeader(404)
			w.Write(resultErr)
		}
	} else {
		result, err := json.Marshal(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(result)
	}
}