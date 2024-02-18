package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/aitalina/nocoin/internal/domain"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var jwtKey = []byte("my_secret_key")

type Restaurant interface {
	Create(ctx context.Context, restaurant domain.Restaurant) error
	GetByID(ctx context.Context, id string) (domain.Restaurant, error)
	GetAll(ctx context.Context) ([]domain.Restaurant, error)
	Update(ctx context.Context, id string, inp domain.UpdateRestaurantInput) error
}

type Profile interface {
	Create(ctx context.Context, profile domain.Profile) error
	GetByID(ctx context.Context, id string) (domain.Profile, error)
	GetPasswordByEmail(ctx context.Context, email string) (domain.Profile, error)
}

type Tip interface {
	Create(ctx context.Context, tip domain.Tip) error
}

type Handler struct {
	restaurantService Restaurant
	profileService    Profile
	tipService        Tip
}

func NewHandler(restaurant Restaurant, profile Profile, tip Tip) *Handler {
	return &Handler{
		restaurantService: restaurant,
		profileService:    profile,
		tipService:        tip,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	restaurants := r.PathPrefix("/restaurants").Subrouter()
	{
		restaurants.HandleFunc("", h.createRestaurant).Methods(http.MethodPost)
		restaurants.HandleFunc("/{id:[0-9a-zA-Z-]+}", h.getRestaurantByID).Methods(http.MethodGet)
	}

	profiles := r.PathPrefix("/profiles").Subrouter()
	{
		profiles.HandleFunc("", h.createProfile).Methods(http.MethodPost)
	}

	tips := r.PathPrefix("/tip").Subrouter()
	{
		tips.HandleFunc("", h.payTip).Methods(http.MethodPost)
	}

	login := r.PathPrefix("/login").Subrouter()
	{
		login.HandleFunc("", h.login).Methods(http.MethodPut)
	}

	return r
}

func (h *Handler) payTip(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newTip domain.Tip
	if err = json.Unmarshal(reqBytes, &newTip); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newTip.Transaction = "ytfytf"
	err = h.tipService.Create(context.TODO(), newTip)

	var message string
	if err != nil {
		message = "Payment error"
		log.Println("payTip() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message = "Payment was successful"

	response, err := json.Marshal(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getRestaurantByID(w http.ResponseWriter, r *http.Request) {
	bool, err := Protected(w, r)
	if !bool {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("getRestaurantByID() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := h.restaurantService.GetByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("getRestaurantByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(book)
	if err != nil {
		log.Println("getRestaurantByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) createRestaurant(w http.ResponseWriter, r *http.Request) {
	bool, err := Protected(w, r)
	if !bool {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newRestaurant domain.Restaurant
	if err = json.Unmarshal(reqBytes, &newRestaurant); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.restaurantService.Create(context.TODO(), newRestaurant)
	if err != nil {
		log.Println("createRestaurant() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getIdFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	id := vars["id"]

	return id, nil
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var profile domain.Profile
	if err = json.Unmarshal(reqBytes, &profile); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profile_, err := h.profileService.GetPasswordByEmail(context.TODO(), profile.Email)
	if err != nil {
		if errors.Is(err, domain.ErrProfileNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("login() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if profile_.Password != profile.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": profile.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Срок действия токена - 24 часа
	})

	// Подпись токена с помощью секретного ключа
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Ошибка создания токена")
		return
	}

	// Отправка токена в заголовке ответа
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Успешный вход в систему")
}

func (h *Handler) createProfile(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var profile domain.Profile
	if err = json.Unmarshal(reqBytes, &profile); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.profileService.Create(context.TODO(), profile)
	if err != nil {
		log.Println("createRestaurant() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	success := "Success"
	response, err := json.Marshal(success)
	if err != nil {
		log.Println("createProfile() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func Protected(w http.ResponseWriter, r *http.Request) (bool, error) {
	// Извлечение токена из заголовка авторизации
	tokenString := r.Header.Get("Authorization")[7:] // Удаление префикса "Bearer "

	// Проверка токена и извлечение данных
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Неверный метод подписи токена")
		}
		return jwtKey, nil
	})
	if err != nil {
		fmt.Fprintln(w, "Недействительный токен")
		return false, nil
	}

	// Проверка типа токена и его валидности
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Fprintln(w, "Защищенный ресурс")
		return true, nil
	}

	fmt.Fprintln(w, "Неверный токен")
	return false, nil
}
