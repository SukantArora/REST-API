package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	crud "github.com/SukantArora/Crud"
	"github.com/gorilla/mux"
)

type emp struct {
	Id    string
	Name  string
	Email string
	Role  string
}

type Handler struct {
	DB *sql.DB
}

func HomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome Home!")
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var employee emp

	// decoding data sent with request and inserting it to post

	_ = json.NewDecoder(r.Body).Decode(&employee)

	crud.InsertData("Employee_Details", DB, employee.Name, employee.Email, employee.Role)
	json.NewEncoder(w).Encode(employee)

}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empID := vars["id"]

	// for _, employee := range employees {
	// 	if employee.Id == empID {
	// 		json.NewEncoder(w).Encode(employee)
	// 	}
	// }

	empIDint, _ := strconv.Atoi(empID)
	employee, _ := crud.GetById(DB, empIDint)
	json.NewEncoder(w).Encode(employee)
}

func GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	data, err := crud.GetAll(DB)
	if err == sql.ErrNoRows {
		w.Write([]byte("No Data Found"))
	}
	if err == nil {
		json.NewEncoder(w).Encode(data)
	}

}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	empID := mux.Vars(r)["id"]
	empIDint, _ := strconv.Atoi(empID)
	// for i, employee := range employees {
	// 	if employee.Id == empID {
	// 		employees = append(employees[:i], employees[i+1:]...)
	// 		fmt.Fprintf(w, "The employee with ID %v has been deleted successfully", empID)
	// 	}
	// }

	crud.DeleteById(DB, empIDint)
	fmt.Fprintf(w, "The employee with ID %v has been deleted successfully", empID)
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	empID := mux.Vars(r)["id"]
	empIDint, _ := strconv.Atoi(empID)
	var updatedEmployee emp

	_ = json.NewDecoder(r.Body).Decode(&updatedEmployee)

	// var tobeUpdated *emp

	// for idx, employee := range employees {
	// 	if employee.Id == empID {
	// 		tobeUpdated = &employees[idx]
	// 		tobeUpdated.Email = updatedEmployee.Email
	// 		tobeUpdated.Name = updatedEmployee.Name
	// 		tobeUpdated.Role = updatedEmployee.Role
	// 		break
	// 	}
	// }

	crud.UpdateById(DB, empIDint, updatedEmployee.Name, updatedEmployee.Email, updatedEmployee.Role)
	fmt.Fprintf(w, "The employee with ID %v has been updated successfully", empID)

}

var DB *sql.DB

func main() {

	DB = crud.DbConn("Employee_Db")

	router := mux.NewRouter()
	router.HandleFunc("/", HomeLink)
	router.HandleFunc("/employee", CreateEmployee).Methods("POST")
	router.HandleFunc("/employees", GetAllEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", GetEmployee).Methods("GET")
	router.HandleFunc("/employees/{id}", UpdateEmployee).Methods("PATCH")
	router.HandleFunc("/employees/{id}", DeleteEmployee).Methods("DELETE")

	http.ListenAndServe(":8080", router)

}
