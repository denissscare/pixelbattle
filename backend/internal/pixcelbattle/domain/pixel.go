package domain

import "time"

type Pixel struct {
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Color     string    `json:"color"`     // "#RRGGBB"
	Author    string    `json:"author"`    // кто изменил
	Timestamp time.Time `json:"timestamp"` // когда изменил
}
