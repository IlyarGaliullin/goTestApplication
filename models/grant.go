package models

type Grant struct {
	Table  string
	Read   bool
	Create bool
	Update bool
	Delete bool
}
