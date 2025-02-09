package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Dorm represents the dorm model
// @Description ข้อมูลหอพัก
type Dorm struct {
	ID            int     `json:"id" example:"1"`
	Name          string  `json:"name" example:"หอพักสุขสบาย A"`
	Location      string  `json:"location" example:"ถนนพหลโยธิน 123"`
	Capacity      int     `json:"capacity" example:"50"`
	PricePerMonth float64 `json:"price_per_month" example:"4500.00"`
	IsAvailable   bool    `json:"is_available" example:"true"`
	CreatedAt     string  `json:"created_at,omitempty"`
	UpdatedAt     string  `json:"updated_at,omitempty"`
}

// DormWithRooms struct สำหรับเก็บข้อมูลหอพักพร้อมห้องพัก
type DormWithRooms struct {
	Dorm
	Rooms []Room `json:"rooms"`
}

// Room struct สำหรับเก็บข้อมูลห้องพัก
type Room struct {
	ID            int     `json:"id"`
	DormID        int     `json:"dorm_id"`
	RoomNumber    string  `json:"room_number"`
	Floor         int     `json:"floor"`
	SizeSqm       float64 `json:"size_sqm"`
	IsOccupied    bool    `json:"is_occupied"`
	PricePerMonth float64 `json:"price_per_month"`
}

// @Summary      สร้างหอพักใหม่
// @Description  สร้างข้อมูลหอพักใหม่ในระบบ
// @Tags         dorms
// @Accept       json
// @Produce      json
// @Param        dorm body Dorm true "ข้อมูลหอพัก"
// @Success      200  {object}  Dorm
// @Failure      400  {object}  string "Invalid input"
// @Failure      500  {object}  string "Server error"
// @Router       /dorms [post]
func CreateDorm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dorm Dorm
		if err := json.NewDecoder(r.Body).Decode(&dorm); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		sqlStatement := `INSERT INTO dorms (name, location, capacity, price_per_month, is_available)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`
		err := db.QueryRow(sqlStatement, dorm.Name, dorm.Location, dorm.Capacity, dorm.PricePerMonth, dorm.IsAvailable).
			Scan(&dorm.ID, &dorm.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dorm)
	}
}

// @Summary      ดึงข้อมูลหอพักทั้งหมด
// @Description  ดึงรายการหอพักทั้งหมดในระบบ
// @Tags         dorms
// @Produce      json
// @Success      200  {array}   Dorm
// @Failure      500  {object}  string "Server error"
// @Router       /dorms [get]
func GetAllDorms(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT id, name, location, capacity, price_per_month, is_available, created_at, updated_at 
			FROM dorms
			ORDER BY id
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var dorms []Dorm
		for rows.Next() {
			var d Dorm
			err := rows.Scan(
				&d.ID,
				&d.Name,
				&d.Location,
				&d.Capacity,
				&d.PricePerMonth,
				&d.IsAvailable,
				&d.CreatedAt,
				&d.UpdatedAt,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dorms = append(dorms, d)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dorms)
	}
}

// GetDormWithRooms ดึงข้อมูลหอพักพร้อมห้องพักทั้งหมด
func GetDormWithRooms(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT 
				d.id, d.name, d.location, d.capacity, d.price_per_month, d.is_available,
				r.id, r.room_number, r.floor, r.size_sqm, r.is_occupied, r.price_per_month
			FROM dorms d
			LEFT JOIN rooms r ON d.id = r.dorm_id
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		dormMap := make(map[int]*DormWithRooms)
		for rows.Next() {
			var d DormWithRooms
			var r Room
			err := rows.Scan(
				&d.ID, &d.Name, &d.Location, &d.Capacity, &d.PricePerMonth, &d.IsAvailable,
				&r.ID, &r.RoomNumber, &r.Floor, &r.SizeSqm, &r.IsOccupied, &r.PricePerMonth,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if dorm, exists := dormMap[d.ID]; exists {
				dorm.Rooms = append(dorm.Rooms, r)
			} else {
				d.Rooms = []Room{r}
				dormMap[d.ID] = &d
			}
		}

		// แปลง map เป็น slice
		dorms := make([]DormWithRooms, 0, len(dormMap))
		for _, d := range dormMap {
			dorms = append(dorms, *d)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dorms)
	}
}

// @Summary      ดึงข้อมูลหอพักตาม ID
// @Description  ดึงข้อมูลหอพักตาม ID ที่ระบุ
// @Tags         dorms
// @Produce      json
// @Param        id   path      int  true  "Dorm ID"
// @Success      200  {object}  Dorm
// @Failure      400  {object}  string "Invalid ID"
// @Failure      404  {object}  string "Not Found"
// @Failure      500  {object}  string "Server error"
// @Router       /dorms/{id} [get]
func GetDorm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		idStr := params["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var d Dorm
		sqlStatement := `SELECT id, name, location, capacity, created_at 
                         FROM dorms WHERE id = $1`
		err = db.QueryRow(sqlStatement, id).Scan(&d.ID, &d.Name, &d.Location, &d.Capacity, &d.CreatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(d)
	}
}

// UPDATE
func UpdateDorm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// (1) ดึงค่า ID จาก URL Param (เช่น /dorms/{id})
		// ถ้าใช้ chi หรือ gorilla/mux:
		// idStr := chi.URLParam(r, "id")          // หรือ mux.Vars(r)["id"]
		// id, err := strconv.Atoi(idStr)
		// ถ้าไม่ได้ใช้ router อื่น ๆ อาจต้อง parse เอง
		params := mux.Vars(r)
		idStr := params["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// (2) Decode JSON จาก Body
		var dorm Dorm
		if err := json.NewDecoder(r.Body).Decode(&dorm); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// อาจเพิ่มการตรวจสอบค่าว่าง / รูปแบบข้อมูลที่จำเป็นก่อนก็ได้

		// (3) เตรียมคำสั่ง SQL และอัปเดตค่า
		sqlStatement := `
			UPDATE dorms
			SET
				name = $1,
				location = $2,
				capacity = $3,
				price_per_month = $4,
				is_available = $5,
				updated_at = NOW()
			WHERE id = $6
			RETURNING id, created_at, updated_at
		`

		// (4) Execute Query พร้อมกับ Scan ค่าที่คืนกลับ (หากต้องการ)
		err = db.QueryRow(
			sqlStatement,
			dorm.Name,
			dorm.Location,
			dorm.Capacity,
			dorm.PricePerMonth,
			dorm.IsAvailable,
			id, // WHERE id = $6
		).Scan(
			&dorm.ID,
			&dorm.CreatedAt,
			&dorm.UpdatedAt,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Dorm not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// (5) ตอบกลับเป็น JSON
		// (เนื่องจากเรา Scan คืนมาได้แล้ว จึงส่ง dorm ที่อัปเดตกลับไปได้เลย)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dorm)
	}
}

// DELETE
func DeleteDorm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		idStr := params["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		sqlStatement := `DELETE FROM dorms WHERE id=$1`
		res, err := db.Exec(sqlStatement, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	r := mux.NewRouter()
	// ... ตั้งค่าเส้นทาง API ของคุณ ...

	// ตั้งค่า CORS ให้เปิดให้ทุกโดเมน
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // อนุญาตให้ทุกโดเมนเข้าถึง
		AllowCredentials: true,
	})

	// ใช้ CORS Middleware
	handler := c.Handler(r)
	http.ListenAndServe(":8080", handler)
}
