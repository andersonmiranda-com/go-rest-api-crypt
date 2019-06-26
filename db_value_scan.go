package main

import (
	_ "github.com/go-sql-driver/mysql"
)

/*


---------------------------------------------------------------------------
NOT IN USE
---------------------------------------------------------------------------


type EncryptedValue string

// Implements driver.Valuer.
func (ev EncryptedValue) Value() (driver.Value, error) {

	UPK := getUserPrivateKey()
	ciphertext, err := encrypt([]byte(ev), UPK)
	if err != nil {
		log.Fatal(err)
	}
	return driver.Value(ciphertext), nil
}

// Implements sql.Scanner. Simplistic -- only handles string and []byte
func (ev *EncryptedValue) Scan(src interface{}) error {
	UPK := getUserPrivateKey()
	plaintext, err := decrypt(src.([]byte), UPK)
	if err != nil {
		log.Fatal(err)
	}
	*ev = EncryptedValue(plaintext)
	return nil
}

*/

// type User struct {
// 	ID       int64  `json:"id"`
// 	Name     string `json:"name"`
// 	LastName string `json:"lastname"`
// 	Password []byte `json:"password"`
// }
/*
type LowercaseString string
type MagicNumber int


// Implements driver.Valuer.
func (ls LowercaseString) Value() (driver.Value, error) {
	return driver.Value(strings.ToLower(string(ls))), nil
}

// Implements sql.Scanner. Simplistic -- only handles string and []byte
func (ls *LowercaseString) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New("Incompatible type for LowercaseString")
	}
	*ls = LowercaseString(strings.ToUpper(source))
	return nil
}

func (num MagicNumber) Value() (driver.Value, error) {
	return int64(num * 2), nil
}

func (num *MagicNumber) Scan(src interface{}) error {
	var source int
	source, _ = strconv.Atoi(string(src.([]byte)))
	source = source / 2
	*num = MagicNumber(source)
	return nil
}

*/
