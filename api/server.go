package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/reward-rabieth/gym/storage"
	"github.com/reward-rabieth/gym/utils"
	"go.uber.org/zap"
)

type APIServer struct {
	ListenAddress string
	Store         storage.Storage
}

func NewApiServer(listenAddress string, store storage.Storage) *APIServer {

	return &APIServer{
		ListenAddress: listenAddress,
		Store:         store,
	}
}

func (s *APIServer) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/member", makeHttpHandleFunc(s.HandleMembers))
	addr := fmt.Sprintf(":%s", s.ListenAddress)
	utils.Logger.Info(fmt.Sprintf("json api server is running on port %s", s.ListenAddress))

	if err := http.ListenAndServe(addr, router); err != nil {
		utils.Logger.Fatal("Server stopped", zap.Error(err))
	}
	utils.Logger.Sync()

}
