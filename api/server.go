package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"

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
	router.HandleFunc("/member", (makeHttpHandleFunc(s.HandleMembers)))
	router.HandleFunc("/member/{id}", WithJwtAuth(makeHttpHandleFunc(s.HandleGetMemberByid), s.Store))
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

func (s *APIServer) HandleGetMemberByid(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		id, err := GetId(r)
		if err != nil {
			return err
		}

		member, err := s.Store.GetMemberByid(id)

		if err != nil {
			return err
		}

		return WriteJson(w, http.StatusOK, member)
	}
	return fmt.Errorf("method is not allowed")
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
	tokenString, err := createJwt(member)
	if err != nil {

		fmt.Println("error in creating jwt", err)
		return err
	}
	fmt.Println("Jwt  token:", tokenString)

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

func createJwt(member *types.Gymmember) (string, error) {

	claims := &jwt.MapClaims{

		"membersId": member.Number,
		"expiresAt": 1500,
	}
	secret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}

func WithJwtAuth(handlefunc http.HandlerFunc, s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("calling with jwt auth middleware ")

		tokenstring := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenstring)

		if err != nil {

			WriteJson(w, http.StatusForbidden, ApiError{Error: "invalid token"})

			return
		}

		if !token.Valid {

			WriteJson(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		memberId, err := GetId(r)

		if err != nil {
			WriteJson(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}
		member, err := s.GetMemberByid(memberId)
		if err != nil {
			WriteJson(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		//	panic(reflect.TypeOf(claims["membersId"]))

		if member.Number != int64(claims["membersId"].(float64)) {
			WriteJson(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}
		if err != nil {
			WriteJson(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		fmt.Println(claims)
		handlefunc(w, r)
	}

}

func validateJWT(tokenString string) (*jwt.Token, error) {

	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

}
func GetId(r *http.Request) (int, error) {
	//conver the id from string to int
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idstr)
	}
	return id, nil

}
