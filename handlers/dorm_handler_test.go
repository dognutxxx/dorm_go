package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
)

func TestGetAllDorms(t *testing.T) {
	// เชื่อมต่อกับฐานข้อมูลทดสอบ
	db, err := sql.Open("postgres", "postgres://myuser:mypassword@localhost:5432/dormdb?sslmode=disable")
	if err != nil {
		t.Fatalf("ไม่สามารถเชื่อมต่อกับฐานข้อมูล: %v", err)
	}
	defer db.Close()

	// สร้าง request และ response recorder
	req, err := http.NewRequest("GET", "/api/dorms", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// เรียกใช้ handler
	handler := GetAllDorms(db)
	handler.ServeHTTP(rr, req)

	// ตรวจสอบ status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// แปลง response เป็น []Dorm
	var dorms []Dorm
	if err := json.NewDecoder(rr.Body).Decode(&dorms); err != nil {
		t.Errorf("ไม่สามารถแปลง response เป็น JSON ได้: %v", err)
	}

	// ตรวจสอบว่ามีข้อมูลหอพักอย่างน้อย 1 รายการ
	if len(dorms) == 0 {
		t.Error("ควรมีข้อมูลหอพักอย่างน้อย 1 รายการ")
	}
}
