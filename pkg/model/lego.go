package model

// LegoSet represents a Lego set metadata.
type LegoSet struct {
	Name    string `json:"name" db:"name"`
	Model   int32  `json:"model" db:"model"`
	Catalog string `json:"catalog" db:"catalog"`
}
