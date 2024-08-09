package router

import (
	"net/http"

	"github.com/Vykiy/house-service/internal/app"
	"github.com/gorilla/mux"
)

func NewRouter(app *app.App, jwtIssuer *JWTIssuer) *mux.Router {
	router := mux.NewRouter()

	handler := NewHandler(app, jwtIssuer)

	middleware := NewMiddleware(jwtIssuer)

	router.Handle("/dummyLogin", http.HandlerFunc(handler.DummyLogin)).Methods("GET")
	router.Handle("/login", http.HandlerFunc(handler.Login)).Methods("POST")
	router.Handle("/register", http.HandlerFunc(handler.Register)).Methods("POST")
	router.Handle("/house/create", middleware.ModeratorAuth(http.HandlerFunc(handler.CreateHouse))).Methods("POST")
	router.Handle("/house/{id}", middleware.UserAuth(http.HandlerFunc(handler.GetFlats))).Methods("GET")
	router.Handle("/flat/create", middleware.UserAuth(http.HandlerFunc(handler.CreateFlat))).Methods("POST")
	router.Handle("/flat/update", middleware.ModeratorAuth(http.HandlerFunc(handler.UpdateFlat))).Methods("POST")
	router.Handle("/house/{id}/subscribe", middleware.UserAuth(http.HandlerFunc(handler.SubscribeToNewFlats))).Methods("POST")

	return router
}
