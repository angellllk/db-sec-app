package core

import "testing"

func Test(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	_, err = db.Exec(`
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
				CREATE TABLE buyers (buyer_nr INTEGER NOT NULL,
					phone_nr INTEGER NOT NULL,
					name VARCHAR(255) NOT NULL,
					address VARCHAR(255) NOT NULL,
					PRIMARY KEY (buyer_nr)
				)
			');

			create_table_if_doesnt_exist('
				CREATE TABLE cars (car_code INTEGER NOT NULL,
					price INTEGER NOT NULL,
					type VARCHAR(255),
					PRIMARY KEY (car_code)
				)
			');

			create_table_if_doesnt_exist('
				CREATE TABLE shops (shop_nr INTEGER NOT NULL,
					address VARCHAR(255),
					name VARCHAR(255),
					PRIMARY KEY (shop_nr)
				)
			');

			create_table_if_doesnt_exist('
				CREATE TABLE stock (car_code INTEGER NOT NULL,
					shop_nr INTEGER NOT NULL,
					quantity INTEGER NOT NULL,
					PRIMARY KEY (car_code, shop_nr)
				)
			');

			create_table_if_doesnt_exist('
				CREATE TABLE sales_details (car_code INTEGER NOT NULL,
					sales_nr INTEGER NOT NULL,
					quantity INTEGER NOT NULL,
					PRIMARY KEY (car_code, sales_nr)
				)
			');

			create_table_if_doesnt_exist('
				CREATE TABLE sales (sales_nr INTEGER NOT NULL,
					shop_nr INTEGER NOT NULL,
					buyer_nr INTEGER NOT NULL,
					sale_date DATE NOT NULL,
					PRIMARY KEY (sales_nr)
				)
			');
		END;`)

	if err != nil {
		t.Fatalf("got error: %v", err)
	}
}
