package core

import "testing"

func TestCreateTableScript(t *testing.T) {
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

	if err != nil {
		t.Fatalf("got error: %v", err)
	}
}

func TestAESEncDecScript(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	_, err = db.Exec(`
	DECLARE
		encrypt_result VARCHAR2(200);
		decrypt_result VARCHAR2(200); 
	BEGIN
		create_table_if_doesnt_exist('crypto', '
			CREATE TABLE crypto (id NUMBER GENERATED AS IDENTITY,
				encrypt_r VARCHAR2(200),
				decrypt_r VARCHAR2(200)
			)
		');

		encrypt1('TEST123', encrypt_result);
		decrypt1(encrypt_result, decrypt_result);

		INSERT INTO crypto (encrypt_r, decrypt_r) VALUES (encrypt_result, decrypt_result);
	END;`)

	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	defer db.Close()
}
