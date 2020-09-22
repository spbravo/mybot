package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	//"math/big"
)

type Data struct {
	Descripcion string 
	Estado int 
	Datos string
	Metadatos string 
}

type Result struct{
	Idema string `json:"idema"` 
	Longitude float64 `json:"lon"`
	Time string `json:"fint"`
	Precision float64 `json:"prec"`
	Altitude float64 `json:"alt"`
	Vmax float64 `json:"vmax"`
	VV float64 `json:"vv"`
	Dv float64 `json:"dv"`
	Latitude float64 `json:"lat"`
	Dmax float64 `json:"dmax"`
	Ubicacion string `json:"ubi"`
	Presion float64 `json:"pres"`
	Hr float64 `json:"hr"`
	PresionMar float64 `json:"pres_nmar"`
	TempMin float64 `json:"tamin"`
	TempActual float64 `json:"ta"`
	TempMax float64 `json:"tamax"`
	Tpr float64 `json:"tpr"`
	Rviento float64 `json:"rviento"`
}

func recuperaDatos(url string) string {
	
		//url := "https://opendata.aemet.es/opendata/api/observacion/convencional/datos/estacion/B228/?api_key=eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJzcGJyYXZvQGdtYWlsLmNvbSIsImp0aSI6IjgyNGM1Y2UzLWFmNzYtNDk0NS1hZDBmLTdhMDk1ZTIyMzJkZCIsImlzcyI6IkFFTUVUIiwiaWF0IjoxNTA5NTU2OTY5LCJ1c2VySWQiOiI4MjRjNWNlMy1hZjc2LTQ5NDUtYWQwZi03YTA5NWUyMjMyZGQiLCJyb2xlIjoiIn0.XAEUT7p_9sXrkavMunL9CtQySwUWOicCIbGfsYxdVZk"
	
		//Aqui se recupera la url para recuperar luego los datos
		req, _ := http.NewRequest("GET", url, nil)
	
		req.Header.Add("cache-control", "no-cache")
		res, _ := http.DefaultClient.Do(req)
	
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		
		var urlDato Data
		err := json.Unmarshal(body, &urlDato)
		if err != nil{ 
			fmt.Println("Estoy en el error %v \n", err)
		}
		
					 
		//Aqui recupero los datos de la estacion	
		resultados := resultsFromJson(contentFromServer(urlDato.Datos))
		
		var tempMax = 0.00
		var hora string
		var ubicacion string
		var lastResult Result
		for _, result := range resultados {
			temper := result.TempActual
			if temper> tempMax {
				tempMax=temper
				hora=result.Time
				ubicacion=result.Ubicacion
			}
			fmt.Printf("%v ($%.2f)\n", result.Ubicacion, temper)
			lastResult = result
		}

	s := fmt.Sprintf("la temperatura máxima en %v ha sido %.2f grados a las %v y la temperatura actual es %.2f", ubicacion, tempMax, hora, lastResult.TempActual) 
	return s
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func contentFromServer(url string) string {
	
	resp, err := http.Get(url)
	checkError(err)
	
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	return string(bytes)
}

func resultsFromJson(content string) []Result {
	results := make([]Result, 0, 24)
	
	decoder := json.NewDecoder(strings.NewReader(content))
	_, err := decoder.Token()
	checkError(err)
	
	var result Result
	for decoder.More() {
		err := decoder.Decode(&result)
		checkError(err)
		results = append(results, result)
	}
	
	return results
}