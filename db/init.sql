-- สร้างตาราง dorms
CREATE TABLE IF NOT EXISTS dorms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(255),
    capacity INT,
    price_per_month DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- สร้างตาราง rooms
CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    dorm_id INT REFERENCES dorms(id),
    room_number VARCHAR(20) NOT NULL,
    floor INT,
    size_sqm DECIMAL(5,2),
    is_occupied BOOLEAN DEFAULT false,
    price_per_month DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- สร้างตาราง tenants (ผู้เช่า)
CREATE TABLE IF NOT EXISTS tenants (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20),
    emergency_contact VARCHAR(100),
    emergency_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- สร้างตาราง contracts (สัญญาเช่า)
CREATE TABLE IF NOT EXISTS contracts (
    id SERIAL PRIMARY KEY,
    room_id INT REFERENCES rooms(id),
    tenant_id INT REFERENCES tenants(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    deposit_amount DECIMAL(10,2),
    monthly_rent DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- สร้างตาราง payments (การชำระเงิน)
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    contract_id INT REFERENCES contracts(id),
    amount DECIMAL(10,2) NOT NULL,
    payment_date DATE NOT NULL,
    payment_type VARCHAR(50),
    payment_status VARCHAR(20),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- สร้างตาราง maintenance_requests (คำขอซ่อมบำรุง)
CREATE TABLE IF NOT EXISTS maintenance_requests (
    id SERIAL PRIMARY KEY,
    room_id INT REFERENCES rooms(id),
    tenant_id INT REFERENCES tenants(id),
    issue_description TEXT NOT NULL,
    priority VARCHAR(20) DEFAULT 'normal',
    status VARCHAR(20) DEFAULT 'pending',
    reported_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_date TIMESTAMP,
    notes TEXT
);

-- เพิ่มข้อมูลตัวอย่างในตาราง dorms
INSERT INTO dorms (name, location, capacity, price_per_month, is_available) VALUES
('หอพักสุขสบาย A', 'ถนนพหลโยธิน 123', 50, 4500.00, true),
('หอพักสุขสบาย B', 'ถนนพหลโยธิน 124', 40, 4000.00, true),
('หอพักนักศึกษา C', 'ถนนวิภาวดี 99', 30, 3500.00, true);

-- เพิ่มข้อมูลตัวอย่างในตาราง rooms
INSERT INTO rooms (dorm_id, room_number, floor, size_sqm, is_occupied, price_per_month) VALUES
(1, 'A101', 1, 25.5, false, 4500.00),
(1, 'A102', 1, 25.5, true, 4500.00),
(1, 'A201', 2, 25.5, false, 4500.00),
(2, 'B101', 1, 22.0, false, 4000.00),
(2, 'B102', 1, 22.0, true, 4000.00),
(3, 'C101', 1, 20.0, false, 3500.00);

-- เพิ่มข้อมูลตัวอย่างในตาราง tenants
INSERT INTO tenants (first_name, last_name, email, phone, emergency_contact, emergency_phone) VALUES
('สมชาย', 'ใจดี', 'somchai@email.com', '081-111-1111', 'สมศรี ใจดี', '081-222-2222'),
('วิชัย', 'รักเรียน', 'wichai@email.com', '082-333-3333', 'วิมล รักเรียน', '082-444-4444'),
('ปานชนก', 'สุขสันต์', 'panchanok@email.com', '083-555-5555', 'ปานตา สุขสันต์', '083-666-6666');

-- เพิ่มข้อมูลตัวอย่างในตาราง contracts
INSERT INTO contracts (room_id, tenant_id, start_date, end_date, deposit_amount, monthly_rent, status) VALUES
(2, 1, '2024-01-01', '2024-12-31', 9000.00, 4500.00, 'active'),
(5, 2, '2024-02-01', '2024-07-31', 8000.00, 4000.00, 'active');

-- เพิ่มข้อมูลตัวอย่างในตาราง payments
INSERT INTO payments (contract_id, amount, payment_date, payment_type, payment_status, notes) VALUES
(1, 4500.00, '2024-01-01', 'เงินสด', 'ชำระแล้ว', 'ค่าเช่าเดือนมกราคม'),
(1, 4500.00, '2024-02-01', 'โอนเงิน', 'ชำระแล้ว', 'ค่าเช่าเดือนกุมภาพันธ์'),
(2, 4000.00, '2024-02-01', 'เงินสด', 'ชำระแล้ว', 'ค่าเช่าเดือนกุมภาพันธ์');

-- เพิ่มข้อมูลตัวอย่างในตาราง maintenance_requests
INSERT INTO maintenance_requests (room_id, tenant_id, issue_description, priority, status, notes) VALUES
(2, 1, 'เครื่องปรับอากาศไม่เย็น', 'high', 'pending', 'รอช่างเข้าตรวจสอบ'),
(5, 2, 'ก๊อกน้ำห้องน้ำรั่ว', 'normal', 'resolved', 'ซ่อมแซมเรียบร้อยแล้ว'),
(2, 1, 'หลอดไฟในห้องดับ', 'low', 'pending', 'จะเปลี่ยนหลอดไฟใหม่'); 