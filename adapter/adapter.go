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
	Value  string  `gorm:"-"`
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

type LabelWValue struct {
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

func (a *Adapter) SearchByLabelValue(labelName string, labelValue interface{}) ([]Image, error) {
	rows, err := a.db.Raw("SELECT images.id FROM labels "+
		"INNER JOIN image_labels "+
		"ON labels.id = image_labels.label_id "+
		"INNER JOIN images "+
		"ON images.id = image_labels.image_id "+
		"WHERE name = ? "+
		"AND image_labels.value = ?", labelName, labelValue).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}

	var images []Image
	a.db.Where(ids).Preload("Labels").Find(&images)

	// Set Label Value
	// FIXME This is N+1 query
	for _, image := range images {
		labelRows, err := a.db.Raw("SELECT name, image_labels.value FROM labels "+
			"INNER JOIN image_labels "+
			"ON labels.id = image_labels.label_id "+
			"INNER JOIN images "+
			"ON images.id = image_labels.image_id "+
			"WHERE image_id = ?", image.ID).Rows()

		if err != nil {
			return nil, err
		}

		for labelRows.Next() {
			var name string
			var value string
			labelRows.Scan(&name, &value)
			for i, label := range image.Labels {
				if label.Name == name {
					image.Labels[i].Value = value
				}
			}
		}

		labelRows.Close()
	}

	return images, nil
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
