package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/afreen23/pokedex/internals/pokecache"
)

var base_url = "https://pokeapi.co/api/v2/location-area/"

type Location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var ID int = 1
var current_id *int = &ID

var cache *pokecache.Cache = pokecache.NewCache(5 * time.Second)

func fetchData(id int) ([]byte, error) {
	fetch_url := fmt.Sprintf("%s/%v", base_url, id)
	// finding in cache first
	res, ok := cache.Get(fetch_url)
	// making GET request if not found in cache
	if !ok {
		res, err := http.Get(fetch_url)
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			return nil, fmt.Errorf("response failed with status code: %d and \nbody: %s\n url: %s", res.StatusCode, body, fetch_url)
		}
		if err != nil {
			return nil, err
		}
		cache.Add(fetch_url, body)
		return body, nil
	}
	return res, nil
}

func printError(err error) {
	fmt.Println(err)
}

func Next20Locations() {
	locations := Location{}
	limit := 20
	if *current_id != 1 {
		limit = limit + *current_id
	}
	for ; *current_id <= limit; *current_id = *current_id + 1 {
		res, err := fetchData(*current_id)
		if err != nil {
			printError(err)
		}
		err = json.Unmarshal(res, &locations)
		if err != nil {
			printError(err)
		}
		fmt.Println(locations.Name)
	}
	*current_id = *current_id - 1
}

func Prev20Locations() {
	if *current_id == 1 {
		fmt.Printf("No previous locations to look back\n")
		return
	}
	locations := Location{}
	limit := *current_id - 20
	for ; *current_id > limit; *current_id = *current_id - 1 {
		res, err := fetchData(*current_id)
		if err != nil {
			printError(err)
		}
		err = json.Unmarshal(res, &locations)
		if err != nil {
			printError(err)
		}
		fmt.Println(locations.Name)
	}
	*current_id = *current_id + 1
}
