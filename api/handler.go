package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/reward-rabieth/gym/types"
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
	// req.Name,
	// req.Age,
	// req.Gender,
	// req.Height,
	// req.Weight,
	// req.Membership,
	// req.StartDate,
	// req.EndDate,
	// req.PersonalTrainer,

	if err != nil {
		return err
	}
	if err := s.Store.CreateMember(member); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, member)
}
