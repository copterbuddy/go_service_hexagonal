package main

import (
	"bank/handler"
	"bank/logs"
	"bank/repository"
	"bank/service"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func main() {
	initTimeZone()
	initConfig()
	db := initDatabase()

	customerRepositoryDB := repository.NewCustomerRepositoryDB(db)
	// customerRepositoryMock := repository.NewCustomerRepositoryMock() //ลบเพื่อไปใช้ gomock
	customerService := service.NewCustomerService(customerRepositoryDB)
	customerHandler := handler.NewCustomerHandler(customerService)

	_ = customerRepositoryDB
	// _ = customerRepositoryMock //ลบเพื่อไปใช้ gomock
	_ = customerService
	_ = customerHandler

	accountRepositoryDB := repository.NewAccountRepositoryDB(db)
	accountService := service.NewAccountService(accountRepositoryDB)
	accountHandler := handler.NewAccountHandler(accountService)

	_ = accountRepositoryDB
	_ = accountService
	_ = accountHandler

	//Get From Db
	/*customers, err := customerRepository.GetAll()
	if err != nil {
		panic(err)
	}

	fmt.Println(customers)

	customer, err := customerRepository.GetById(1)
	if err != nil {
		panic(err)
	}

	fmt.Println(customer)*/
	//End Get From Db

	//Get From Service
	/*customers, err := customerService.GetCustomers()
	if err != nil {
		panic(err)
	}

	fmt.Println(customers)

	customer, err := customerService.GetCustomer(1)
	if err != nil {
		panic(err)
	}

	fmt.Println(customer)*/
	//End Get From Service

	router := mux.NewRouter()

	router.HandleFunc("/customers", customerHandler.GetCustomers).Methods(http.MethodGet)

	router.HandleFunc("/customers/{customerID:[0-9]+}", customerHandler.GetCustomer).Methods(http.MethodGet)

	router.HandleFunc("/customers/{customerID:[0-9]+}/accounts", accountHandler.GetAccounts).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customerID:[0-9]+}/accounts", accountHandler.NewAccount).Methods(http.MethodPost)

	// log.Printf("Banking service started at port %v", viper.GetInt("app.port"))
	logs.Info("Banking service started at port " + viper.GetString("app.port"))
	http.ListenAndServe(fmt.Sprintf(":%v", viper.GetInt("app.port")), router)
}

func initConfig() {
	//เรียก ค่า config จาก yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv() //Check Env ก่อน config
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}

func initDatabase() *sqlx.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.database"),
	)

	// logs.Info("username = " + viper.GetString("db.username"))
	// logs.Info("password = " + viper.GetString("db.password"))
	// logs.Info("host = " + viper.GetString("db.host"))
	// logs.Info("port = " + viper.GetString("db.port"))
	// logs.Info("database = " + viper.GetString("db.database"))
	db, err := sqlx.Open(viper.GetString("db.driver"), dsn)
	if err != nil {
		panic(err)
	}

	//set timeout
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
