package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
    "strconv"
    "io/ioutil"
	"github.com/julienschmidt/httprouter"
	"github.com/taniachanda86/Assignment3/uber"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	 // "sort"
)


// LocationController represents the controller for operating on the InputAddress resource
type LocationController struct {
		session *mgo.Session
	}


type InputAddress struct {
		Name   string        `json:"name"`
		Address string 		`json:"address"`
		City string			`json:"city"`
		State string		`json:"state"`
		Zip string			`json:"zip"`
	}



type OutputAddress struct {

		Id     bson.ObjectId `json:"_id" bson:"_id,omitempty"`
		Name   string        `json:"name"`
		Address string 		`json:"address"`
		City string			`json:"city" `
		State string		`json:"state"`
		Zip string			`json:"zip"`

		Coordinate struct{
			Lat string 		`json:"lat"`
			Lang string 	`json:"lang"`
		}
	}

//------The total structure for google response--------------------------

type GoogleResponse struct {
	Results []GoogleResult
}

type GoogleResult struct {

	Address      string               `json:"formatted_address"`
	AddressParts []GoogleAddressPart `json:"address_components"`
	Geometry     Geometry
	Types        []string
}

type GoogleAddressPart struct {

	Name      string `json:"long_name"`
	ShortName string `json:"short_name"`
	Types     []string
}

type Geometry struct {

	Bounds   Bounds
	Location Point
	Type     string
	Viewport Bounds
}
type Bounds struct {
	NorthEast, SouthWest Point
}

type Point struct {
	Lat float64
	Lng float64
}
//-------------------------//---------------------//----------------------------------
//---------Adding model for Trip planner----------------//

type TripPostInput struct{
	Starting_from_location_id   string    `json:"starting_from_location_id"`
	Location_ids []string
}

type TripPostOutput struct{
	Id     bson.ObjectId 				  `json:"_id" bson:"_id,omitempty"`
	Status string  						  `json:"status"`
	Starting_from_location_id   string    `json:"starting_from_location_id"`
	Best_route_location_ids []string
	Total_uber_costs int			  `json:"total_uber_costs"`
	Total_uber_duration int			  `json:"total_uber_duration"`
	Total_distance float64				  `json:"total_distance"`

}

type UberOutput struct{
	Cost int
	Duration int
	Distance float64
}

type TripPutOutput struct{
	Id     bson.ObjectId 				  `json:"_id" bson:"_id,omitempty"`
	Status string  						  `json:"status"`
	Starting_from_location_id   string    `json:"starting_from_location_id"`
	Next_destination_location_id   string    `json:"next_destination_location_id"`
	Best_route_location_ids []string
	Total_uber_costs int			  `json:"total_uber_costs"`
	Total_uber_duration int			  `json:"total_uber_duration"`
	Total_distance float64			  `json:"total_distance"`
	Uber_wait_time_eta int 			  `json:"uber_wait_time_eta"`

}

type Struct_for_put struct{
	trip_route []string
	trip_visits map[string]int
}

type Final_struct struct{
	theMap map[string]Struct_for_put
}

//------------Ending trip planner model-----------//
// NewLocationController provides a reference to a LocationController with provided mongo session
func NewLocationController(s *mgo.Session) *LocationController {
	return &LocationController{s}
}

//The func to find google's response-----------------------------------------------
func getGoogLocation(address string) OutputAddress{
	client := &http.Client{}

	reqURL := "http://maps.google.com/maps/api/geocode/json?address="
	reqURL += url.QueryEscape(address)
	reqURL += "&sensor=false";
	fmt.Println("URL formed: "+ reqURL)
	req, err := http.NewRequest("GET", reqURL , nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in sending req to google: ", err);	
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error in reading response: ", err);	
	}

	var res GoogleResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("error in unmashalling response: ", err);	
	}

	var ret OutputAddress
	ret.Coordinate.Lat = strconv.FormatFloat(res.Results[0].Geometry.Location.Lat,'f',7,64)
	ret.Coordinate.Lang = strconv.FormatFloat(res.Results[0].Geometry.Location.Lng,'f',7,64)

	return ret;
}

//-----------------------------------------------------------------------------------


// GetLocation retrieves an individual location resource
func (uc LocationController) GetLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("location_id")
	// fmt.Println(id)
	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }

    // Grab id
    oid := bson.ObjectIdHex(id)
	var o OutputAddress
	if err := uc.session.DB("go_273").C("Locations").FindId(oid).One(&o); err != nil {
        w.WriteHeader(404)
        return
    }
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(o)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}


