package models

import "time"

type EventType string

const (
	CitizenReportCreated EventType = "citizen_report_created"
	CitizenReportUpdated EventType = "citizen_report_updated"
	CitizenReportDeleted EventType = "citizen_report_deleted"
)

// type ProductEvent struct {
// 	EventID     string    `json:"event_id"`
// 	EventType   EventType `json:"event_type"`
// 	Timestamp   time.Time `json:"timestamp"`
// 	ProductID   int       `json:"product_id,omitempty"`
// 	ProductData *Product  `json:"product_data,omitempty"`

// 	ProducerID string `json:"producer_id"`
// 	Sequence   int64  `json:"sequence"`
// }

type CitizenReportEvent struct {
	EventID    string         `json:"event_id"`
	EventType  EventType      `json:"event_type"`
	Timestamp  time.Time      `json:"timestamp"`
	ReportID   string         `json:"report_id,omitempty"`
	ReportData *CitizenReport `json:"report_data,omitempty"`

	ProducerID string `json:"producer_id"`
	Sequence   int64  `json:"sequence"`
}
