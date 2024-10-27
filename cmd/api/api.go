package api

import (
	"log"
	"net/http"

	"github.com/BerkCicekler/e-commerce-audio-api/repository"
	"github.com/BerkCicekler/e-commerce-audio-api/service/category"
	"github.com/BerkCicekler/e-commerce-audio-api/service/image"
	"github.com/BerkCicekler/e-commerce-audio-api/service/product"
	"github.com/BerkCicekler/e-commerce-audio-api/service/user"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer{
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run(mongoDatabase *mongo.Database) error {
	router:= mux.NewRouter()
	subRouter:= router.PathPrefix("/api/v1").Subrouter()

	imageHandler := image.ImageServiceHandler{}
	imageHandler.RegisterRoutes(subRouter)

	userRepository := repository.UsersRepo{
		MongoCollection: mongoDatabase.Collection("users"),
	}
	userHandler := user.UserServiceNewHandler(userRepository)
	userHandler.RegisterRoutes(subRouter)


	categoriesRepository := repository.CategoriesRepo{
		MongoCollection: mongoDatabase.Collection("categories"),
	}
	categoriesHandler := category.CategoriesServiceNewHandler(&categoriesRepository)
	categoriesHandler.RegisterRoutes(subRouter)

	productRepository := repository.ProductRepo{
		MongoCollection: mongoDatabase.Collection("products"),
	}
	productsHandler := product.ProductServiceNewHandler(&productRepository)
	productsHandler.RegisterRoutes(subRouter)

	log.Println("Listening on", s.addr)


	return http.ListenAndServe(s.addr, router)
}