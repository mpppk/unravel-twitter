package adapter

import "github.com/jinzhu/gorm"

type Image struct {
	gorm.Model
	Url         string
	Description string
	Labels      []Label `gorm:"many2many:imageLabel;"`
}

type Label struct {
	gorm.Model
	Name   string
	Value  string
	Images []Image `gorm:"many2many:imageLabel;"`
}
