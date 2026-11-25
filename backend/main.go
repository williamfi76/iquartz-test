package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// OBJECTS
type Product struct {
	ID        int     `json:"id"`
	Sku       string  `json:"sku"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	isActive  bool    `json:"is_active"`
	createdAt string  `json:"created_at"`
}

type StockMovement struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Note      string `json:"note"`
	createdAt string `json:"created_at"`
}

type SalesOrder struct {
	ID          int    `json:"id"`
	Customer    string `json:"customer"`
	status      string `json:"status"`
	total       int    `json:"total"`
	createdAt   string `json:"created_at"`
	cancelledAt string `json:"cancelled_at"`
}

type SalesOrderItems struct {
	ID        int    `json:"id"`
	OrderID   int    `json:"order_id"`
	Sku       string `json:"sku"`
	Quantity  int    `json:"quantity"`
	UnitPrice int    `json:"unit_price"`
}

func main() {
	//CONNECT TO DB
	db, err := sql.Open("postgres", os.Getenv("URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ROUTERS
	router := mux.NewRouter()
	router.HandleFunc("api/products?search={search}&page={page}&limit={limit}", getProducts(db)).Methods("GET")

	// router.HandleFunc("api/products", addNewProduct(db)).Methods("POST")
	// router.HandleFunc("api/products/{id}", updateProduct(db)).Methods("PUT")
	// router.HandleFunc("api/products/{id}", deleteProduct(db)).Methods("DELETE")
	// router.HandleFunc("api/products/{id}/stock", addNewProduct(db)).Methods("POST")
	// router.HandleFunc("api/products/{id}/movements", getProductHistory(db)).Methods("GET")
	// router.HandleFunc("api/orders", addNewOrder(db)).Methods("POST")
	// router.HandleFunc("api/orders", cancelOrder(db)).Methods("POST")
	// router.HandleFunc("api/orders/{id}", getOrderDetails(db)).Methods("GET")

	log.Println("Server running on :8080")

	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))
	log.Fatal(http.ListenAndServe(":8080", enhancedRouter))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Check if the request is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware
		next.ServeHTTP(w, r)
	})

}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set JSON Content-Type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// FUNCTIONS
func getProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Fetching products from the database...")

		// QUERY PARAMETERS
		search := r.URL.Query().Get("search")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

		rows, err := db.Query("SELECT * FROM products WHERE name ILIKE $1 OR sku ILIKE $1 LIMIT $2 OFFSET $3", "%"+search+"%", limit, (page-1)*limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		products := []Product{}
		for rows.Next() {
			var product Product
			err := rows.Scan(&product.ID, &product.Sku, &product.Name, &product.Price, &product.Stock, &product.isActive, &product.createdAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, product)
		}
		json.NewEncoder(w).Encode(products)
	}
}

// func addNewProduct(db *sql.DB) http.HandlerFunc
// func updateProduct(db *sql.DB) http.HandlerFunc
// func deleteProduct(db *sql.DB) http.HandlerFunc
// func getProductHistory(db *sql.DB) http.HandlerFunc
// func addNewOrder(db *sql.DB) http.HandlerFunc
// func cancelOrder(db *sql.DB) http.HandlerFunc
// func getOrderDetails(db *sql.DB) http.HandlerFunc
