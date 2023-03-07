package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/reward-rabieth/gym/storage"
	"github.com/reward-rabieth/gym/types"
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
	router.HandleFunc("/member", makeHttpHandleFunc(s.HandleCreateMember))
	router.HandleFunc("/exercise", makeHttpHandleFunc(s.HandleExercises))
	addr := fmt.Sprintf(":%s", s.ListenAddress)
	utils.Logger.Info(fmt.Sprintf("json api server is running on port %s", s.ListenAddress))

	if err := http.ListenAndServe(addr, router); err != nil {
		utils.Logger.Fatal("Server stopped", zap.Error(err))
	}
	utils.Logger.Sync()

}
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-type", "Application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type Apifunc func(w http.ResponseWriter, r *http.Request) error
type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandleFunc(f Apifunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}

}
func (s *APIServer) HandleMembers(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.HandleGetMembers(w, r)
	}

	if r.Method == "POST" {

		return s.HandleCreateMember(w, r)
	}
	return fmt.Errorf("method is not allowed %s", r.Method)
}
func (s *APIServer) HandleExercises(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.HandleGetExercises(w, r)
	}

	if r.Method == "POST" {

		return s.HandleCreateExercise(w, r)
	}
	return fmt.Errorf("method is not allowed %s", r.Method)
}

func (s *APIServer) HandleGetMembers(w http.ResponseWriter, r *http.Request) error {

	members, err := s.Store.GetMembers()
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, members)
}

func (s *APIServer) HandleCreateMember(w http.ResponseWriter, r *http.Request) error {

	req := new(types.CreateGymMemberRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	member, err := types.NewGymMember(*req)

	if err != nil {
		return err
	}
	if err := s.Store.CreateMember(member); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, member)
}
func (s *APIServer) HandleCreateExercise(w http.ResponseWriter, r *http.Request) error {
	req := new(types.CreateNewExerciseRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		return err
	}

	exercise, err := types.NewExercise(*req)
	if err != nil {
		return err
	}
	if err := s.Store.CreateExercise(exercise); err != nil {
		return err
	}
	return WriteJson(w, http.StatusCreated, exercise)

}

func (s *APIServer) HandleGetExercises(w http.ResponseWriter, r *http.Request) error {
	exercise, err := s.Store.GetExercises()
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, exercise)
}
