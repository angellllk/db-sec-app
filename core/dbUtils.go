package core

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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

func addConstraint(db *sql.DB, table string, constraint string, key string) error {
	stmt := fmt.Sprintf(`SELECT table_name, constraint_name, status, owner
	FROM all_constraints
	WHERE r_owner = 'ANGEL'
	AND constraint_type = 'R'
	AND r_constraint_name in
	 (
	   SELECT constraint_name from all_constraints
	   WHERE constraint_type in ('P', 'U')
	   AND table_name = '%s'
	   AND owner = 'ANGEL'
	 )
	ORDER BY table_name, constraint_name`, strings.ToUpper(table))
	log.Println(stmt)
	rows, errExist := db.Query(stmt)
	if errExist != nil {
		return errExist
	}
	defer rows.Close()

	qConst := ""
	for rows.Next() {
		err := rows.Scan(&qConst)
		if err != nil {
			return err
		}

		return nil
	}

	_, errF := db.Exec(fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT %s
			FOREIGN KEY (%s) REFERENCES cars (%s)`, table, constraint, key, key))
	if errF != nil {
		return errF
	}

	return nil
}

func CreateTables(db *sql.DB) error {
	_, errC := db.Exec(`
	DECLARE
		PROCEDURE create_table_if_doesnt_exist(
			p_create_table_query VARCHAR2
		) IS
		BEGIN
			EXECUTE IMMEDIATE p_create_table_query;
		EXCEPTION
			WHEN OTHERS THEN
			-- suppresses "name is already being used" exception
			IF SQLCODE = -955 THEN
				NULL; 
			END IF;
		END;
		
		BEGIN
			create_table_if_doesnt_exist('
				CREATE TABLE buyers (buyer_nr NUMBER GENERATED AS IDENTITY,
					phone_nr INTEGER NOT NULL,
					name VARCHAR2(255) NOT NULL,
					address VARCHAR2(255) NOT NULL,
					PRIMARY KEY (buyer_nr)
				)
			');

			create_table_if_doesnt_exist('
				CREATE TABLE cars (car_code NUMBER GENERATED AS IDENTITY,
					price INTEGER NOT NULL,
					type VARCHAR2(255),
					PRIMARY KEY (car_code)
				)
			');

			create_table_if_doesnt_exist('
				CREATE TABLE shops (shop_nr NUMBER GENERATED AS IDENTITY,
					address VARCHAR2(255),
					name VARCHAR2(255),
					PRIMARY KEY (shop_nr)
				)
			');

			create_table_if_doesnt_exist('
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

			create_table_if_doesnt_exist('
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

			create_table_if_doesnt_exist('
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
