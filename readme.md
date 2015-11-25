
##Assignment 3: CMPE-273 (Fall 15)

Author: Minerva Dutta

Description: Trip Planner

Input: Locations from a database

Output: Best route (metric: cost and duration)

External API's: Uber's price estimates

Usage:
-------
##Plan a trip 
POST        /trips   

Request:
{
   
    "starting_from_location_id: "10040",
    "location_ids" : [ "10000", "10002", "20005", "30000" ]
    
}

//Response: array (best_route_location_ids) that sorts location id's Uber cost and duration.

Response: HTTP 201
{

     "id" : "4038",
     “status” : “planning”,
     "starting_from_location_id: "10040",
     "best_route_location_ids" : [ "30000", "10000", "10002", "20005" ],
     "total_uber_costs" : 113,
     "total_uber_duration" : 115,
     "total_distance" : 27.5
  
}

### Check a trip 
GET        /trips/{trip_id}

Request:  GET             /trips/4083

Response:
{

    "id" : "4083",
    "status" : "planning",
    "starting_from_location_id: "10040",
    "best_route_location_ids" : [ "30000", "10000", "10002", "20005" ],
    "total_uber_costs" : 113,
     "total_uber_duration" : 115,
     "total_distance" : 27.5
     
}

## Start Trip 
//Request destination 1 (after which destination 2 to n)
PUT        /trips/{trip_id}/request

Request:  PUT             /trips/4083/request

Response:
{

     "id" : "4083",
     "status" : "requesting",
     "starting_from_location_id”: "10040",
     "next_destination_location_id”: "30000",
     "best_route_location_ids" : [ "30000", "10000", "10002", "20005" ],
     "total_uber_costs" : 113,
     "total_uber_duration" : 115,
     "total_distance" : 27.5,
     "uber_wait_time_eta" : 7
     
}

## After destination n status updated to complete 

Response:
{

     "id" : "4083",
     "status" : "completed",
     "starting_from_location_id”: "10040",
     "next_destination_location_id”: "",
     "best_route_location_ids" : [ "30000", "10000", "10002", "20005" ],
     "total_uber_costs" : 113,
     "total_uber_duration" : 115,
     "total_distance" : 27.5,
     "uber_wait_time_eta" : 0 
     
}
