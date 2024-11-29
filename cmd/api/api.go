package api

import (
	"log"
	"net/http"

	"github.com/BerkCicekler/e-commerce-audio-api/repository"
	"github.com/BerkCicekler/e-commerce-audio-api/service/basket"
	"github.com/BerkCicekler/e-commerce-audio-api/service/category"
	"github.com/BerkCicekler/e-commerce-audio-api/service/image"
	"github.com/BerkCicekler/e-commerce-audio-api/service/product"
	"github.com/BerkCicekler/e-commerce-audio-api/service/user"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIServer struct {
	addr    string
	baseDir string
}

func NewAPIServer(addr, baseDir string) *APIServer {
	return &APIServer{
		addr:    addr,
		baseDir: baseDir,
	}
}

func (s *APIServer) Run(mongoDatabase *mongo.Database) error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	imageHandler := image.ImageServiceHandler{
		BaseDir: s.baseDir,
	}
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

	basketRepository := repository.BasketRepo{
		MongoCollection: mongoDatabase.Collection("basket"),
	}
	basketHandler := basket.BasketServiceNewHandler(&basketRepository)
	basketHandler.RegisterRoutes(subRouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