// GetTrip retrieves an individual trip resource
func (uc LocationController) GetTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("trip_id")
	// fmt.Println(id)
	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }

    // Grab id
    oid := bson.ObjectIdHex(id)
	var tO TripPostOutput
	if err := uc.session.DB("go_273").C("Trips").FindId(oid).One(&tO); err != nil {
        w.WriteHeader(404)
        return
    }
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(tO)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}


// CreateLocation creates a new Location resource
func (uc LocationController) CreateLocation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var u InputAddress
	var oA OutputAddress

	json.NewDecoder(r.Body).Decode(&u)	
//Trying to get the lat lang!!!--------------------
	googResCoor := getGoogLocation(u.Address + "+" + u.City + "+" + u.State + "+" + u.Zip);
    fmt.Println("resp is: ", googResCoor.Coordinate.Lat, googResCoor.Coordinate.Lang);
	
	// oA.Id = bson.NewObjectId()
	oA.Name = u.Name
	oA.Address = u.Address
	oA.City= u.City
	oA.State= u.State
	oA.Zip = u.Zip
	oA.Coordinate.Lat = googResCoor.Coordinate.Lat
	oA.Coordinate.Lang = googResCoor.Coordinate.Lang

	// Write the user to mongo
	uc.session.DB("go_273").C("Locations").Insert(oA)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(oA)
	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}


