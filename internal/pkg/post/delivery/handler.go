package delivery

import (
	"encoding/json"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/post"
	"github.com/gorilla/mux"
	"net/http"
)

type PostHandler struct {
	UseCase   post.Usecase
}

func (uh *PostHandler) Find(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	related := r.URL.Query().Get("related")
	postFull, err := uh.UseCase.Find(id, related)
	if err != nil {
		res, err1 := json.Marshal(err)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write(res)
	} else {
		res, err1 := json.Marshal(postFull)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(res)
	}
}

func (uh *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	message := models.PostMessage{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post, err1 := uh.UseCase.Update(id, message.Message)
	if err1 != nil {
		res, err := json.Marshal(err1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write(res)
	} else {
		res, err := json.Marshal(post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(res)
	}
}
