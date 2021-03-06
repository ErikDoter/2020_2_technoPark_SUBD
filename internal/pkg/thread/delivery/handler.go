package delivery

import (
	"encoding/json"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/thread"
	"github.com/gorilla/mux"
	"net/http"
)

type ThreadHandler struct {
	UseCase   thread.Usecase
}

func (uh *ThreadHandler) Find(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]
	thread, err := uh.UseCase.Find(slug)
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
		result, err1 := json.Marshal(thread)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(result)
	}
}

func (uh *ThreadHandler) CreatePosts (w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]
	posts := models.Posts{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	postsAnswer, err1 := uh.UseCase.CreatePosts(slug, posts)
	if err1 != nil {
		if err1.Message == "can't find thread" {
			result, err := json.Marshal(err1)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			w.Write(result)
		} else {
			result, err := json.Marshal(err1)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(409)
			w.Write(result)
		}
	} else {
		res, err := json.Marshal(postsAnswer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write(res)
	}
}

func (uh *ThreadHandler) Update(w http.ResponseWriter, r *http.Request) {
	thread := models.ThreadUpdate{}
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t, err1 := uh.UseCase.Update(slug, thread.Message, thread.Title)
	if err1 != nil {
		result, err := json.Marshal(err1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write(result)
	} else {
		res, err := json.Marshal(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(res)
	}
}

func (uh *ThreadHandler) Vote(w http.ResponseWriter, r *http.Request) {
	vote := models.ThreadVote{}
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t, err1 := uh.UseCase.Vote(slug, vote.Nickname, vote.Voice)
	if err1 != nil {
		result, err := json.Marshal(err1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write(result)
	} else {
		res, err := json.Marshal(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(res)
	}
}


func (uh *ThreadHandler) Posts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]
	limit := r.URL.Query().Get("limit")
	desc := r.URL.Query().Get("desc")
	since := r.URL.Query().Get("since")
	sort := r.URL.Query().Get("sort")
	p, err := uh.UseCase.Posts(slug, limit, since, sort, desc)
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
		res, err1 := json.Marshal(p)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(res)
	}
}