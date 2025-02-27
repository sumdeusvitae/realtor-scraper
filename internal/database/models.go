// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"
)

type Address struct {
	ID         int64
	Line       sql.NullString
	Line2      sql.NullString
	City       sql.NullString
	Country    sql.NullString
	PostalCode sql.NullString
	State      sql.NullString
	StateCode  sql.NullString
}

type Agent struct {
	ID                   string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	ProfileUrl           sql.NullString
	FirstName            sql.NullString
	LastName             sql.NullString
	NickName             sql.NullString
	PersonName           sql.NullString
	Title                sql.NullString
	Slogan               sql.NullString
	Email                sql.NullString
	Description          sql.NullString
	Video                sql.NullString
	Photo                sql.NullString
	Website              sql.NullString
	AgentRating          sql.NullInt64
	RecommendationsCount sql.NullInt64
	ReviewCount          sql.NullInt64
	FirstMonth           sql.NullInt64
	FirstYear            sql.NullInt64
	LastUpdated          sql.NullTime
	AddressID            sql.NullInt64
	BrokerID             sql.NullInt64
	OfficeID             sql.NullInt64
}

type AgentDesignation struct {
	AgentID       sql.NullString
	DesignationID sql.NullInt64
}

type AgentFeedLicense struct {
	AgentID       sql.NullString
	FeedLicenseID sql.NullInt64
}

type AgentLanguage struct {
	AgentID    sql.NullString
	LanguageID sql.NullInt64
}

type AgentMarketingArea struct {
	AgentID sql.NullString
	AreaID  sql.NullInt64
}

type AgentMultipleListingService struct {
	AgentID                  sql.NullString
	MultipleListingServiceID sql.NullInt64
}

type AgentPhone struct {
	AgentID sql.NullString
	PhoneID sql.NullInt64
}

type AgentServedArea struct {
	AgentID sql.NullString
	AreaID  sql.NullInt64
}

type AgentSpecialization struct {
	AgentID          sql.NullString
	SpecializationID sql.NullInt64
}

type AgentUserLanguage struct {
	AgentID    sql.NullString
	LanguageID sql.NullInt64
}

type AgentZip struct {
	AgentID sql.NullString
	ZipID   sql.NullInt64
}

type Area struct {
	ID        int64
	Name      sql.NullString
	StateCode sql.NullString
}

type Broker struct {
	ID            int64
	FulfillmentID sql.NullInt64
	Name          sql.NullString
	Photo         sql.NullString
	Video         sql.NullString
}

type Designation struct {
	ID   int64
	Name sql.NullString
}

type FeedLicense struct {
	ID            int64
	Country       sql.NullString
	LicenseNumber sql.NullString
	StateCode     sql.NullString
}

type Language struct {
	ID   int64
	Name sql.NullString
}

type ListingsDatum struct {
	ID              int64
	Count           sql.NullInt64
	Min             sql.NullInt64
	Max             sql.NullInt64
	LastListingDate sql.NullTime
	AgentID         sql.NullString
	Constraint      interface{}
}

type MultipleListingService struct {
	ID               int64
	Abbreviation     sql.NullString
	InactivationDate sql.NullTime
	LicenseNumber    sql.NullString
	MemberID         sql.NullString
	Type             sql.NullString
	IsPrimary        sql.NullBool
}

type Office struct {
	ID            int64
	Name          sql.NullString
	Photo         sql.NullString
	Website       sql.NullString
	Email         sql.NullString
	Slogan        sql.NullString
	Video         sql.NullString
	FulfillmentID sql.NullInt64
	AddressID     sql.NullInt64
}

type OfficeFeedLicense struct {
	OfficeID      sql.NullInt64
	FeedLicenseID sql.NullInt64
}

type OfficePhone struct {
	OfficeID sql.NullInt64
	PhoneID  sql.NullInt64
}

type Phone struct {
	ID      int64
	Ext     sql.NullString
	Number  sql.NullString
	Type    sql.NullString
	IsValid sql.NullBool
}

type RawAgent struct {
	ID      int64
	Data    sql.NullString
	AgentID sql.NullString
}

type Request struct {
	ID             int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Offset         int64
	ResultsPerPage int64
}

type SalesDatum struct {
	ID           int64
	Count        sql.NullInt64
	Min          sql.NullInt64
	Max          sql.NullInt64
	LastSoldDate sql.NullTime
	AgentID      sql.NullString
	Constraint   interface{}
}

type SocialMedia struct {
	ID      int64
	Type    sql.NullString
	Href    sql.NullString
	AgentID sql.NullString
}

type Specialization struct {
	ID   int64
	Name sql.NullString
}

type Zip struct {
	ID      int64
	ZipCode sql.NullString
}
