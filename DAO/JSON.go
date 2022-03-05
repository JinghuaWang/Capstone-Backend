package DAO

import (
	"database/sql/driver"
	"encoding/json"
)

type Strings []string

func (c Strings) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *Strings) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}
