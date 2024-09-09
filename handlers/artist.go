package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Artist struct {
	ID               int      `json:"id"`
	Image            string   `json:"image"`
	Name             string   `json:"name"`
	Members          []string `json:"members"`
	CreationDate     int      `json:"creationDate"`
	FirstAlbum       string   `json:"firstAlbum"`
	LocationsJson    string   `json:"locations"`
	ConcertDatesJson string   `json:"concertDates"`
	RelationsJson    string   `json:"relations"`
}

type DatesLocations map[string][]string

type Relations struct {
	ID             int            `json:"id"`
	DatesLocations DatesLocations `json:"datesLocations"`
}

// FetchArtists fetches and processes the list of artists from the API
func FetchArtists() ([]Artist, error) {
	const url = "https://groupietrackers.herokuapp.com/api/artists"

	// Perform the HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artists: %w", err)
	}
	defer res.Body.Close()

	// Check for a successful response code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response into the artist slice
	var artists []Artist
	if err := json.Unmarshal(body, &artists); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artists: %w", err)
	}

	return artists, nil
}

// FetchRelations fetches the concert dates and locations for a specific artist
func FetchRelations(artistID int) (DatesLocations, error) {
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", artistID)

	// Perform the HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch relations for artist %d: %w", artistID, err)
	}
	defer res.Body.Close()

	// Check for a successful response code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response into the relations struct
	var relations Relations
	if err := json.Unmarshal(body, &relations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal relations for artist %d: %w", artistID, err)
	}

	return relations.DatesLocations, nil
}

