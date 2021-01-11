package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//date for covid19api
type CovidData struct {
	Confirmed int64
}
type CountryPopulation struct {
	Body CountryData
}

//data for country population
type CountryData struct {
	Population int64
}

func main() {
	http.HandleFunc("/covid19chances/", calculatechances)
	http.ListenAndServe(":8000", nil)

}

//gets the number of confirmed covid19cases
//gets the country name and returns the data with type CovidData
func getcovidcases(cname string) []CovidData {

	url := "https://covid1910.p.rapidapi.com/data/confirmed/country/" + strings.ToLower(cname)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-host", "covid1910.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "1a76ca5476msh3925d740f2910e5p1a25f7jsn0013a8934bc6")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	var data []CovidData
	json.Unmarshal([]byte(body), &data)
	//fmt.Println(data)
	return data
	/*
		fmt.Fprintf(w, "Number of confirmed cases: ")
		fmt.Fprintln(w, strconv.Itoa(int(data[0].Provinces[0].Confirmed)))
		fmt.Fprintf(w, "Number of recovered: ")
		fmt.Fprintln(w, strconv.Itoa(int(data[0].Provinces[0].Recovered)))
		fmt.Fprintf(w, "Number of Deaths: ")
		fmt.Fprintln(w, strconv.Itoa(int(data[0].Provinces[0].Deaths)))*/
}

//gets the country population bases on the cname and returns the population with CountryPopulation type
func getpopulation(cname string) CountryPopulation {
	url := "https://world-population.p.rapidapi.com/population?country_name=" + cname

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-host", "world-population.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "1a76ca5476msh3925d740f2910e5p1a25f7jsn0013a8934bc6")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var countrydata CountryPopulation
	json.Unmarshal([]byte(body), &countrydata)
	//fmt.Println(res)
	return countrydata
}

type Response struct {
	Country             string  `json:country`
	ConfirmedCases      int64   `json:confirmed`
	Yourchanceofgetting float64 `json:yourchanceofgetting`
}

//calculates the chance of getting the virus and sends a json response
func calculatechances(w http.ResponseWriter, req *http.Request) {
	country, ok := req.URL.Query()["country"]
	//country := coun
	if !ok || len(country[0]) < 1 {
		fmt.Fprintf(w, "404 country not provided")
		return
	}
	var population CountryPopulation
	population = getpopulation(string(country[0]))
	//fmt.Fprintln(w, population.Body.Population)

	var covidcases []CovidData
	covidcases = getcovidcases(string(country[0]))

	var chances float64 = (float64(covidcases[0].Confirmed) / float64(population.Body.Population)) * 100
	//fmt.Fprintln(w, covidcases[0].Confirmed)
	//fmt.Fprintln(w, population.Body.Population)
	//fmt.Fprintln(w, chances)
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Country:             string(country[0]),
		ConfirmedCases:      covidcases[0].Confirmed,
		Yourchanceofgetting: chances}

	json.NewEncoder(w).Encode(response)
}
