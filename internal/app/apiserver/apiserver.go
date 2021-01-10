package apiserver

import (
	"fmt"
	"github.com/gorilla/mux"
	userHandler "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/user/delivery"
	userRep "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/user/repository"
	userUC "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/user/usecase"
	forumHandler "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/forum/delivery"
	forumRep "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/forum/repository"
	forumUC "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/forum/usecase"
	threadHandler "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/thread/delivery"
	threadRep "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/thread/repository"
	threadUC "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/thread/usecase"
	postHandler "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/post/delivery"
	postRep "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/post/repository"
	postUC "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/post/usecase"
	serviceHandler "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/service/delivery"
	serviceRep "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/service/repository"
	serviceUC "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/service/usecase"
	"github.com/jackc/pgx"
	"log"
	"net/http"
)

type APIServer struct {
	config    *Config
	router    *mux.Router
	store     *pgx.ConnPool
}

func New(config *Config) *APIServer {
	return &APIServer{
		config:    config,
		router:    mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	config := pgx.ConnPoolConfig{
		ConnConfig:     pgx.ConnConfig{
			Host:                 "localhost",
			Port:                 5432,
			Database:             "docker",
			User:                 "docker",
			Password:             "docker",
			TLSConfig:            nil,
			UseFallbackTLS:       false,
			FallbackTLSConfig:    nil,
			Logger:               nil,
			LogLevel:             0,
			Dial:                 nil,
			RuntimeParams:        nil,
			OnNotice:             nil,
			CustomConnInfo:       nil,
			CustomCancel:         nil,
			PreferSimpleProtocol: false,
			TargetSessionAttrs:   "",
		},
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	connPool, err := pgx.NewConnPool(config)
	if err != nil {
		log.Fatal(err)
	}
	s.store = connPool
	defer s.store.Close()

	s.configureRouter()

	fmt.Println("starting api server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) InitHandler() (userHandler.UserHandler, forumHandler.ForumHandler, threadHandler.ThreadHandler, postHandler.PostHandler, serviceHandler.ServiceHandler){
	UserRep := userRep.NewUserRepository(s.store)
	UserUC := userUC.NewUserUseCase(UserRep)
	UserHandler := userHandler.UserHandler{
		UseCase: UserUC,
	}
	ForumRep := forumRep.NewForumRepository(s.store)
	ForumUC := forumUC.NewForumUseCase(ForumRep)
	ForumHandler := forumHandler.ForumHandler{
		UseCase: ForumUC,
	}
	ThreadRep := threadRep.NewThreadRepository(s.store)
	ThreadUC := threadUC.NewThreadUseCase(ThreadRep)
	ThreadHandler := threadHandler.ThreadHandler{
		UseCase: ThreadUC,
	}
	PostRep := postRep.NewPostRepository(s.store)
	PostUC := postUC.NewPostUseCase(PostRep)
	PostHandler := postHandler.PostHandler{
		UseCase: PostUC,
	}
	ServiceRep := serviceRep.NewServiceRepository(s.store)
	ServiceUC := serviceUC.NewServiceUseCase(ServiceRep)
	ServiceHandler := serviceHandler.ServiceHandler{
		UseCase: ServiceUC,
	}
	return UserHandler, ForumHandler, ThreadHandler, PostHandler, ServiceHandler
}

func (s *APIServer) configureRouter() {
	user, forum, thread, post, service := s.InitHandler()
	//User routes ...
	s.router.HandleFunc("/api/user/{nickname}/profile", user.FindByNickname).Methods("GET")
	s.router.HandleFunc("/api/user/{nickname}/profile", user.Update).Methods("POST")
	s.router.HandleFunc("/api/user/{nickname}/create", user.Create)
	//forum
	s.router.HandleFunc("/api/forum/create", forum.Create).Methods("POST")
	s.router.HandleFunc("/api/forum/{slug}/details", forum.Find)
	s.router.HandleFunc("/api/forum/{slug}/users", forum.FindUsers).Methods("GET")
	s.router.HandleFunc("/api/forum/{slug}/create", forum.CreateThread).Methods("POST")
	s.router.HandleFunc("/api/forum/{slug}/threads", forum.ShowThreads).Methods("GET")
	//thread
	s.router.HandleFunc("/api/thread/{slug_or_id}/details", thread.Find).Methods("GET")
	s.router.HandleFunc("/api/thread/{slug_or_id}/create", thread.CreatePosts).Methods("POST")
	s.router.HandleFunc("/api/thread/{slug_or_id}/details", thread.Update).Methods("POST")
	s.router.HandleFunc("/api/thread/{slug_or_id}/vote", thread.Vote).Methods("POST")
	s.router.HandleFunc("/api/thread/{slug_or_id}/posts", thread.Posts).Methods("GET")

	//post
	s.router.HandleFunc("/api/post/{id}/details", post.Find).Methods("GET")
	s.router.HandleFunc("/api/post/{id}/details", post.Update).Methods("POST")

	//service
	s.router.HandleFunc("/api/service/status", service.Status).Methods("GET")
	s.router.HandleFunc("/api/service/clear", service.Clear)
}