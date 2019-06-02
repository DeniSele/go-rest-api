package main

import (
	"GoTask/database"
	"GoTask/handler"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Server running...")

	database.Init()
	r := mux.NewRouter()

	r.HandleFunc("/users", handler.GetUsers).Methods("GET")
	r.HandleFunc("/user/{secondname}", handler.GetUserBySName).Methods("GET")
	r.HandleFunc("/users", handler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	r.HandleFunc("/account/{id}", handler.GetBalance).Methods("GET")
	r.HandleFunc("/account", handler.AddAccount).Methods("POST")
	r.HandleFunc("/account/{id}", handler.CloseAccount).Methods("DELETE")

	r.HandleFunc("/transactions", handler.GetTransactions).Methods("GET")
	r.HandleFunc("/transaction/{id}", handler.TransferMoney).Methods("PUT")
	r.HandleFunc("/transaction/{id}", handler.DeclineTransaction).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
