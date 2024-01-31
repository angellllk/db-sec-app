package core

import (
	"database/sql"
	"errors"
	"fmt"
)

func ParseData(data RegisterBuyer) (Buyer, error) {
	var record Buyer

	if len(data.Name) == 0 {
		return Buyer{}, errors.New("invalid name")
	}

	if data.PhoneNr == 0 {
		return Buyer{}, errors.New("invalid phone number")
	}

	if len(data.Address) == 0 {
		return Buyer{}, errors.New("invalid address")
	}

	record.Name = data.Name
	record.Address = data.Address
	record.PhoneNr = data.PhoneNr

	return record, nil
}

func AddBuyer(db *sql.DB, record Buyer) error {
	addStmt := fmt.Sprintf(`INSERT INTO buyers (phone_nr, name, address) 
	VALUES('%d', '%s', '%s')`, record.PhoneNr, record.Name, record.Address)

	_, errAdd := db.Exec(addStmt)
	if errAdd != nil {
		return errAdd
	}

	return nil
}
