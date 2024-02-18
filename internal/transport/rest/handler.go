package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gitlab.com/aitalina/nocoin/internal/domain"
	"gitlab.com/aitalina/nocoin/internal/dto"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var jwtKey = []byte("my_secret_key")

type Restaurant interface {
	Create(ctx context.Context, restaurant domain.Restaurant) error
	GetByOwnerId(ctx context.Context, id string) (domain.Restaurant, error)
}

type Profile interface {
	Create(ctx context.Context, profile domain.Profile) error
	FindProfileByEmail(ctx context.Context, email string) (domain.Profile, error)
}

type Tip interface {
	Create(ctx context.Context, tip domain.Tip) error
	FindAllByOwnerId(ctx context.Context, ownerId string) ([]dto.Tip, error)
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

	profiles := r.PathPrefix("/profile").Subrouter()
	{
		profiles.HandleFunc("", h.createProfileAndRestaurant).Methods(http.MethodPost)
	}

	tips := r.PathPrefix("/tip").Subrouter()
	{
		tips.HandleFunc("", h.payTip).Methods(http.MethodPost)
	}

	login := r.PathPrefix("/login").Subrouter()
	{
		login.HandleFunc("", h.loginAndGetRestaurantInfo).Methods(http.MethodPut)
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
	newTip.CreateDate = time.Now()
	//тут должна быть транзакция в банк (stripe)
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

func (h *Handler) createProfileAndRestaurant(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var profileDto dto.Profile
	if err = json.Unmarshal(reqBytes, &profileDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profile := domain.Profile{
		Name:     profileDto.Name,
		Email:    profileDto.Email,
		Password: profileDto.Password,
		Role:     "OWNER",
	}

	err = h.profileService.Create(context.TODO(), profile)
	if err != nil {
		log.Println("createRestaurant() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	profileFromDB, err := h.profileService.FindProfileByEmail(context.TODO(), profileDto.Email)
	if err != nil {
		log.Println("createRestaurant() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newRestaurant := domain.Restaurant{
		Name:    profileDto.RestaurantName,
		OwnerId: profileFromDB.ID,
	}
	err = h.restaurantService.Create(context.TODO(), newRestaurant)
	if err != nil {
		log.Println("createRestaurant() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	restaurantFromDB, err := h.restaurantService.GetByOwnerId(context.TODO(), profileFromDB.ID)
	if err != nil {
		log.Println("createRestaurant() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(restaurantFromDB.ID)
	if err != nil {
		log.Println("createProfileAndRestaurant() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (h *Handler) loginAndGetRestaurantInfo(w http.ResponseWriter, r *http.Request) {
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

	//ищем профиль в бд
	profileFromDB, err := h.profileService.FindProfileByEmail(context.TODO(), profile.Email)
	if err != nil {
		if errors.Is(err, domain.ErrProfileNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("loginAndGetRestaurantInfo() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//проверяем корректность пароля
	if profileFromDB.Password != profile.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//ищем чаевые по владельцу
	tips, err := h.tipService.FindAllByOwnerId(context.TODO(), profileFromDB.ID)
	if err != nil {
		if errors.Is(err, domain.ErrRestaurantNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("getRestaurantInfoWithTipSum() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	restaurantFromDB, err := h.restaurantService.GetByOwnerId(context.TODO(), profileFromDB.ID)
	if err != nil {
		log.Println("createRestaurant() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tipInfo := dto.TipInfo{
		Tips:         tips,
		RestaurantId: restaurantFromDB.ID,
	}

	response, err := json.Marshal(tipInfo)
	if err != nil {
		log.Println("getRestaurantInfoWithTipSum() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)

	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"username": profile.Email,
	//	"exp":      time.Now().Add(time.Hour * 24).Unix(), // Срок действия токена - 24 часа
	//})
	//
	//// Подпись токена с помощью секретного ключа
	//tokenString, err := token.SignedString(jwtKey)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	fmt.Fprintln(w, "Ошибка создания токена")
	//	return
	//}
	//
	//// Отправка токена в заголовке ответа
	//w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintln(w, "Успешный вход в систему")
}

//func Protected(w http.ResponseWriter, r *http.Request) (bool, error) {
//	// Извлечение токена из заголовка авторизации
//	tokenString := r.Header.Get("Authorization")[7:] // Удаление префикса "Bearer "
//
//	// Проверка токена и извлечение данных
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("Неверный метод подписи токена")
//		}
//		return jwtKey, nil
//	})
//	if err != nil {
//		fmt.Fprintln(w, "Недействительный токен")
//		return false, nil
//	}
//
//	// Проверка типа токена и его валидности
//	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		fmt.Fprintln(w, "Защищенный ресурс")
//		return true, nil
//	}
//
//	fmt.Fprintln(w, "Неверный токен")
//	return false, nil
//}
