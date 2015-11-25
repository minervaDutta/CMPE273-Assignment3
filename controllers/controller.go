package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
  	"strconv"
  	"io/ioutil"
	"Assignments/assnmnt3/httprouter"
	"Assignments/assnmnt3/uber"
)



type LocationController struct {
		session *mgo.Session
	}

// Controller for operating the InputAddress resource below


type InputAddress struct {
		Name   string   `json:"name"`
		Address string 	`json:"address"`
		City string	`json:"city"`
		State string	`json:"state"`
		Zip string	`json:"zip"`
	}



type OutputAddress struct {

		Id     bson.ObjectId `json:"_id" bson:"_id,omitempty"`
		Name   string        `json:"name"`
		Address string 	     `json:"address"`
		City string	     `json:"city" `
		State string	     `json:"state"`
		Zip string	     `json:"zip"`

		Coordinate struct{
			Lat string 	`json:"lat"`
			Lang string 	`json:"lang"`
		}
	}

//######################struct for google response##################################

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

//#####################The trip planner struct###############################

type TripPostInput struct{
	Starting_from_location_id   string    `json:"starting_from_location_id"`
	Location_ids []string
}

type TripPostOutput struct{
	Id     bson.ObjectId 				`json:"_id" bson:"_id,omitempty"`
	Status string  					`json:"status"`
	Starting_from_location_id   string    		`json:"starting_from_location_id"`
	Best_route_location_ids []string
	Total_uber_costs int			  	`json:"total_uber_costs"`
	Total_uber_duration int				`json:"total_uber_duration"`
	Total_distance float64				`json:"total_distance"`

}

type UberOutput struct{
	Cost int
	Duration int
	Distance float64
}

type TripPutOutput struct{
	Id     bson.ObjectId 			 `json:"_id" bson:"_id,omitempty"`
	Status string  				 `json:"status"`
	Starting_from_location_id   string    	 `json:"starting_from_location_id"`
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

//###################################################################################


func NewLocationController(s *mgo.Session) *LocationController {
	return &LocationController{s}
}
//  reference to a LocationController with given mongo session

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

	//The func to find google's response

	var ret OutputAddress
	ret.Coordinate.Lat = strconv.FormatFloat(res.Results[0].Geometry.Location.Lat,'f',7,64)
	ret.Coordinate.Lang = strconv.FormatFloat(res.Results[0].Geometry.Location.Lng,'f',7,64)

	return ret;
}




func (uc LocationController) GetLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// GetLocation retrieves an individual location resource
	id := p.ByName("location_id")

	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }

  oid := bson.ObjectIdHex(id)
	var o OutputAddress
	if err := uc.session.DB("go_273").C("Locations").FindId(oid).One(&o); err != nil {
        w.WriteHeader(404)
        return
    }

	uj, _ := json.Marshal(o)
	// Marshal interface into JSON structure
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}



func (uc LocationController) GetTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// GetTrip retrieves an individual trip resource
	id := p.ByName("trip_id")

	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }


  oid := bson.ObjectIdHex(id)
	var tO TripPostOutput
	if err := uc.session.DB("go_273").C("Trips").FindId(oid).One(&tO); err != nil {
        w.WriteHeader(404)
        return
    }
	uj, _ := json.Marshal(tO)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}



func (uc LocationController) CreateLocation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// CreateLocation creates a new Location resource
	var u InputAddress
	var oA OutputAddress

	json.NewDecoder(r.Body).Decode(&u)
	googResCoor := getGoogLocation(u.Address + "+" + u.City + "+" + u.State + "+" + u.Zip);
    fmt.Println("resp is: ", googResCoor.Coordinate.Lat, googResCoor.Coordinate.Lang);
	//get lat and lang
	oA.Name = u.Name
	oA.Address = u.Address
	oA.City= u.City
	oA.State= u.State
	oA.Zip = u.Zip
	oA.Coordinate.Lat = googResCoor.Coordinate.Lat
	oA.Coordinate.Lang = googResCoor.Coordinate.Lang


	uc.session.DB("go_273").C("Locations").Insert(oA)
	uj, _ := json.Marshal(oA)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}



func (uc LocationController) CreateTrip(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//create a new trip
	var tI TripPostInput
	var tO TripPostOutput
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

    for len(tI.Location_ids)>0{

			for _, loc := range tI.Location_ids{
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
		            min_cost = value // replace with smaller vals
		            indexNeeded = index
		        }
		    }

			cost_total += min_cost
			duration_total += duration_array[indexNeeded]
			distance_total += distance_array[indexNeeded]

			tO.Best_route_location_ids = append(tO.Best_route_location_ids, tI.Location_ids[indexNeeded])

			starting_id = bson.ObjectIdHex(tI.Location_ids[indexNeeded])
			if err := uc.session.DB("go_273").C("Locations").FindId(starting_id).One(&start); err != nil {
       			w.WriteHeader(404)
        		return
    		}
    		tI.Location_ids = append(tI.Location_ids[:indexNeeded], tI.Location_ids[indexNeeded+1:]...)


    		start_Lat = start.Coordinate.Lat
    		start_Lang = start.Coordinate.Lang


    		cost_array = cost_array[:0]
    		duration_array = duration_array[:0]
    		distance_array = distance_array[:0]
    		// reset arrays

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



	uc.session.DB("go_273").C("Trips").Insert(tO)


	uj, _ := json.Marshal(tO)

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




func (uc LocationController) UpdateTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//UpdateTrip updates an existing location resource
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

    var trip_visited []string
    var trip_not_visited []string

  	if err := uc.session.DB("go_273").C("Trip_internal_data").FindId(id).One(&internal); err != nil {
    	for index, loc := range theStruct.trip_route{
    		if index == 0{
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
  	if trip_completed == 1 {
  		fmt.Println("Entering the trip completed if statement")
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


func (uc LocationController) RemoveLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//  removes an existing location resource
	id := p.ByName("location_id")



	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}// Verify id is ObjectId, otherwise bail

	oid := bson.ObjectIdHex(id)


	if err := uc.session.DB("go_273").C("Locations").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}	// Remove user


	w.WriteHeader(200)
}


func (uc LocationController) UpdateLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var i InputAddress
	var o OutputAddress

	id := p.ByName("location_id")

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

	googResCoor := getGoogLocation(i.Address + "+" + i.City + "+" + i.State + "+" + i.Zip);
    fmt.Println("resp is: ", googResCoor.Coordinate.Lat, googResCoor.Coordinate.Lang);


	o.Address = i.Address
	o.City = i.City
	o.State = i.State
	o.Zip = i.Zip
	o.Coordinate.Lat = googResCoor.Coordinate.Lat
	o.Coordinate.Lang = googResCoor.Coordinate.Lang


	c := uc.session.DB("go_273").C("Locations")

	id2 := bson.M{"_id": oid}
	err := c.Update(id2, o)
	if err != nil {
		panic(err)
	}


	uj, _ := json.Marshal(o)

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}
