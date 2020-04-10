package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Car struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Ref     string             `bson:"ref"`
	Title   string             `bson:"title"`
	Brand   string             `bson:"brand"`
	Model   string             `bson:"model"`
	Images  []Image            `bson:"images"`
	URL     string             `bson:"url"`
	Price   float64            `bson:"price"`
	Specs   string             `bson:"specs"`
	Options string             `bson:"options"`
}

// path: div.pics > div.thumbnails > ul > li > img.fiche-vignette > src
type Image struct {
	Thumbnail string `bson:"thumbnail"`
	Original  string `bson:"original"`
}
