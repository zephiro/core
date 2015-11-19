package components

import (
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

func (component *ComponentModel) SetDI(di *Module) {

	component.di = di
}

func (component *ComponentModel) UpdatePrice(prices map[string]float64) {

	database := component.di.Mongo.Database
	set := bson.M{"store.updated_at": time.Now(), "activated": true}

	for key, price := range prices {

		set["store.prices." + key] = price
	}

	err := database.C("components").Update(bson.M{"_id": component.Id}, bson.M{"$set": set})

	if err != nil {
		panic(err)
	}

	go component.UpdateAlgolia()
}

func (component *ComponentModel) UpdateAlgolia() {

	index := component.di.Search.Get("components")

	// Compose algolia item
	full_name := component.FullName 

	if full_name == "" {
		full_name = component.Name
	}

	var image string

	if len(component.Images) > 0 {

		image = component.Images[0].Path
		image = strings.Replace(image, "full/", "", -1)
	}

	object := make(map[string]interface{})
	object["objectID"] = component.Id.Hex()
	object["name"] = component.Name
	object["full_name"] = full_name
	object["part_number"] = component.PartNumber
	object["slug"] = component.Slug
	object["image"] = image
	object["type"] = component.Type
	object["activated"] = true

	_, err := index.UpdateObject(object)

	if err != nil {
		panic(err)
	}	
}