// CreateTrip creates a new Trip 
func (uc LocationController) CreateTrip(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var tI TripPostInput
	var tO TripPostOutput
	// var cost_map map[string]int
	var cost_array []int
	var duration_array []int
	var distance_array []float64
	cost_total := 0
	duration_total := 0
	distance_total := 0.0
	
	json.NewDecoder(r.Body).Decode(&tI)	

	starting_id:= bson.ObjectIdHex(tI.Starting_from_location_id)
	var start OutputAddress
	if err := uc.session.DB("go_273").C("Locations").FindId(starting_id).One(&start); err != nil {
       	w.WriteHeader(404)
        return
    }
    start_Lat := start.Coordinate.Lat
    start_Lang := start.Coordinate.Lang
    // Location_ids := tI.Location_ids

    for len(tI.Location_ids)>0{
	
			for _, loc := range tI.Location_ids{
				// var cost_array []int
				id := bson.ObjectIdHex(loc)
				var o OutputAddress
				if err := uc.session.DB("go_273").C("Locations").FindId(id).One(&o); err != nil {
		       		w.WriteHeader(404)
		        	return
		    	}
		    	loc_Lat := o.Coordinate.Lat
		    	loc_Lang := o.Coordinate.Lang
		    	
		    	getUberResponse := uber.Get_uber_price(start_Lat, start_Lang, loc_Lat, loc_Lang)
		    	fmt.Println("Uber Response is: ", getUberResponse.Cost, getUberResponse.Duration, getUberResponse.Distance );
		    	cost_array = append(cost_array, getUberResponse.Cost)
		    	duration_array = append(duration_array, getUberResponse.Duration)
		    	distance_array = append(distance_array, getUberResponse.Distance)
		    	
			}
			fmt.Println("Cost Array", cost_array)

			min_cost:= cost_array[0]
			var indexNeeded int
			for index, value := range cost_array {
		        if value < min_cost {
		            min_cost = value // found another smaller value, replace previous value in min
		            indexNeeded = index
		        }
		    }
			// fmt.Println("Min Cost", min_cost)
			// // fmt.Println(indexNeeded)
			// // fmt.Println(tI.Location_ids[indexNeeded])
			// fmt.Println("Best", tO.Best_route_location_ids)

			cost_total += min_cost
			duration_total += duration_array[indexNeeded]
			distance_total += distance_array[indexNeeded]

			tO.Best_route_location_ids = append(tO.Best_route_location_ids, tI.Location_ids[indexNeeded])
			// fmt.Println("Best", tO.Best_route_location_ids)

			starting_id = bson.ObjectIdHex(tI.Location_ids[indexNeeded])
			if err := uc.session.DB("go_273").C("Locations").FindId(starting_id).One(&start); err != nil {
       			w.WriteHeader(404)
        		return
    		}
    		tI.Location_ids = append(tI.Location_ids[:indexNeeded], tI.Location_ids[indexNeeded+1:]...)
			// fmt.Println("Af Location ids", tI.Location_ids)

    		start_Lat = start.Coordinate.Lat
    		start_Lang = start.Coordinate.Lang

    		// Re-initializing the arrays------
    		cost_array = cost_array[:0]
    		duration_array = duration_array[:0]
    		distance_array = distance_array[:0]
    		// fmt.Println("Cost Array", cost_array)

	}


	Last_loc_id := bson.ObjectIdHex(tO.Best_route_location_ids[len(tO.Best_route_location_ids)-1])
	var o2 OutputAddress
	if err := uc.session.DB("go_273").C("Locations").FindId(Last_loc_id).One(&o2); err != nil {
		w.WriteHeader(404)
		return
	}
	last_loc_Lat := o2.Coordinate.Lat
	last_loc_Lang := o2.Coordinate.Lang

	ending_id:= bson.ObjectIdHex(tI.Starting_from_location_id)
	var end OutputAddress
	if err := uc.session.DB("go_273").C("Locations").FindId(ending_id).One(&end); err != nil {
       	w.WriteHeader(404)
        return
    }
    end_Lat := end.Coordinate.Lat
    end_Lang := end.Coordinate.Lang
		    	
	getUberResponse_last := uber.Get_uber_price(last_loc_Lat, last_loc_Lang, end_Lat, end_Lang)


	tO.Id = bson.NewObjectId()
	tO.Status = "planning"
	tO.Starting_from_location_id = tI.Starting_from_location_id
	tO.Total_uber_costs = cost_total + getUberResponse_last.Cost
	tO.Total_distance = distance_total + getUberResponse_last.Distance
	tO.Total_uber_duration = duration_total + getUberResponse_last.Duration
	

	// Write the user to mongo
	uc.session.DB("go_273").C("Trips").Insert(tO)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(tO)
	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

type Internal_data struct{
	Id string               `json:"_id" bson:"_id,omitempty"`
	Trip_visited []string  `json:"trip_visited"`
	Trip_not_visited []string  `json:"trip_not_visited"`
	Trip_completed int        `json:"trip_completed"`
}

// func giveLatLang() {
	
// }

//UpdateTrip updates an existing location resource
func (uc LocationController) UpdateTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var theStruct Struct_for_put
	var final Final_struct
	final.theMap = make(map[string]Struct_for_put)

	var tPO TripPutOutput
	var internal Internal_data

	id := p[0].Value
	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }
    oid := bson.ObjectIdHex(id)
	if err := uc.session.DB("go_273").C("Trips").FindId(oid).One(&tPO); err != nil {
        w.WriteHeader(404)
        return
    }


	theStruct.trip_route = tPO.Best_route_location_ids
    theStruct.trip_route = append([]string{tPO.Starting_from_location_id}, theStruct.trip_route...)
    fmt.Println("The route array is: ", theStruct.trip_route)
    theStruct.trip_visits = make(map[string]int)

    // theStruct.trip_route = list_location_ids
    var trip_visited []string 
    var trip_not_visited []string

  	if err := uc.session.DB("go_273").C("Trip_internal_data").FindId(id).One(&internal); err != nil {
    	for index, loc := range theStruct.trip_route{
    		if index == 0{
    		// fmt.Println("Coming here.....................")
    			theStruct.trip_visits[loc] = 1
    			trip_visited = append(trip_visited, loc)
    		}else{
    			theStruct.trip_visits[loc] = 0
    			trip_not_visited = append(trip_not_visited, loc)
    		}
    	}
    	internal.Id = id
    	internal.Trip_visited = trip_visited
    	internal.Trip_not_visited = trip_not_visited
    	internal.Trip_completed = 0
    	uc.session.DB("go_273").C("Trip_internal_data").Insert(internal)

    }else {
    	for _, loc_id := range internal.Trip_visited {
    		theStruct.trip_visits[loc_id] = 1
    	}
    	for _, loc_id := range internal.Trip_not_visited {
    		theStruct.trip_visits[loc_id] = 0
    	}
    }


  	fmt.Println("Trip visit map ", theStruct.trip_visits)
  	final.theMap[id] = theStruct


  	last_index := len(theStruct.trip_route) - 1
  	trip_completed := internal.Trip_completed
  	// last_elem = theStruct.trip_route[last_index]
  		// fmt.Println("Trip completed ==", trip_completed)
  	if trip_completed == 1 {
  		fmt.Println("Entering the trip completed if statement")
  		// tpost.Status = "completed"
  		tPO.Status = "completed"

		uj, _ := json.Marshal(tPO)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		fmt.Fprintf(w, "%s", uj)
		return
	}

	for i, location := range theStruct.trip_route{
	  	if  (theStruct.trip_visits[location] == 0){
	  		tPO.Next_destination_location_id = location
	  		nextoid := bson.ObjectIdHex(location)
			var o OutputAddress
			if err := uc.session.DB("go_273").C("Locations").FindId(nextoid).One(&o); err != nil {
        		w.WriteHeader(404)
        		return
    		}
    		nlat := o.Coordinate.Lat
    		nlang:= o.Coordinate.Lang

	  		if i == 0 {
	  			starting_point := theStruct.trip_route[last_index]
	  			startingoid := bson.ObjectIdHex(starting_point)
				var o OutputAddress
				if err := uc.session.DB("go_273").C("Locations").FindId(startingoid).One(&o); err != nil {
        			w.WriteHeader(404)
        			return
    			}
    			slat := o.Coordinate.Lat
    			slang:= o.Coordinate.Lang


	  			eta := uber.Get_uber_eta(slat, slang, nlat, nlang)
	  			tPO.Uber_wait_time_eta = eta
	  			trip_completed = 1
	  		}else {
	  			starting_point2 := theStruct.trip_route[i-1]
	  			startingoid2 := bson.ObjectIdHex(starting_point2)
				var o OutputAddress
				if err := uc.session.DB("go_273").C("Locations").FindId(startingoid2).One(&o); err != nil {
        			w.WriteHeader(404)
        			return
    			}
    			slat := o.Coordinate.Lat
    			slang:= o.Coordinate.Lang
	  			eta := uber.Get_uber_eta(slat, slang, nlat, nlang)
	  			tPO.Uber_wait_time_eta = eta
	  		}	

	  		fmt.Println("Starting Location: ", tPO.Starting_from_location_id)
	  		fmt.Println("Next destination: ", tPO.Next_destination_location_id)
	  		theStruct.trip_visits[location] = 1
	  		if i == last_index {
	  			theStruct.trip_visits[theStruct.trip_route[0]] = 0
	  		}
	  		break
	  	}
	}

	// fmt.Println("After break.......")
	trip_visited  = trip_visited[:0]
	trip_not_visited  = trip_not_visited[:0]
	for location, visit := range theStruct.trip_visits{
		if visit == 1 {
			trip_visited = append(trip_visited, location)
		}else {
			trip_not_visited = append(trip_not_visited, location)
		} 
	}

	internal.Id = id
	internal.Trip_visited = trip_visited
	internal.Trip_not_visited = trip_not_visited
	fmt.Println("Trip Visisted", internal.Trip_visited)
	fmt.Println("Trip Not Visisted", internal.Trip_not_visited)
	internal.Trip_completed = trip_completed

	c := uc.session.DB("go_273").C("Trip_internal_data")
	id2 := bson.M{"_id": id}
	err := c.Update(id2, internal)
	if err != nil {
		panic(err)
	}

    uj, _ := json.Marshal(tPO)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)

}

