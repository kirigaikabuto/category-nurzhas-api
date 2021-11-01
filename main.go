package main

import (
	"fmt"
	"github.com/gorilla/mux"
	category_nurzhas_store "github.com/kirigaikabuto/category-nurzhas-store"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"log"
	"net/http"
)

var (
	postgresUser            = "oaxlkqvpikdard"
	postgresPassword        = "79a272cdf4249041aa90183895ff92d9b2d1e6bd69cd5165552f98c6f0e634bd"
	postgresDatabaseName    = "dd4k5rjp3rmvg1"
	postgresHost            = "ec2-44-194-54-123.compute-1.amazonaws.com"
	postgresPort            = 5432
	postgresParams          = ""
	s3endpoint              = "https://s3.us-east-2.amazonaws.com"
	s3bucket                = "setdata"
	s3accessKey             = "AKIA54CP6OJQEHUI6KFO"
	s3secretKey             = "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"
	s3uploadedFilesBasePath = "https://setdata.s3.us-east-2.amazonaws.com"
	s3region                = "us-east-2"
	port                    = "8080"
)

func main() {
	postgreConfig := category_nurzhas_store.PostgresConfig{
		Host:             postgresHost,
		Port:             postgresPort,
		User:             postgresUser,
		Password:         postgresPassword,
		Database:         postgresDatabaseName,
		Params:           postgresParams,
		ConnectionString: "",
	}
	store, err := category_nurzhas_store.NewPostgresCategoryStore(postgreConfig)
	if err != nil {
		log.Fatal(err)
		return
	}
	usersStore, err := category_nurzhas_store.NewPostgresUsersStore(postgreConfig)
	if err != nil {
		log.Fatal(err)
		return
	}
	usersService := category_nurzhas_store.NewUserService(usersStore)
	if err != nil {
		log.Fatal(err)
		return
	}
	usersHttpEndpoints := category_nurzhas_store.NewHttpEndpoints(setdata_common.NewCommandHandler(usersService))

	s3, err := category_nurzhas_store.NewS3Uploader(
		s3endpoint,
		s3accessKey,
		s3secretKey,
		s3bucket,
		s3uploadedFilesBasePath,
		s3region)
	if err != nil {
		log.Fatal(err)
		return
	}
	service := category_nurzhas_store.NewCategoryService(store, s3)
	ch := setdata_common.NewCommandHandler(service)
	httpEndpoints := category_nurzhas_store.NewHttpEndpoints(ch)
	router := mux.NewRouter()
	router.Methods("PUT").Path("/prices-file").HandlerFunc(httpEndpoints.MakeUploadPricesFile())
	router.Methods("GET").Path("/prices-file").HandlerFunc(httpEndpoints.MakeGetUploadPricesFile())

	router.Methods("POST").Path("/category").HandlerFunc(httpEndpoints.MakeCreateCategoryEndpoint())
	router.Methods("GET").Path("/category").HandlerFunc(httpEndpoints.MakeGetCategoryEndpoint())
	router.Methods("PUT").Path("/category/image").HandlerFunc(httpEndpoints.MakeUploadCategoryImageEndpoint())
	router.Methods("GET").Path("/category/list").HandlerFunc(httpEndpoints.MakeListCategoryEndpoint())

	router.Methods("POST").Path("/users/register").HandlerFunc(usersHttpEndpoints.MakeRegisterUserEndpoint())

	fmt.Println("api is running on port " + port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
