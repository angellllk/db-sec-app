package db

import (
	"database/sql"
	"fmt"
	"log"

	go_ora "github.com/sijms/go-ora/v2"
)

func ConnectDB() (*sql.DB, error) {
	port := 1521
	connStr := go_ora.BuildUrl("localhost", port, "orclpdb", "angel", "angel", nil)
	db, err := sql.Open("oracle", connStr)
	if err != nil {
		fmt.Print(err)
	}

	// check for error
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// return db instance
	return db, nil
}

func CreateTables(db *sql.DB) {
	_, errC := db.Exec(`CREATE TABLE buyers (buyer_nr INTEGER NOT NULL,
		phone_nr INTEGER NOT NULL,
		name VARCHAR(255) NOT NULL,
		address VARCHAR(255) NOT NULL,
		PRIMARY KEY (buyer_nr)
	)`)
	if errC != nil {
		log.Println(errC)
	}

	_, errC = db.Exec(`CREATE TABLE cars (car_code INTEGER NOT NULL,
		price INTEGER NOT NULL,
		type VARCHAR(255),
		PRIMARY KEY (car_code)
	)`)
	if errC != nil {
		log.Println(errC)
	}

	_, errC = db.Exec(`CREATE TABLE stock (car_code INTEGER NOT NULL,
		shop_nr INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		PRIMARY KEY (car_code, shop_nr)
	)`)
	if errC != nil {
		log.Println(errC)
	}

	_, errC = db.Exec(`CREATE TABLE shops (shop_nr INTEGER NOT NULL,
		address VARCHAR(255),
		name VARCHAR(255),
		PRIMARY KEY (shop_nr)
	)`)
	if errC != nil {
		log.Println(errC)
	}

	_, errC = db.Exec(`CREATE TABLE sales_details (car_code INTEGER NOT NULL,
		sales_nr INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		PRIMARY KEY (car_code, sales_nr)
	)`)
	if errC != nil {
		log.Println(errC)
	}

	_, errC = db.Exec(`CREATE TABLE sales (sales_nr INTEGER NOT NULL,
		shop_nr INTEGER NOT NULL,
		buyer_nr INTEGER NOT NULL,
		sale_date DATE NOT NULL,
		PRIMARY KEY (sales_nr)
	)`)
	if errC != nil {
		log.Println(errC)
	}

	_, errF := db.Exec(`ALTER TABLE stock ADD CONSTRAINT cars_stock
		FOREIGN KEY (car_code) REFERENCES cars (car_code)`)
	if errF != nil {
		log.Println(errF)
	}

	_, errF = db.Exec(`ALTER TABLE stock ADD CONSTRAINT shops_stock
		FOREIGN KEY (shop_nr) REFERENCES shops (shop_nr)`)
	if errF != nil {
		log.Println(errF)
	}

	_, errF = db.Exec(`ALTER TABLE sales ADD CONSTRAINT buyers_sales
		FOREIGN KEY (buyer_nr) REFERENCES buyers (buyer_nr)`)
	if errF != nil {
		log.Println(errF)
	}

	_, errF = db.Exec(`ALTER TABLE sales ADD CONSTRAINT shops_sales
		FOREIGN KEY (shop_nr) REFERENCES shops (shop_nr)`)
	if errF != nil {
		log.Println(errF)
	}

	_, errF = db.Exec(`ALTER TABLE sales_details ADD CONSTRAINT sales_sales_details
		FOREIGN KEY (sales_nr) REFERENCES sales (sales_nr)`)
	if errF != nil {
		log.Println(errF)
	}

	_, errF = db.Exec(`ALTER TABLE sales_details ADD CONSTRAINT cars_sales_details
		FOREIGN KEY (car_code) REFERENCES cars (car_code)`)
	if errF != nil {
		log.Println(errF)
	}
}

func DeleteTables(db *sql.DB) {
	_, errD := db.Exec(`ALTER TABLE stock DROP CONSTRAINT cars_stock`)
	if errD != nil {
		log.Println(errD)
	}

	_, errD = db.Exec(`ALTER TABLE stock DROP CONSTRAINT shops_stock`)
	if errD != nil {
		log.Println(errD)
	}

	_, errD = db.Exec(`ALTER TABLE sales DROP CONSTRAINT buyers_sales`)
	if errD != nil {
		log.Println(errD)
	}

	_, errD = db.Exec(`ALTER TABLE sales DROP CONSTRAINT shops_sales`)
	if errD != nil {
		log.Println(errD)
	}

	_, errD = db.Exec(`ALTER TABLE sales_details DROP CONSTRAINT sales_sales_details`)
	if errD != nil {
		log.Println(errD)
	}

	_, errD = db.Exec(`ALTER TABLE sales_details DROP CONSTRAINT cars_sales_details`)
	if errD != nil {
		log.Println(errD)
	}

	_, errD = db.Exec(`DROP TABLE buyers`)
	if errD != nil {
		log.Println(errD)
	}

	_, errD = db.Exec(`DROP TABLE cars`)
	if errD != nil {
		log.Println(errD)
	}

	_, errD = db.Exec(`DROP TABLE stock`)
	if errD != nil {
		log.Println(errD)
	}

	_, errD = db.Exec(`DROP TABLE shops`)
	if errD != nil {
		log.Println(errD)
	}

	_, errD = db.Exec(`DROP TABLE sales_details`)
	if errD != nil {
		log.Println(errD)
	}

	_, errD = db.Exec(`DROP TABLE sales`)
	if errD != nil {
		log.Println(errD)
	}
}
