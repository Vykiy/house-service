package tests

import (
	"os"
	"testing"

	"github.com/Vykiy/house-service/internal/app"
	"github.com/Vykiy/house-service/internal/models"
	"github.com/Vykiy/house-service/internal/repository"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestIntegration(t *testing.T) {
	dbConnection := os.Getenv("DB_CONNECTION")

	db, err := sqlx.Connect("postgres", dbConnection)
	if err != nil {
		t.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("ошибка пинга базы данных: %v", err)
	}

	repo := repository.NewRepository(db)

	app := app.NewApp(repo)

	const (
		mail     = "bla"
		password = "blabla"
	)

	userID, err := app.CreateUser(mail, password, models.UserTypeModerator)
	if err != nil {
		t.Fatalf("ошибка создания пользователя: %v", err)
	}

	ok, userType, err := app.CheckUserPassword(userID, password)
	if err != nil {
		t.Fatalf("ошибка проверки пароля: %v", err)
	}

	if !ok {
		t.Fatalf("неверный пароль")
	}

	if userType != models.UserTypeModerator {
		t.Fatalf("неверный тип пользователя")
	}

	const (
		address   = "foo"
		developer = "bar"
		yearBuilt = 2021
	)

	house, err := app.CreateHouse(address, developer, yearBuilt)
	if err != nil {
		t.Fatalf("ошибка создания дома: %v", err)
	}

	if house.Address != address {
		t.Fatalf("неверный адрес дома")
	} else if house.Developer != developer {
		t.Fatalf("неверный застройщик")
	} else if house.YearBuilt != yearBuilt {
		t.Fatalf("неверный год постройки")
	}

	const (
		price = 1000000
		rooms = 3
	)

	flat, err := app.CreateFlat(house.ID, price, rooms)
	if err != nil {
		t.Fatalf("ошибка создания квартиры: %v", err)
	}

	if flat.Price != price {
		t.Fatalf("неверная цена квартиры")
	} else if flat.Rooms != rooms {
		t.Fatalf("неверное количество комнат")
	}

	flats, err := app.GetFlats(house.ID)
	if err != nil {
		t.Fatalf("ошибка получения квартир: %v", err)
	}

	if len(flats) != 1 {
		t.Fatalf("неверное количество квартир")
	}

	if flats[0].Price != price {
		t.Fatalf("неверная цена квартиры")
	} else if flats[0].Rooms != rooms {
		t.Fatalf("неверное количество комнат")
	}

	updatedFlat, err := app.UpdateFlat(flat.ID, models.FlatStatusOnModeration)
	if err != nil {
		t.Fatalf("ошибка обновления статуса квартиры: %v", err)
	}

	if updatedFlat.Status != models.FlatStatusOnModeration {
		t.Fatalf("неверный статус квартиры")
	}
}
