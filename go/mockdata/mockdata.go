package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {
	// connect to db
	// driverName := os.Getenv("DRIVER_NAME")
	driverName := "mysql"
	// dbHost := os.Getenv("DB_HOST")
	dbHost := "bob:101@tcp(localhost:3306)/HelloWorld"
	db, err := connectDB(driverName, dbHost)
	if err != nil {
		log.Fatal("connectDB: ", err)
	}

	// create table if it not exist
	err = createTableIfItNotExist(db)
	if err != nil {
		log.Fatal("create table: ", err)
	}

	// generate mockdata and insert to table
	err = insertMockItemToDB(db)
	if err != nil {
		log.Fatal("insert mock: ", err)
	}
}

func connectDB(driverName string, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTableIfItNotExist(db *sql.DB) error {
	err := checkIfTableIsExist(db)
	if err != nil {
		fmt.Println("Check if table exist Fail")
		return err
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS branch (
			branch_sk BINARY(16) NOT NULL, 
			merchant_id CHAR(20) COLLATE utf8_bin NOT NULL, 
			branch_id VARCHAR(20) COLLATE utf8_bin NOT NULL, 
			name VARCHAR(2000) COLLATE utf8_bin DEFAULT NULL, 
			number VARCHAR(25) COLLATE utf8_bin DEFAULT NULL, 
			branch_province_code VARCHAR(5) COLLATE utf8_bin DEFAULT NULL, 
			branch_mobile VARCHAR(50) COLLATE utf8_bin DEFAULT NULL, 
			address_no VARCHAR(500) COLLATE utf8_bin DEFAULT NULL, 
			building VARCHAR(500) COLLATE utf8_bin DEFAULT NULL, 
			street VARCHAR(500) COLLATE utf8_bin DEFAULT NULL, 
			moo VARCHAR(500) COLLATE utf8_bin DEFAULT NULL, 
			soi VARCHAR(500) COLLATE utf8_bin DEFAULT NULL, 
			sub_district VARCHAR(256) COLLATE utf8_bin DEFAULT NULL, 
			sub_district_code VARCHAR(20) COLLATE utf8_bin DEFAULT NULL COMMENT 'Sub district code', 
			district VARCHAR(256) COLLATE utf8_bin DEFAULT NULL, 
			district_code VARCHAR(20) COLLATE utf8_bin DEFAULT NULL COMMENT 'District code', 
			province VARCHAR(100) COLLATE utf8_bin DEFAULT NULL, 
			zip_code VARCHAR(10) COLLATE utf8_bin DEFAULT NULL, 
			latitude VARCHAR(100) COLLATE utf8_bin DEFAULT NULL, 
			longitude VARCHAR(100) COLLATE utf8_bin DEFAULT NULL, 
			created_date DATETIME DEFAULT CURRENT_TIMESTAMP, 
			created_by VARCHAR(200) COLLATE utf8_bin DEFAULT NULL, 
			updated_date DATETIME DEFAULT NULL, 
			updated_by VARCHAR(200) COLLATE utf8_bin DEFAULT NULL, 
			is_update_location CHAR(1) COLLATE utf8_bin DEFAULT NULL, 
			last_location_date DATETIME DEFAULT NULL, 
			tour_company_name VARCHAR(100) COLLATE utf8_bin DEFAULT NULL, 
			tour_package_name VARCHAR(200) COLLATE utf8_bin DEFAULT NULL, 
			name_en VARCHAR(2000) COLLATE utf8_bin DEFAULT NULL, 
			last_branch_name_updated_date DATETIME DEFAULT NULL, 
			use_same_address_for_delivery_address BIT(1) DEFAULT b'1', 
			use_same_address_for_tax_address TINYINT(1) DEFAULT NULL, 
			vat_branch_code VARCHAR(10) COLLATE utf8_bin DEFAULT NULL, 
			PRIMARY KEY (merchant_id, branch_id), 
			KEY IDX_branch_id (branch_sk), 
			KEY IDX_branch_mobile (branch_mobile)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin`

	_, err = db.Exec(createTableQuery)
	return err
}

func checkIfTableIsExist(db *sql.DB) error {
	rows, err := db.Query("SHOW TABLES LIKE 'branch'")
	if err != nil {
		return err // Handle error checking the table
	}
	defer rows.Close()

	if rows.Next() {
		// Table exists, drop it before re-creating
		dropTableQuery := "DROP TABLE branch"
		_, err := db.Exec(dropTableQuery)
		if err != nil {
			return err // Handle error dropping the table
		}
	}
	return nil
}

func insertMockItemToDB(db *sql.DB) error {
	stmt, err := db.Prepare(`INSERT INTO branch (branch_sk, merchant_id, branch_id, name, number, branch_province_code, branch_mobile, address_no, building, street, moo, soi, sub_district, sub_district_code, district,district_code, province, zip_code, latitude, longitude, created_date, created_by,updated_date,updated_by,is_update_location,last_location_date,tour_company_name,tour_package_name,name_en,last_branch_name_updated_date,use_same_address_for_delivery_address,use_same_address_for_tax_address,vat_branch_code) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		fmt.Printf("Prepare Fail")
		return err
	}
	defer stmt.Close()

	for i := 1; i < 1000; i++ {
		branch_sk := UNHEX(REPLACEALL(fmt.Sprintf("4f3c4f5b-c5b4-46f8-ba2b-e7ab0d14%s", formatDigit(i)), "-", ""))
		branch_id := formatDigit(i)
		number := formatDigit(i)
		name := fmt.Sprintf("NEW MOCKS 2 สาขาที่ %d", i)
		_, err := stmt.Exec(branch_sk, "W202411151038589998r", branch_id, name, number, "010", "omrt3cyg4eXJAm4ZcuM6Vd72a03akTbds7Q=", "77/782", "99", "ถนน2", "หมู่2", "ซอย2", "แพรกษา", "10140100", "แพรกษา2", "1014", "กรุงเทพมหานคร", "10110", "15.877076727392777", "100.99253810942888", time.Now(), nil, time.Now(), "10007", "Y", time.Now(), nil, nil, "", time.Now(), 0, nil, nil)
		if err != nil {
			fmt.Printf("Insert Fail")
			return err
		}
	}
	return nil
}

func REPLACEALL(text string, replaced string, replacewith string) string {
	newText := strings.ReplaceAll(text, replaced, replacewith)
	return newText
}

func UNHEX(hexString string) []byte {
	dhex, _ := hex.DecodeString(hexString)
	return dhex
}

func formatDigit(num int) string {
	if num < 1 || num > 9999 {
		return "Invalid input. num should be between 1 and 9999."
	}

	return fmt.Sprintf("%04d", num)
}
