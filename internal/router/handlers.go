package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Vykiy/house-service/internal/app"
	"github.com/Vykiy/house-service/internal/models"
	"github.com/google/uuid"
)

var dummyUserID = uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

type Handler struct {
	app       *app.App
	jwtIssuer *JWTIssuer
}

func NewHandler(app *app.App, jwtIssuer *JWTIssuer) *Handler {
	return &Handler{app: app, jwtIssuer: jwtIssuer}
}

func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	userType := r.URL.Query().Get("user_type")

	if userType != string(models.UserTypeUser) && userType != string(models.UserTypeModerator) {
		http.Error(w, "тип пользователя не поддерживается", http.StatusBadRequest)
		return
	}

	token, err := h.jwtIssuer.IssueToken(models.UserType(userType), dummyUserID)
	if err != nil {
		http.Error(w, "ошибка создания токена", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	credentials := struct {
		ID       uuid.UUID `json:"id"`
		Password string    `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	successful, userType, err := h.app.CheckUserPassword(credentials.ID, credentials.Password)
	if err != nil {
		http.Error(w, "ошибка проверки пароля", http.StatusInternalServerError)
		return
	}

	if !successful {
		http.Error(w, "неверный пароль", http.StatusForbidden)
		return
	}

	token, err := h.jwtIssuer.IssueToken(models.UserType(userType), credentials.ID)
	if err != nil {
		http.Error(w, "ошибка создания токена", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	registrationData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"user_type"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&registrationData); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	userType := models.UserType(registrationData.UserType)

	if userType != models.UserTypeUser && userType != models.UserTypeModerator {
		http.Error(w, "тип пользователя не поддерживается", http.StatusBadRequest)
		return
	}

	userID, err := h.app.CreateUser(registrationData.Email, registrationData.Password, userType)
	if err != nil {
		http.Error(w, "ошибка создания пользователя", http.StatusInternalServerError)
		return
	}

	userIDJson, err := json.Marshal(struct {
		ID uuid.UUID `json:"user_id"`
	}{ID: userID})
	if err != nil {
		http.Error(w, "ошибка создания ответа", http.StatusInternalServerError)
		return
	}

	w.Write(userIDJson)
}

func (h *Handler) CreateHouse(w http.ResponseWriter, r *http.Request) {
	createHouseData := struct {
		Address   string `json:"address"`
		YearBuilt int    `json:"year"`
		Developer string `json:"developer"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&createHouseData); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if createHouseData.YearBuilt < 0 {
		http.Error(w, "неверный год постройки", http.StatusBadRequest)
		return
	}

	house, err := h.app.CreateHouse(createHouseData.Address, createHouseData.Developer, createHouseData.YearBuilt)
	if err != nil {
		http.Error(w, "ошибка создания дома", http.StatusInternalServerError)
		return
	}

	houseJson, err := json.Marshal(house)
	if err != nil {
		http.Error(w, "ошибка создания ответа", http.StatusInternalServerError)
		return
	}

	w.Write(houseJson)
}

func (h *Handler) GetFlats(w http.ResponseWriter, r *http.Request) {
	houseIDString := r.URL.Query().Get("house_id")
	if houseIDString == "" {
		http.Error(w, "не указан ID дома", http.StatusBadRequest)
		return
	}

	houseID, err := strconv.Atoi(houseIDString)
	if err != nil {
		http.Error(w, "неверный формат ID дома", http.StatusBadRequest)
		return
	}

	flats, err := h.app.GetFlats(houseID)
	if err != nil {
		http.Error(w, "ошибка получения квартир", http.StatusInternalServerError)
		return
	}

	flatsJson, err := json.Marshal(flats)
	if err != nil {
		http.Error(w, "ошибка создания ответа", http.StatusInternalServerError)
		return
	}

	w.Write(flatsJson)

}

func (h *Handler) CreateFlat(w http.ResponseWriter, r *http.Request) {
	createFlatData := struct {
		HouseID int `json:"house_id"`
		Price   int `json:"price"`
		Rooms   int `json:"rooms"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&createFlatData); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if createFlatData.Price < 0 {
		http.Error(w, "неверная цена", http.StatusBadRequest)
		return
	} else if createFlatData.Rooms < 1 {
		http.Error(w, "неверное количество комнат", http.StatusBadRequest)
		return
	}

	flat, err := h.app.CreateFlat(createFlatData.HouseID, createFlatData.Price, createFlatData.Rooms)
	if err != nil {
		http.Error(w, "ошибка создания квартиры", http.StatusInternalServerError)
		return
	}

	flatJson, err := json.Marshal(flat)
	if err != nil {
		http.Error(w, "ошибка создания ответа", http.StatusInternalServerError)
		return
	}

	w.Write(flatJson)

}

func (h *Handler) UpdateFlat(w http.ResponseWriter, r *http.Request) {
	updateFlatData := struct {
		FlatID int               `json:"flat_id"`
		Status models.FlatStatus `json:"status"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&updateFlatData); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if updateFlatData.Status != models.FlatStatusApproved && updateFlatData.Status != models.FlatStatusDeclined && updateFlatData.Status != models.FlatStatusOnModeration {
		http.Error(w, "неверный статус квартиры", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(userIDCtxKey).(uuid.UUID)
	if !ok {
		http.Error(w, "не указан ID пользователя", http.StatusBadRequest)
		return
	}

	ok, err := h.app.CheckFlatModerator(updateFlatData.FlatID, userID)
	if err != nil {
		http.Error(w, "ошибка проверки модератора", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "квартира уже модерируется другим сотрудником", http.StatusForbidden)
		return
	}

	flat, err := h.app.UpdateFlat(updateFlatData.FlatID, updateFlatData.Status)
	if err != nil {
		http.Error(w, "ошибка обновления квартиры", http.StatusInternalServerError)
		return
	}

	flatJson, err := json.Marshal(flat)
	if err != nil {
		http.Error(w, "ошибка создания ответа", http.StatusInternalServerError)
		return
	}

	w.Write(flatJson)
}

func (h *Handler) SubscribeToNewFlats(w http.ResponseWriter, r *http.Request) {
	houseIDString := r.URL.Query().Get("house_id")
	if houseIDString == "" {
		http.Error(w, "не указан ID дома", http.StatusBadRequest)
		return
	}

	houseID, err := strconv.Atoi(houseIDString)
	if err != nil {
		http.Error(w, "неверный формат ID дома", http.StatusBadRequest)
		return
	}

	subscriptionData := struct {
		Email string `json:"email"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&subscriptionData); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := h.app.SubscribeToNewFlats(houseID, subscriptionData.Email); err != nil {
		http.Error(w, "ошибка подписки на новые квартиры", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Успешно оформлена подписка"))
}
