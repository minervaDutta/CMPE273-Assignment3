package main

import (
	"net/http"
	"Assignments/assnmnt3/httprouter"
	"Assignments/assnmnt3/controllers"
	"gopkg.in/mgo.v2"
)

func main() {
	//***************************************************The RESTful calls************************************************

	r := httprouter.New()
	// make a new router
	uc := controllers.NewLocationController(getSession())
	// uc = LocationController instance
	r.GET("/locations/:location_id", uc.GetLocation)
  	// Get - location resource
	r.GET("/trips/:trip_id", uc.GetTrip)
	// Get - trip resourc
	r.POST("/locations", uc.CreateLocation)
	// Create a new address

	r.POST("/trips", uc.CreateTrip)
	// Create a new trip
	r.PUT("/locations/:location_id", uc.UpdateLocation)
	// Update an address
	r.PUT("/trips/:trip_id/request", uc.UpdateTrip)
	// Update an trip
	r.DELETE("/locations/:location_id", uc.RemoveLocation)
	// Remove an existing address
	http.ListenAndServe("localhost:8080", r)
	// Start server
}

//**********************************************************************************************************************

func getSession() *mgo.Session {
	//create a new mongo session and panics if connection error occurs
	s, err := mgo.Dial("mongodb://admin:admin@ds045464.mongolab.com:45464/go_273")
	// Connect to local mongo
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)
	return s
}
