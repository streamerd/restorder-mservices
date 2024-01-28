package place

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"encore.dev/pubsub"
	"encore.dev/storage/sqldb"
)

type Menu struct {
	ID       int       `json:"_id,omitempty" bson:"_id,omitempty"`
	Place    Place     `json:"place" bson:"place"`
	Products []Product `json:"products"`
}

type Place struct {
	ID      int    `json:"id,omitempty" bson:"_id,omitempty"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Address struct {
	ID          int    `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID      int    `json:"userID" bson:"userID"`
	Default     bool   `json:"default"`
	Street      string `json:"street"`
	Number      string `json:"number"`
	ZipCode     string `json:"zipCode"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Description string `json:"description"`
}

type Product struct {
	ID           int      `json:"_id,omitempty" bson:"_id,omitempty"`
	Place        int      `json:"place" bson:"place"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Price        float32  `json:"price"`
	Images       []string `json:"images"`
	Categories   []string `json:"categories"`
	ImpactRating float32  `json:"impact_rating"`
	TasteRating  float32  `json:"taste_rating"`
	ImpactUrl    string   `json:"impact_url"`
}

// AddPlace are the parameters for adding a place to be listed.
type AddPlace struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Add adds a new place to the list of registered places.
//
//encore:api public method=POST path=/place
func (s *Service) Add(ctx context.Context, p *AddPlace) (*Place, error) {
	place := &Place{Name: p.Name, Address: p.Address}
	if err := s.db.Create(place).Error; err != nil {
		return nil, err
	}
	if _, err := SiteAddedTopic.Publish(ctx, place); err != nil {
		return nil, err
	}
	println("===========")
	print(place.ID)
	println("===========")
	return place, nil
}

type ListResponse struct {
	Places []*Place `json:"place"`
}

// Get gets a place by id.
//
//encore:api public method=GET path=/place/:placeID
func (s *Service) Get(ctx context.Context, placeID int) (*Place, error) {
	var place Place
	if err := s.db.Where("id = $1", placeID).First(&place).Error; err != nil {
		return nil, err
	}
	return &place, nil
}

// Delete deletes a place by id.
//
//encore:api public method=DELETE path=/place/:placeID
func (s *Service) Delete(ctx context.Context, placeID int) error {
	return s.db.Delete(&Place{ID: placeID}).Error
}

// List lists the places.
//
//encore:api public method=GET path=/place
func (s *Service) List(ctx context.Context) (*ListResponse, error) {
	var places []*Place
	if err := s.db.Find(&places).Error; err != nil {
		return nil, err
	}
	return &ListResponse{Places: places}, nil
}

//encore:service
type Service struct {
	db *gorm.DB
}

var placeDB = sqldb.Named("place").Stdlib()

func initService() (*Service, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: placeDB,
	}))
	if err != nil {
		return nil, err
	}
	return &Service{db: db}, nil
}

// topic
var SiteAddedTopic = pubsub.NewTopic[*Place]("place-added", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})
