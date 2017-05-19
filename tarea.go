package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type Dirreq struct {
	Origin     string `json:"origin"`
	Desination string `json:"destination"`
}

type Restaurantreq struct {
	Location string `json:"location"`
}

func getDirections(w http.ResponseWriter, p *http.Request) {
	var orig Dirreq
	//orig.Origen = "Hannover"
	//orig.Destino = "Luxembourg"
	_ = json.NewDecoder(p.Body).Decode(&orig)

	c, erro := maps.NewClient(maps.WithAPIKey("AIzaSyAkMHjwcRtPz9trV2DUb6N6R3E_m7aH91s"))
	if erro != nil {
		log.Fatalf("fatal error: %s", erro)
	}

	r := &maps.DirectionsRequest{
		Origin:      orig.Origin,
		Destination: orig.Desination,
	}

	listarutas, _, errr := c.Directions(context.Background(), r)
	if errr != nil {
		log.Fatalf("fatal error: %s", errr)
	}

	rutaStr := new(bytes.Buffer)
	rutaStr.WriteString("{\n\"ruta\":[\n")
	json.NewDecoder(p.Body).Decode(&listarutas)

	for ind := 0; ind < len(listarutas[0].Legs[0].Steps); ind++ {
		rutaStr.WriteString("{\n   \"lat\": ")
		rutaStr.WriteString(strconv.FormatFloat(listarutas[0].Legs[0].Steps[ind].StartLocation.Lat, 'f', 4, 64))
		rutaStr.WriteString(",")
		rutaStr.WriteString(" \n   \"long\": ")
		rutaStr.WriteString(strconv.FormatFloat(listarutas[0].Legs[0].Steps[ind].StartLocation.Lng, 'f', 4, 64))
		rutaStr.WriteString("\n},\n")
		if ind == (len(listarutas[0].Legs[0].Steps) - 1) {
			rutaStr.WriteString("{\n   \"lat\": ")
			rutaStr.WriteString(strconv.FormatFloat(listarutas[0].Legs[0].Steps[ind].StartLocation.Lat, 'f', 4, 64))
			rutaStr.WriteString(",")
			rutaStr.WriteString(" \n   \"long\": ")
			rutaStr.WriteString(strconv.FormatFloat(listarutas[0].Legs[0].Steps[ind].StartLocation.Lng, 'f', 4, 64))
			rutaStr.WriteString("\n}")
		}
	}
	rutaStr.WriteString("\n ]\n}")
	fmt.Fprintf(w, rutaStr.String())
	//pretty.Println(buff.String())
}

func getRestaurants(w http.ResponseWriter, p *http.Request) {
	var locat Restaurantreq
	_ = json.NewDecoder(p.Body).Decode(&locat)

	dClient, err := maps.NewClient(maps.WithAPIKey("AIzaSyAkMHjwcRtPz9trV2DUb6N6R3E_m7aH91s"))

	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := &maps.DirectionsRequest{
		Origin:      locat.Location,
		Destination: locat.Location,
	}

	route, _, errr := dClient.Directions(context.Background(), r)

	if errr != nil {
		log.Fatalf("fatal error: %s", errr)
	}

	json.NewDecoder(p.Body).Decode(&route)

	cl, _ := maps.NewClient(maps.WithAPIKey("AIzaSyAkMHjwcRtPz9trV2DUb6N6R3E_m7aH91s"))

	var latitud float64 = route[0].Legs[0].Steps[0].StartLocation.Lat
	var longitud float64 = route[0].Legs[0].Steps[0].StartLocation.Lng

	req := &maps.NearbySearchRequest{
		Location: &maps.LatLng{latitud, longitud},
		Radius:   10000,
		Type:     maps.PlaceTypeRestaurant,
		Keyword:  "restaurant",
		//poner maps.placetyperestaurant o "restaurant" de alguna manera retorna hoteles, que asumo que poseen restaurantes tambien
		//poner keyword restaurants lo mejora
	}

	listarestau, _ := cl.NearbySearch(context.Background(), req)
	json.NewDecoder(p.Body).Decode(&listarestau)

	rstrnts := new(bytes.Buffer)
	rstrnts.WriteString("{\n\"Restaurantes\":[\n")

	for ind := 0; ind < len(listarestau.Results); ind++ {
		rstrnts.WriteString("{\n\"nombre\": ")
		rstrnts.WriteString("\"" + listarestau.Results[ind].Name + "\",\n")
		rstrnts.WriteString("\"lat\": ")
		rstrnts.WriteString(strconv.FormatFloat(listarestau.Results[ind].Geometry.Location.Lat, 'f', 4, 64))
		rstrnts.WriteString(",\n")
		rstrnts.WriteString("\"long\": ")
		rstrnts.WriteString(strconv.FormatFloat(listarestau.Results[ind].Geometry.Location.Lng, 'f', 4, 64))
		rstrnts.WriteString("\n  },\n")
	}

	rstrnts.WriteString(" ]\n}")
	fmt.Fprintf(w, rstrnts.String())
	//pretty.Println(rstrnts.String())

}

func main() {
	srvs := mux.NewRouter()
	srvs.HandleFunc("/ej1", getDirections).Methods("POST")
	srvs.HandleFunc("/ej2", getRestaurants).Methods("POST")

	log.Println("Listening to http://localhost:8686...")
	log.Fatal(http.ListenAndServe(":8686", srvs))
}
