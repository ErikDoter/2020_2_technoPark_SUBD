package delivery

import (
	"encoding/json"
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
		w.WriteHeader(404)
		w.Write(res)
	} else {
		result, err1 := json.Marshal(thread)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(result)
	}
}
