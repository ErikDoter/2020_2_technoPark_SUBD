package apiserver

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/store"
	userHandler "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/user/delivery"
	userRep "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/user/repository"
	userUC "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/user/usecase"
	"net/http"
)

type APIServer struct {
	config    *Config
	router    *mux.Router
	store     *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config:    config,
		router:    mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {

	if err := s.configureStore(); err != nil {
		return err
	}

	s.configureRouter()

	fmt.Println("starting api server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) InitHandler() (userHandler.UserHandler){
	UserRep := userRep.NewUserRepository(s.store.Db)
	UserUC := userUC.NewUserUseCase(UserRep)
	UserHandler := userHandler.UserHandler{
		UseCase: UserUC,
	}
	return UserHandler
}

func (s *APIServer) configureRouter() {
	user := s.InitHandler()
	//User routes ...
	s.router.HandleFunc("/api/user/{nickname:[A-z0-9](?:[-_\\.]?[A-z0-9])*}/profile", user.FindByNickname).Methods("GET")
	s.router.HandleFunc("/api/user/{nickname:[A-z0-9](?:[-_\\.]?[A-z0-9])*}/profile", user.Update).Methods("POST")
	s.router.HandleFunc("/api/user/{nickname:[A-z0-9](?:[-_\\.]?[A-z0-9])*}/create", user.Create)
}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}