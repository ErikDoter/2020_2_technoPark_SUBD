package delivery

import (
	"encoding/json"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/service"
	"net/http"
)

type ServiceHandler struct {
	UseCase   service.Usecase
}


func (uh *ServiceHandler) Status(w http.ResponseWriter, r *http.Request) {
	status := uh.UseCase.Status()
	res, err := json.Marshal(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}

func (uh *ServiceHandler) Clear(w http.ResponseWriter, r *http.Request) {
	uh.UseCase.Clear()
	w.WriteHeader(200)
}
