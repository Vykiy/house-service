package app

import (
	"context"
	"fmt"
	"log"

	"github.com/Vykiy/house-service/internal/models"
	"github.com/Vykiy/house-service/internal/repository"
	"github.com/Vykiy/house-service/internal/sender"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	repository *repository.Repository
	sender     *sender.Sender
}

func NewApp(repository *repository.Repository) *App {
	return &App{repository: repository}
}

func (a *App) CreateUser(email, password string, userType models.UserType) (uuid.UUID, error) {
	passwordHash, err := hashAndSalt([]byte(password))
	if err != nil {
		log.Println(fmt.Errorf("хеширование пароля: %v", err))
		return uuid.Nil, err
	}
	userID, err := a.repository.CreateUser(email, passwordHash, userType)
	if err != nil {
		log.Println(fmt.Errorf("создание пользователя: %v", err))
		return uuid.Nil, err
	}

	return userID, nil
}

func (a *App) CheckUserPassword(userID uuid.UUID, password string) (bool, models.UserType, error) {
	user, err := a.repository.GetUser(userID)
	if err != nil {
		log.Println(fmt.Errorf("получение пользователя: %v", err))
		return false, user.UserType, err
	}

	if !comparePasswords([]byte(user.PasswordHash), []byte(password)) {
		return false, user.UserType, nil
	}

	return true, user.UserType, nil
}

func (a *App) CreateHouse(address, developer string, yearBuilt int) (models.House, error) {
	house, err := a.repository.CreateHouse(address, developer, yearBuilt)
	if err != nil {
		log.Println(fmt.Errorf("создание дома: %v", err))
		return models.House{}, err
	}

	return house, nil
}

func (a *App) GetFlats(houseID int) ([]models.Flat, error) {
	flats, err := a.repository.GetFlats(houseID)
	if err != nil {
		log.Println(fmt.Errorf("получение квартир: %v", err))
		return nil, err
	}

	return flats, nil
}

func (a *App) CreateFlat(houseID int, price, rooms int) (models.Flat, error) {
	flat, err := a.repository.CreateFlat(houseID, price, rooms)
	if err != nil {
		log.Println(fmt.Errorf("создание квартиры: %v", err))
		return models.Flat{}, err
	}

	subscribers, err := a.repository.GetSubscribers(houseID)
	if err != nil {
		log.Println(fmt.Errorf("получение подписчиков: %v", err)) // не хотим прерывать выполнение функции из-за ошибки
	}

	for _, subscriber := range subscribers {
		go a.sender.SendEmail(context.Background(), subscriber, fmt.Sprintf("В доме №%d появилась новая квартира!", houseID))
	}

	return flat, nil
}

func (a *App) UpdateFlat(flatID int, status models.FlatStatus) (models.Flat, error) {
	flat, err := a.repository.UpdateFlat(flatID, status)
	if err != nil {
		log.Println(fmt.Errorf("обновление квартиры: %v", err))
		return models.Flat{}, err
	}

	return flat, nil
}

func (a *App) CheckFlatModerator(flatID int, userID uuid.UUID) (bool, error) {
	moderatorID, err := a.repository.GetFlatModerator(flatID)
	if err != nil {
		log.Println(fmt.Errorf("получение модератора квартиры: %v", err))
		return false, err
	}

	if moderatorID == uuid.Nil {
		return true, nil
	}

	return userID == moderatorID, nil
}

func (a *App) SubscribeToNewFlats(houseID int, email string) error {
	if err := a.repository.SubscribeToNewFlats(houseID, email); err != nil {
		log.Println(fmt.Errorf("подписка на новые квартиры: %v", err))
		return err
	}

	return nil
}

func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePasswords(hashedPwd, plainPwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	return err == nil
}