// RemoveLocation removes an existing location resource
func (uc LocationController) RemoveLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("location_id")
	// fmt.Println(id)

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}
	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err := uc.session.DB("go_273").C("Locations").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}

//UpdateLocation updates an existing location resource
func (uc LocationController) UpdateLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var i InputAddress
	var o OutputAddress

	id := p.ByName("location_id")
	// fmt.Println(id)
	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }
    oid := bson.ObjectIdHex(id)
	
	if err := uc.session.DB("go_273").C("Locations").FindId(oid).One(&o); err != nil {
        w.WriteHeader(404)
        return
    }	

	json.NewDecoder(r.Body).Decode(&i)	
    //Trying to get the lat lang!!!--------------------
	googResCoor := getGoogLocation(i.Address + "+" + i.City + "+" + i.State + "+" + i.Zip);
    fmt.Println("resp is: ", googResCoor.Coordinate.Lat, googResCoor.Coordinate.Lang);

	
	o.Address = i.Address
	o.City = i.City
	o.State = i.State
	o.Zip = i.Zip
	o.Coordinate.Lat = googResCoor.Coordinate.Lat
	o.Coordinate.Lang = googResCoor.Coordinate.Lang

	// Write the user to mongo
	c := uc.session.DB("go_273").C("Locations")
	
	id2 := bson.M{"_id": oid}
	err := c.Update(id2, o)
	if err != nil {
		panic(err)
	}
	
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(o)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}