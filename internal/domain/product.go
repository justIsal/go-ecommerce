package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model             
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Price       float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int        `gorm:"not null" json:"stock"`
	ImageURL    string     `gorm:"type:varchar(255)" json:"image_url"`
}