package main

import (
	"ecommerce/db"
	"ecommerce/handler"
	"ecommerce/middleware"
	"ecommerce/repository"
	"ecommerce/services"
	"ecommerce/utils"
	"fmt"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db.ConnectDb()
	database := db.GetDb()

	productRepo := repository.NewProductRepo(database)
	userRepo := repository.NewUserRepo(database)
	productService := services.NewProductService(productRepo)
	userService := services.NewUserService(userRepo)
	productHandler := handler.NewProductHander(productService)
	userHandler := handler.NewUserHandler(userService)

	r := chi.NewRouter()
	verifier := utils.JWTVerifier{}

	r.Post("/login", userHandler.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middleware.Auth(verifier, next)
		})

		r.Post("/products", productHandler.CreateProduct)
		r.Get("/products/{id}", productHandler.GetProductByID)
		r.Get("/products", productHandler.GetAllProducts)
		r.Put("/products/{id}", productHandler.UpdateProduct)
		r.Delete("/products/{id}", productHandler.DeleteProducts)
	})

	r.Post("/users", userHandler.RegisterUser)
	r.Get("/users/{id}", userHandler.GetUserByID)
	r.Get("/users", userHandler.GetAllUsers)
	r.Put("/users/{id}", userHandler.UpdateUser)
	r.Delete("/users/{id}", userHandler.DeleteUser)

	fmt.Println("Server started on : 8080")
	http.ListenAndServe(":8080", r)
}
