package uber

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	// "net/url"
    // "strconv"
    "io/ioutil"
	// "gopkg.in/mgo.v2"
	 // "gopkg.in/mgo.v2/bson"
)

// List of price estimates
type PriceEstimates struct {
	Prices         []PriceEstimate `json:"prices"`
}

// Uber price estimate
type PriceEstimate struct {
	ProductId       string  `json:"product_id"`
	CurrencyCode    string  `json:"currency_code"`
	DisplayName     string  `json:"display_name"`
	Estimate        string  `json:"estimate"`
	LowEstimate     int     `json:"low_estimate"`
	HighEstimate    int     `json:"high_estimate"`
	SurgeMultiplier float64 `json:"surge_multiplier"`
	Duration        int     `json:"duration"`
	Distance        float64 `json:"distance"`
}

type UberOutput struct{
	Cost int
	Duration int
	Distance float64
}

type UberETA struct{
  Request_id string `json:"request_id"`
  Status string	    `json:"status"`
  Vehicle string	`json:"vehicle"`
  Driver string		`json:"driver"`
  Location string	`json:"location"`
  ETA int			`json:"eta"`
  SurgeMultiplier float64 `json:"surge_multiplier"`
}

func Get_uber_price(startLat, startLon, endLat, endLon string) UberOutput{
	client := &http.Client{}
	reqURL := fmt.Sprintf("https://sandbox-api.uber.com/v1/estimates/price?start_latitude=%s&start_longitude=%s&end_latitude=%s&end_longitude=%s&server_token=sh4TxBDenZAnaa0XMFQXgyGEI4hOt60ZRflX-3rb", startLat, startLon, endLat, endLon)
	fmt.Println("URL formed: "+ reqURL)
	// res, err := http.GET(reqURL,)
	req, err := http.NewRequest("GET", reqURL , nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in sending req to Uber: ", err);	
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error in reading response: ", err);	
	}

	var res PriceEstimates
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("error in unmashalling response: ", err);	
	}

	var uberOutput UberOutput
	uberOutput.Cost = res.Prices[0].LowEstimate
	uberOutput.Duration = res.Prices[0].Duration
	uberOutput.Distance = res.Prices[0].Distance

	return uberOutput

}


func Get_uber_eta(startLat, startLon, endLat, endLon string) int{

	var jsonStr = []byte(`{"start_latitude":"` + startLat + `","start_longitude":"` + startLon + `","end_latitude":"` + endLat + `","end_longitude":"` + endLon + `","product_id":"04a497f5-380d-47f2-bf1b-ad4cfdcb51f2"}`)
	reqURL := "https://sandbox-api.uber.com/v1/requests"
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicHJvZmlsZSIsInJlcXVlc3QiXSwic3ViIjoiNzcyMjZiYWMtMzJiMC00YzMzLWEwNWYtYWI0ODBjNTUyOTg3IiwiaXNzIjoidWJlci11czEiLCJqdGkiOiJkZDEzMTQ1My1kNzFlLTQ5NWUtOThjNi03ZTE4NDllMDkwYjciLCJleHAiOjE0NTA1OTkzOTAsImlhdCI6MTQ0ODAwNzM5MCwidWFjdCI6IjF3dk0zOHdvdUR2Z1g4V29jM25SdDZadmNjcjhyTyIsIm5iZiI6MTQ0ODAwNzMwMCwiYXVkIjoiT0Rkb0VkeUxqLVA5VzdkMXkzSk5JSldsUDA5anZTbXQifQ.JmZ2EN8bcx0eo3K_B51Wk29amm7_u8FV0l5mdQm1s5i7RNuI9t9pVh7t0To28OnvOs_RCj5nIIz5t7ggZgPoed9mmVU7rndXDec-O6cbjUdY45zAJ3ZsW_joNiCircyOsyDpLXtDvCi5J93TTrLJ6vJ9JGx1Slg3mxLoI4-VrVZIDoAOnTg5AmQxop0-yrgOZpombkf4jK4APBZArygEpHr7QDBXgxMkjnIlu_Nxftb8BDtOS9XX-8fRcC4A5WZQT-f8mMpH4kRr80_YEgMi4CgXRmsgdhTDYr9-6G1VOdQCNJODbFJno4_jCDDCl6OX-UTyfvcabjzdXZwCIviF-Q")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in sending req to Uber: ", err);	
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error in reading response: ", err);	
	}

	var res UberETA
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("error in unmashalling response: ", err);	
	}
	eta:= res.ETA
	return eta
	
}
	
