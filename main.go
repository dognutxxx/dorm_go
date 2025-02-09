package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/dognutxxx/dorm_go/go-app/docs"
	"github.com/dognutxxx/dorm_go/go-app/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title           Dorm Management API
// @version         1.0
// @description     API สำหรับระบบจัดการหอพัก
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api
// @schemes   http

func main() {
	// อ่าน environment variables สำหรับเชื่อมต่อ DB
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "myuser")
	dbPass := getEnv("DB_PASS", "mypassword")
	dbName := getEnv("DB_NAME", "dormdb")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("เชื่อมต่อฐานข้อมูลล้มเหลว: %v", err)
	}
	defer db.Close()

	// ทดสอบ ping DB
	if err := db.Ping(); err != nil {
		log.Fatalf("ping DB ไม่ผ่าน: %v", err)
	}

	fmt.Println("เชื่อมต่อ PostgreSQL สำเร็จ!")

	// Enable CORS
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// สร้าง router
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	// Swagger endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
		httpSwagger.DomID("#swagger-ui"),
	))

	// API endpoints
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/dorms", handlers.GetAllDorms(db)).Methods("GET", "OPTIONS")
	api.HandleFunc("/createDorms", handlers.CreateDorm(db)).Methods("POST", "OPTIONS")
	api.HandleFunc("/dorms/{id}", handlers.GetDorm(db)).Methods("GET", "OPTIONS")
	api.HandleFunc("/updateDorms/{id}", handlers.UpdateDorm(db)).Methods("POST", "OPTIONS")
	api.HandleFunc("/deleteDorms/{id}", handlers.DeleteDorm(db)).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/dorms/with-rooms", handlers.GetDormWithRooms(db)).Methods("GET", "OPTIONS")

	// เปิดเซิร์ฟเวอร์บน port 8080
	log.Println("Server is running on port 8080")
	log.Println("Swagger UI is available at http://localhost:8080/swagger/index.html")
	http.ListenAndServe(":8080", r)
}

// ฟังก์ชันช่วยอ่านค่า Environment Variable
func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func UpdateDorm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// โค้ดสำหรับอัปเดตหอพัก
	}
}
