package adapter

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

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

type Adapter struct {
	db *gorm.DB
}

func (a *Adapter) Close() {
	a.db.Close()
}

func (a *Adapter) AddImage(image *Image) *gorm.DB {
	return a.db.FirstOrCreate(image)
}

func (a *Adapter) AddLabelsToImage(image *Image, names []string) error {
	for _, name := range names {
		var label Label
		err := a.db.FirstOrCreate(&label, Label{Name: name}).Error
		if err != nil {
			return err
		}

		err = a.db.FirstOrCreate(image, Image{Url: image.Url, Description: image.Description}).
			First(image).
			Association("Labels").
			Append([]Label{label}).Error

		if err != nil {
			return err
		}

	}
	return nil
}

func New(debugMode bool) (*Adapter, error) {
	db, err := gorm.Open("sqlite3", "test2.db")
	if err != nil {
		return nil, err
	}
	db.LogMode(debugMode)
	err = db.AutoMigrate(&Image{}, &Label{}).Error
	return &Adapter{db}, err
}
