package main

import (
	// Standard library packages
	"net/http"

	// Third party packages
	"github.com/julienschmidt/httprouter"
	"github.com/minervaDutta/Assignment3/controllers"
	"gopkg.in/mgo.v2"
)

func main() {
	// Instantiate a new router
	r := httprouter.New()

	// Get a LocationController instance
	uc := controllers.NewLocationController(getSession())

	// Get a location resource
	r.GET("/locations/:location_id", uc.GetLocation)

	// Get a trip resource
	r.GET("/trips/:trip_id", uc.GetTrip)

	// Create a new address
	r.POST("/locations", uc.CreateLocation)

	// Create a new trip
	r.POST("/trips", uc.CreateTrip)

	// Update an address
	r.PUT("/locations/:location_id", uc.UpdateLocation)

	// Update an trip
	r.PUT("/trips/:trip_id/request", uc.UpdateTrip)


	// Remove an existing address
	r.DELETE("/locations/:location_id", uc.RemoveLocation)

	// Fire up the server
	http.ListenAndServe("localhost:8080", r)
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://minervaDutta:amma123@ds041144.mongolab.com:31144/go_273")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)
	// Deliver session
	return s
}
