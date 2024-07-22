package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Role     string `gorm:"not null"`
}

type Category struct {
	gorm.Model
	Name         string `gorm:"uniqueIndex;not null"`
	DisplayOrder int    `gorm:"not null"`
	MenuItems    []MenuItem
}

type MenuItem struct {
	gorm.Model
	CategoryID  uint
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	ImageURL    string
	IsAvailable bool `gorm:"not null;default:true"`
	AddOns      []AddOn
}

type AddOn struct {
	gorm.Model
	MenuItemID uint
	Name       string  `gorm:"not null"`
	Price      float64 `gorm:"not null"`
}

type Order struct {
	gorm.Model
	TableNumber  string  `gorm:"not null"`
	Status       string  `gorm:"not null"`
	TotalAmount  float64 `gorm:"not null"`
	OrderDetails []OrderDetail
}

type OrderDetail struct {
	gorm.Model
	OrderID             uint
	MenuItemID          uint
	Quantity            int     `gorm:"not null"`
	UnitPrice           float64 `gorm:"not null"`
	Subtotal            float64 `gorm:"not null"`
	SpecialInstructions string
	SelectedAddOns      []SelectedAddOn
}

type SelectedAddOn struct {
	gorm.Model
	OrderDetailID uint
	AddOnID       uint
	Name          string  `gorm:"not null"`
	Price         float64 `gorm:"not null"`
}

type TimeRange struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

type DaySchedule struct {
	Ranges []TimeRange `json:"ranges"`
}

type WeekSchedule struct {
	Monday    DaySchedule `json:"monday"`
	Tuesday   DaySchedule `json:"tuesday"`
	Wednesday DaySchedule `json:"wednesday"`
	Thursday  DaySchedule `json:"thursday"`
	Friday    DaySchedule `json:"friday"`
	Saturday  DaySchedule `json:"saturday"`
	Sunday    DaySchedule `json:"sunday"`
}

type SpecialDate struct {
	Date     string      `json:"date"` // Format: "2006-01-02"
	Schedule DaySchedule `json:"schedule"`
}

type OpeningHours struct {
	WeekSchedule  WeekSchedule  `json:"week_schedule"`
	SpecialDates  []SpecialDate `json:"special_dates"`
	HolidayClosed bool          `json:"holiday_closed"`
}

type RestaurantInfo struct {
	gorm.Model
	Name         string       `json:"name" gorm:"uniqueIndex"`
	Description  string       `json:"description"`
	Address      string       `json:"address"`
	Phone        string       `json:"phone"`
	Email        string       `json:"email"`
	LogoURL      string       `json:"logo_url"`
	BannerURL    string       `json:"banner_url"`
	OpeningHours OpeningHours `gorm:"type:jsonb" json:"opening_hours"`
}

// Scan implements the sql.Scanner interface for OpeningHours
func (oh *OpeningHours) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &oh)
}

// Value implements the driver.Valuer interface for OpeningHours
func (oh OpeningHours) Value() (driver.Value, error) {
	return json.Marshal(oh)
}

// IsOpen checks if the restaurant is open at the given time
func (oh OpeningHours) IsOpen(t time.Time) bool {
	// Check for special dates first
	dateStr := t.Format("2006-01-02")
	for _, specialDate := range oh.SpecialDates {
		if specialDate.Date == dateStr {
			return isDayScheduleOpen(specialDate.Schedule, t)
		}
	}

	// Check regular week schedule
	var daySchedule DaySchedule
	switch t.Weekday() {
	case time.Monday:
		daySchedule = oh.WeekSchedule.Monday
	case time.Tuesday:
		daySchedule = oh.WeekSchedule.Tuesday
	case time.Wednesday:
		daySchedule = oh.WeekSchedule.Wednesday
	case time.Thursday:
		daySchedule = oh.WeekSchedule.Thursday
	case time.Friday:
		daySchedule = oh.WeekSchedule.Friday
	case time.Saturday:
		daySchedule = oh.WeekSchedule.Saturday
	case time.Sunday:
		daySchedule = oh.WeekSchedule.Sunday
	}

	return isDayScheduleOpen(daySchedule, t)
}

func isDayScheduleOpen(schedule DaySchedule, t time.Time) bool {
	currentTime := t.Format("15:04")
	for _, timeRange := range schedule.Ranges {
		if currentTime >= timeRange.Open && currentTime < timeRange.Close {
			return true
		}
	}
	return false
}
