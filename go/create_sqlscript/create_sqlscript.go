package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func main() {
	//create file
	filename := "mock-sql-insert.txt"

	// Check if the file already exists
	if _, err := os.Stat(filename); err == nil {
		// File exists, remove it
		err := os.Remove(filename)
		if err != nil {
			fmt.Println("Error removefile:", err)
			return
		}
		fmt.Println("Existing file removed successfully.")
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write content using a for loop
	for i := 0; i <= 999; i++ {
		text := formatMockBranch(i)
		_, err := file.WriteString(text + "\n")
		if err != nil {
			fmt.Println("Error writing to file branch: ", err)
			return
		}

		text = formatMockMerchant_x_branch(i)
		_, err = file.WriteString(text + "\n")
		if err != nil {
			fmt.Println("Error writing to file merchant_x_branch: ", err)
			return
		}
	}
	fmt.Println("Content written to the file successfully!")
}

func formatMockBranch(num int) string {
	u1 := uuid.New()
	branch_sk := fmt.Sprintf(`UNHEX(REPLACE("%s","-",""))`, u1)
	branch_id := formatDigit(num)
	name := fmt.Sprintf("NEW MOCKS 2 สาขาที่ %d", num)
	number := formatDigit(num)

	text := fmt.Sprintf(`INSERT INTO branch(branch_sk, merchant_id,branch_id, name, number, branch_province_code, branch_mobile, address_no, building, street, moo, soi, sub_district, sub_district_code, district,district_code, province, zip_code, latitude, longitude, created_date, created_by,updated_date,updated_by,is_update_location,last_location_date,tour_company_name,tour_package_name,name_en,last_branch_name_updated_date,use_same_address_for_delivery_address,use_same_address_for_tax_address,vat_branch_code) VALUES ( %s,"W202411151038589998r", "%s", "%s", "%s", "010", "omrt3cyg4eXJAm4ZcuM6Vd72a03akTbds7Q=", "77/782", "99", "ถนน2", "หมู่2", "ซอย2", "แพรกษา", "10140100", "แพรกษา2", "1014", "กรุงเทพมหานคร", "10110", "15.877076727392777", "100.99253810942888", now(), null, now(), "10007", "Y", now(), null, null, "", now(), 0, null, null)`, string(branch_sk), branch_id, name, number)
	return text
}

func formatMockMerchant_x_branch(num int) string  {
	u1 := uuid.New()
	branch_sk := fmt.Sprintf(`UNHEX(REPLACE("%s","-",""))`, u1)
	text := fmt.Sprintf(`INSERT INTO merchant_x_branch(merchant_id,branch_sk,mobile_no,created_date,created_by,updated_date,updated_by) VALUES ("W202411151038589998r", %s,"omrt3cyg4eXJAm4ZcuM6Vd72a03akTbds7Q=", now(), "10007", now(), "10007")`, branch_sk)
	return text
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
