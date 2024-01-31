package core

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

func PrepareDB(db *sql.DB) error {
	errProc := CreateProcedures(db)
	if errProc != nil {
		return errProc
	}

	db.Close()

	connStr := go_ora.BuildUrl("localhost", 1521, "orclpdb", "angel", "angel", nil)
	db, err := sql.Open("oracle", connStr)
	if err != nil {
		return err
	}

	errCreate := CreateTables(db)
	if errCreate != nil {
		return err
	}

	return nil
}

func CreateProcedures(db *sql.DB) error {
	_, errC := db.Exec(
		`CREATE OR REPLACE PROCEDURE create_table_if_doesnt_exist(
			p_table_name VARCHAR2,
			create_table_query VARCHAR2
		) AUTHID CURRENT_USER IS
			n NUMBER;
		BEGIN
			SELECT COUNT(*) INTO n FROM user_tables WHERE table_name = UPPER(p_table_name);
			IF (n = 0) THEN
				EXECUTE IMMEDIATE create_table_query;
			END IF;
		END;`)
	if errC != nil {
		log.Println(errC)
		return errC
	}

	_, errC = db.Exec(
		`CREATE or REPLACE PROCEDURE encrypt1(text IN VARCHAR2, encrypted_text OUT VARCHAR2) AS
			key VARCHAR2(8) := '12345678';
			raw_text RAW(100);
			raw_key RAW(100);
			op_mode PLS_INTEGER;
	BEGIN
		raw_text := utl_i18n.string_to_raw(text, 'AL32UTF8');
		raw_key := utl_i18n.string_to_raw(key, 'AL32UTF8');

		op_mode := dbms_crypto.encrypt_des + dbms_crypto.pad_zero + dbms_crypto.chain_ecb;

		encrypted_text := dbms_crypto.encrypt(raw_text, op_mode, raw_key);
		dbms_output.put_line('Encryption result: ' || encrypted_text);
	END;
	`)
	if errC != nil {
		log.Println(errC)
		return errC
	}

	_, errC = db.Exec(
		`CREATE or REPLACE PROCEDURE decrypt1(encrypted_text IN VARCHAR2, decrypted_text OUT VARCHAR2) AS
		key VARCHAR2(8) := '12345678';
		--raw_text RAW(100);
		raw_key RAW(100);
		op_mode PLS_INTEGER;
	BEGIN
		raw_key := utl_i18n.string_to_raw(key, 'AL32UTF8');
		op_mode := dbms_crypto.encrypt_des + dbms_crypto.pad_zero + dbms_crypto.chain_ecb;

		decrypted_text := utl_i18n.raw_to_char(dbms_crypto.decrypt(encrypted_text, op_mode, raw_key), 'AL32UTF8');
		dbms_output.put_line('Decryption result: ' || decrypted_text);
	END;`)
	if errC != nil {
		log.Println(errC)
		return errC
	}

	return nil
}

func CreateTables(db *sql.DB) error {
	_, errC := db.Exec(`
		BEGIN
			create_table_if_doesnt_exist('buyers', '
				CREATE TABLE buyers (buyer_nr NUMBER GENERATED AS IDENTITY,
					phone_nr INTEGER NOT NULL,
					name VARCHAR2(255) NOT NULL,
					address VARCHAR2(255) NOT NULL,
					PRIMARY KEY (buyer_nr)
				)
			');

			create_table_if_doesnt_exist('cars', '
				CREATE TABLE cars (car_code NUMBER GENERATED AS IDENTITY,
					price INTEGER NOT NULL,
					type VARCHAR2(255),
					PRIMARY KEY (car_code)
				)
			');
			
			create_table_if_doesnt_exist('shops', '
				CREATE TABLE shops (shop_nr NUMBER GENERATED AS IDENTITY,
					address VARCHAR2(255),
					name VARCHAR2(255),
					PRIMARY KEY (shop_nr)
				)
			');

			create_table_if_doesnt_exist('stock', '
				CREATE TABLE stock (car_code NUMBER GENERATED AS IDENTITY,
					shop_nr INTEGER,
					quantity INTEGER NOT NULL,
					PRIMARY KEY (car_code, shop_nr),
					CONSTRAINT cars_stock
						FOREIGN KEY (car_code) REFERENCES cars (car_code),
					CONSTRAINT shops_stock
						FOREIGN KEY (shop_nr) REFERENCES shops (shop_nr)
				)
			');

			create_table_if_doesnt_exist('sales_details', '
				CREATE TABLE sales_details (car_code INTEGER,
					sales_nr INTEGER,
					quantity INTEGER NOT NULL,
					PRIMARY KEY (car_code, sales_nr),
					CONSTRAINT sales_sales_details
						FOREIGN KEY (sales_nr) REFERENCES sales (sales_nr),
					CONSTRAINT cars_sales_details
						FOREIGN KEY (car_code) REFERENCES cars (car_code)
				)
			');

			create_table_if_doesnt_exist('sales', '
				CREATE TABLE sales (sales_nr NUMBER GENERATED AS IDENTITY,
					shop_nr INTEGER NOT NULL,
					buyer_nr INTEGER NOT NULL,
					sale_date DATE NOT NULL,
					PRIMARY KEY (sales_nr),
					CONSTRAINT buyers_sales
						FOREIGN KEY (buyer_nr) REFERENCES buyers (buyer_nr),
					CONSTRAINT shops_sales
						FOREIGN KEY (shop_nr) REFERENCES shops (shop_nr)
				)
			');
		END;`)

	if errC != nil {
		log.Println(errC)
		return errC
	}

	return nil
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
