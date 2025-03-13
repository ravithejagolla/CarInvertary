package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Car struct {
	ID    int64
	Name  string
	Model string
	Brand string
	Year  int64
	Price float64
}

var Cars = make(map[int64]string)

var mu sync.Mutex

func carHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	entity := strings.TrimPrefix(path, "/cars")
	entity = strings.Trim(entity, "/")

	switch r.Method {
	case "POST":
		if entity == "" {
			createHandler(w, r)
		} else {
			http.Error(w, "IncorrectPostMethod", http.StatusBadRequest)

		}
	case "GET":
		if entity == "" {
			http.Error(w, "Dont Have data create", http.StatusBadRequest)
		} else {
			id, _ := strconv.Atoi(entity)
			getHandler(w, r, int64(id))
		}
	case "DELETE":
		if entity == "" {
			http.Error(w, "IncorrectDeleteMethod", http.StatusBadRequest)
		} else {
			id, _ := strconv.Atoi(entity)
			deleteHandler(w, r, int64(id))
		}
	default:
		http.Error(w, "IncorrectMethod", http.StatusBadRequest)

	}
}

//Create Car to Our Inventary System

func createHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var car Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id = rand.Intn(10000)
	car.ID = int64(id)

	carJSON, err := json.Marshal(car)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Cars[car.ID] = string(carJSON)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(car)

}

func getHandler(w http.ResponseWriter, r *http.Request, id int64) {
	mu.Lock()
	defer mu.Unlock()

	car, exists := Cars[id]
	if !exists {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, Cars[id])
	w.Write([]byte(car))

}

func deleteHandler(w http.ResponseWriter, r *http.Request, id int64) {
	mu.Lock()
	defer mu.Unlock()

	var _, ok = Cars[id]
	if !ok {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	delete(Cars, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}

func main() {
	http.HandleFunc("/cars", carHandler)
	http.HandleFunc("/cars/", carHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

