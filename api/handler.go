package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/reward-rabieth/gym/types"
	"github.com/reward-rabieth/gym/utils"
)

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

func (s *APIServer) HandleGetExercises(w http.ResponseWriter, r *http.Request) error {
	exerciseData, err := s.FetchExercises()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	json.NewEncoder(w).Encode(exerciseData)
	return err
}
func (s *APIServer) FetchExercises() (types.Exercise, error) {
	exUrl := utils.Loadconfig().Exurl

	resp, err := http.Get(exUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var exerciseData types.Exercise
	err = json.NewDecoder(resp.Body).Decode(&exerciseData)
	if err != nil {
		log.Fatal(err)
	}

	return exerciseData, nil
}
