package handler

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

var PATH = []string{}

type Artist struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
	Description  string
}

type Dates struct {
	Index []struct {
		Id    int
		Dates []string
	}
}

type Locations struct {
	Index []struct {
		ID        int
		Locations []string
		Dates     string
	}
}

type Loc struct {
	Artists         []string
	Location        string
	ListOfLocations []string
	//Dates    []string
}

type Relations struct {
	Index []struct {
		ID             int
		DatesLocations map[string][]string //interface{}
	}
}

type Data struct {
	ListOfArtists []Artist
	PageNumber    []int
}

var Maintemp = OpenTemplate("index")
var ArtistTemp = OpenTemplate("artist")
var FormRoute = []string{"pages"}
var ListOfArtist []Artist

func RoutingHandler(rw http.ResponseWriter, r *http.Request) {
	PATH = GetUrl(r)
	if PATH[0] == "locations" {
		LocationsHandler(rw, r)
	}
	if PATH[0] == "artists" {
		if len(PATH) > 1 {
			ArtistHandler(rw, r)
		} else {
			AllArtistsHandler(rw, r)
		}
	} else {
		Error404Handler(rw, r)
	}

}

func AllArtistsHandler(w http.ResponseWriter, r *http.Request) {
	isPaginated := false
	page := []int{}
	listmp := []Artist{}

	if len(r.URL.Query()["page"]) == 0 {
		SearchInApi("artists", &ListOfArtist)
	}

	if r.Method == "POST" {
		for i := 0; i < len(ListOfArtist); i++ {
			if strings.Contains(strings.ToUpper(ListOfArtist[i].Name), strings.ToUpper(r.FormValue("artists"))) {
				listmp = append(listmp, ListOfArtist[i])
			}
		}
		ListOfArtist = listmp
	}

	for i := 1; i <= len(ListOfArtist)/12+1; i++ {
		if len(r.URL.Query()["page"]) > 0 && fmt.Sprintf("%d", i) == r.URL.Query()["page"][0] {
			isPaginated = true
			if i*12 > len(ListOfArtist) {
				listmp = ListOfArtist[12*(i-1):]
			} else {
				listmp = ListOfArtist[12*(i-1) : 12*i]
			}
		}
		page = append(page, i)
	}

	if !isPaginated && len(ListOfArtist) > 12 {
		ListOfArtist = ListOfArtist[:12]
	} else {
		ListOfArtist = listmp
	}

	Maintemp.Execute(w, Data{ListOfArtist, page})
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	var artist Artist
	SearchInApi(fmt.Sprintf("artists/%s", PATH[1]), &artist)
	GetWiki(&artist)
	ArtistTemp.Execute(w, artist)
}

func LocationsHandler(rw http.ResponseWriter, r *http.Request) {
	Maintemp = OpenTemplate("locations")
	var locations Locations
	var location string
	var indexes []int
	var ArtistsinArea []string
	var listOfListsOfLocations [][]string
	var listOfLocations []string

	SearchInApi("locations", &locations)
	listOfArtist := &[]Artist{}
	SearchInApi("artists", listOfArtist)
	var listOfDates Dates
	SearchInApi("dates", &listOfDates)
	var listOfRelations Relations
	SearchInApi("relation", &listOfRelations)
	fmt.Println(listOfRelations.Index[0].DatesLocations["los_angeles-usa"])

	for i := 0; i < len(locations.Index); i++ {
		listOfListsOfLocations = append(listOfListsOfLocations, locations.Index[i].Locations)
	}
	for j := 0; j < 51; j++ {
		for s := 0; s < len(listOfListsOfLocations[j]); s++ {
			listOfLocations = append(listOfLocations, listOfListsOfLocations[j][s])
		}
	}
	sort.Sort(sort.StringSlice(listOfLocations)) //Sort List of Locations alphabetically
	fmt.Println("1:", len(listOfLocations))
	fmt.Println("2:", len(RemoveDuplicateStr(listOfLocations)))
	if r.Method == "POST" {
		fmt.Println(r.FormValue("locations"))
		for i := 0; i <= 51; i++ {
			for j := 0; j < len(locations.Index[i].Locations); j++ {
				if strings.Contains(strings.ToUpper(locations.Index[i].Locations[j]), strings.ToUpper(strings.ReplaceAll(r.FormValue("locations"), " ", "_"))) {
					location = "https:www.google.com/maps/embed/v1/place?key=AIzaSyAXXPpGp3CYZDcUSiE2YRlNID4ybzoZa7o&q=" + locations.Index[i].Locations[j]
					//date = append(date, listOfDates.Index[i].Dates[j])
					indexes = append(indexes, i)
					for s := 0; s < len(listOfRelations.Index[i].DatesLocations[locations.Index[i].Locations[j]]); s++ {
						ArtistsinArea = append(ArtistsinArea, (*listOfArtist)[i].Name+" in "+listOfRelations.Index[i].DatesLocations[locations.Index[i].Locations[j]][s]) //listOfDates.Index[i].Dates[j])
					}
				}
			}
		}
	}
	start := Loc{ArtistsinArea, location, RemoveDuplicateStr(listOfLocations)}
	Maintemp.Execute(rw, start)
}

func RemoveDuplicateStr(strSlice []string) []string { //We use this function to remove duplicate strings in an array of strings
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func Error404Handler(rw http.ResponseWriter, r *http.Request) {
	OpenTemplate("err404").Execute(rw, nil)
}
