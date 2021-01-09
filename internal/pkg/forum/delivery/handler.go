package delivery

import (
	"encoding/json"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/forum"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ForumHandler struct {
	UseCase   forum.Usecase
}

func (uh *ForumHandler) Create(w http.ResponseWriter, r *http.Request) {
	f := models.Forum{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	forum, err1 := uh.UseCase.Create(f.Title, f.User, f.Slug)
	if err1 != nil {
		if err1.Message == "exist" {
			res, err := json.Marshal(forum)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(409)
			w.Write(res)
		} else {
			res, err := json.Marshal(err1)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			w.Write(res)
		}
	} else {
		res, err := json.Marshal(forum)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write(res)
	}
}

func (uh *ForumHandler) Find(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	forum, err := uh.UseCase.Find(slug)
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
		result, err1 := json.Marshal(forum)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(result)
	}
}

func (uh *ForumHandler) FindUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	var limit int
	var desc bool
	l := r.URL.Query().Get("limit")
	d := r.URL.Query().Get("desc")
	s := r.URL.Query().Get("since")
	if l == "" {
		limit = 100
	} else {
		limit, _ = strconv.Atoi(l)
	}
	if d == "" || d == "false" {
		desc = false
	} else {
		desc = true
	}
	users, err := uh.UseCase.FindUsers(slug, s, desc, limit)
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
		result, err1 := json.Marshal(users)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(result)
	}
}

func (uh *ForumHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	t := models.Thread{}
	vars := mux.Vars(r)
	slug := vars["slug"]
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.Forum = slug
	thread, err1 := uh.UseCase.CreateThread(t.Forum, t.Title, t.Author, t.Message, t.Created, t.Slug)
	if err1 != nil {
		if err1.Message == "can't find" {
			res, err := json.Marshal(err1)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			w.Write(res)
		} else {
			res, err := json.Marshal(thread)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(409)
			w.Write(res)
		}
	} else {
		thread.Slug = t.Slug
		res, err := json.Marshal(thread)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write(res)
	}
}

func (uh *ForumHandler) ShowThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	limit := r.URL.Query().Get("limit")
	desc := r.URL.Query().Get("desc")
	since := r.URL.Query().Get("since")
	threads, err := uh.UseCase.ShowThreads(slug, limit, since, desc)
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
		res, err1 := json.Marshal(threads)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(res)
	}
}