package main

import (
	_ "bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Address struct {
	ID          int    `json:"id"`
	StreetName  string `json:"streetName"`
	City        string `json:"city"`
	State       string `json:"state"`
	Customer_ID int    `json:"customerId"`
}
type Customer struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	DOB  string  `json:"dob"`
	Addr Address `json:"addr"`
}

var (
	db *sql.DB
)

func DateSubstract(d1 string) int {
	d1Slice := strings.Split(d1, "/")

	newDate := d1Slice[2] + "/" + d1Slice[1] + "/" + d1Slice[0]
	myDate, err := time.Parse("2006/01/02", newDate)

	if err != nil {
		panic(err)
	}

	return int(time.Now().Unix() - myDate.Unix())
}

func handleGetName(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:Bhavani@123go@/Customer_service")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := r.URL.Query()
	s := query.Get("name")
	fmt.Println(s, " s is")
	var c Customer
	var result []Customer
	rows, err := db.Query("SELECT * from customer INNER JOIN address ON customer.ID = address.customer_ID where customer.Name=?", s)

	if len(s) == 0 {
		rows, err = db.Query("SELECT * from customer INNER JOIN address ON customer.ID = address.customer_ID")
	}
	for rows.Next() {
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.Customer_ID); err != nil {
			log.Fatal(err)
		}
		result = append(result, c)
	}

	json.NewEncoder(w).Encode(result)

}

func handleGetId(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:Bhavani@123go@/Customer_service")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var result []Customer
	vars := mux.Vars(r)
	id := vars["id"]
	a, _ := strconv.Atoi(id)
	rows, err := db.Query("SELECT * FROM customer INNER JOIN address ON customer.ID = address.customer_ID")
	if a != 0 {
		rows, err = db.Query("SELECT * FROM customer INNER JOIN address ON customer.ID = address.customer_ID where customer.ID=?", a)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var c Customer
	for rows.Next() {

		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.Customer_ID); err != nil {
			log.Fatal(err)
		}
		result = append(result, c)
	}

	json.NewEncoder(w).Encode(result)

}

func handlePost(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:Bhavani@123go@/Customer_service")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var idr []interface{}
	var idr1 []interface{}
	var c Customer
	body, _ := ioutil.ReadAll(r.Body)
	err3 := json.Unmarshal(body, &c)
	if err3 != nil {
		log.Fatal(err3)
	}
	age := DateSubstract(c.DOB)
	if age/(365*24*3600) < 18 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Customer{})
		return
	}
	query := "insert into customer(Name, DOB) values(?,?)"
	query1 := "insert into address (StreetName,City,State,customer_ID) values(?,?,?,?)"
	idr = append(idr, &c.Name)
	idr = append(idr, &c.DOB)
	idr1 = append(idr1, &c.Addr.StreetName)
	idr1 = append(idr1, &c.Addr.City)
	idr1 = append(idr1, &c.Addr.State)
	row, err1 := db.Exec(query, idr...)
	if err1 != nil {
		log.Fatal(err1)
	}
	id, err2 := row.LastInsertId()
	if err2 == nil {
		idr1 = append(idr1, int(id))
	}
	row, err = db.Exec(query1, idr1...)
	if err != nil {
		log.Fatal(err)
	}
	id1, err2 := row.LastInsertId()
	if err2 != nil {
		log.Fatal(err2)
	}

	query = "select * from customer INNER JOIN address ON customer.ID=address.customer_ID and customer.ID=? and address.ID =?"
	var idd []interface{}
	idd = append(idd, id, id1)
	rows, err := db.Query(query, idd...)
	var result Customer
	for rows.Next() {
		rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.Customer_ID)
		result = c
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:Bhavani@123go@/Customer_service")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var c Customer
	var c1 Customer

	v, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(v, &c)
	if err != nil {
		log.Fatal(err)
	}
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := db.Query("SELECT * from customer INNER JOIN address on customer.ID=address.customer_ID and customer.ID=?", id)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&c1.ID, &c1.Name, &c1.DOB, &c1.Addr.ID, &c1.Addr.StreetName, &c1.Addr.City, &c1.Addr.State, &c1.Addr.Customer_ID); err != nil {
			log.Fatal(err)
		}
	}
	if c.DOB != "" || c.ID != 0 || c.Addr.ID != 0 || c.Addr.Customer_ID != 0 {
		fmt.Println(c.DOB, c.Addr.ID, c.Addr.Customer_ID, c.ID)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(c1)
		return
	}
	if (c.Addr.City == "") || (c.Addr.State == "") || (c.Addr.StreetName == "") || (c.Name == "") {
		fmt.Println(c.Addr.City)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(c1)
		return
	}
	_, _ = db.Exec("UPDATE customer SET customer.Name=? where customer.ID=?", c.Name, id)
	_, _ = db.Exec("UPDATE address SET address.StreetName=? , address.City=?, address.State=? where address.customer_ID=?", c.Addr.StreetName, c.Addr.City, c.Addr.State, id)
	rows, err = db.Query("SELECT * from customer INNER JOIN address on customer.ID=address.customer_ID and customer.ID=?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&c1.ID, &c1.Name, &c1.DOB, &c1.Addr.ID, &c1.Addr.StreetName, &c1.Addr.City, &c1.Addr.State, &c1.Addr.Customer_ID); err != nil {
			log.Fatal(err)
		}
	}
	json.NewEncoder(w).Encode(c1)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:Bhavani@123go@/Customer_service")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var c Customer
	vars := mux.Vars(r)
	id := vars["id"]

	val, _ := strconv.Atoi(id)
	_, err1 := db.Exec("DELETE from address where address.customer_ID=?", val)
	_, err1 = db.Exec("DELETE from customer where customer.ID=?", val)
	if err1 != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * from customer INNER JOIN address ON customer.ID=address.customer_ID where customer.ID=?", val)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.Customer_ID); err != nil {
			log.Fatal(err)
		}
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(c)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/customer", handleGetName).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id:[0-9]+}", handleGetId).Methods(http.MethodGet)
	r.HandleFunc("/customer", handlePost).Methods(http.MethodPost)
	r.HandleFunc("/customer/{id:[0-9]+}", handlePut).Methods(http.MethodPut)
	r.HandleFunc("/customer/{id:[0-9]+}", handleDelete).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", r))
}
