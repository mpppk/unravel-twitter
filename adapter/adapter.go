package adapter

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Image struct {
	gorm.Model
	Url         string
	Description string
	Labels      []Label `gorm:"many2many:image_labels;"`
}

type Label struct {
	gorm.Model
	Name   string
	Images []Image `gorm:"many2many:image_labels;"`
}

type ImageLabel struct {
	ImageId uint
	LabelId uint
	Value   string
}

type NewLabel struct {
	Name  string
	Value string
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

func (a *Adapter) AddLabelsToImage(image *Image, newLabels []NewLabel) error {
	for _, newLabel := range newLabels {
		var label Label
		err := a.db.FirstOrCreate(&label, Label{Name: newLabel.Name}).Error
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

		var imageLabel ImageLabel
		a.db.Model(&ImageLabel{}).
			Where(&ImageLabel{ImageId: image.ID, LabelId: label.ID}).
			First(&imageLabel).
			Update(&ImageLabel{Value: newLabel.Value})
	}
	return nil
}

func New(debugMode bool) (*Adapter, error) {
	db, err := gorm.Open("sqlite3", "test2.db")
	if err != nil {
		return nil, err
	}
	db.LogMode(debugMode)
	err = db.AutoMigrate(&Image{}, &Label{}, &ImageLabel{}).Error
	return &Adapter{db}, err
}
