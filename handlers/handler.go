package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

// HomeHandler serves the homepage by rendering a list of artists.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if err := validateMethod(r, http.MethodGet); err != nil {
		displayError(w, "Method Not Allowed", err.Error(), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		displayError(w, "Page Not Found", "The page you are looking for is unavailable", http.StatusNotFound)
		return
	}

	artists, err := FetchArtists()
	if err != nil {
		displayError(w, "Internal Server Error", "Error fetching the artists", http.StatusInternalServerError)
		return
	}

	if err := renderTemplate(w, "home.html", artists); err != nil {
		displayError(w, "Internal Server Error", "Error rendering the homepage", http.StatusInternalServerError)
	}
}

// ArtistHandler serves individual artist pages based on the artist ID.
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	if err := validateMethod(r, http.MethodGet); err != nil {
		displayError(w, "Method Not Allowed", err.Error(), http.StatusMethodNotAllowed)
		return
	}

	artistIDStr, err := extractArtistID(r.URL.Path)
	if err != nil {
		displayError(w, "Invalid Artist ID", err.Error(), http.StatusBadRequest)
		return
	}

	artistID, err := strconv.Atoi(artistIDStr)
	if err != nil || artistID < 1 || artistID > 52 {
		displayError(w, "Artist Not Found", "Artist with the given ID does not exist", http.StatusNotFound)
		return
	}

	artist, err := findArtistByID(artistID)
	if err != nil {
		displayError(w, "Internal Server Error", err.Error(), http.StatusInternalServerError)
		return
	}

	datesLocations, err := FetchRelations(artistID)
	if err != nil {
		displayError(w, "Internal Server Error", "Error fetching artist relations.", http.StatusInternalServerError)
		return
	}

	data := struct {
		Artist         Artist
		DatesLocations map[string][]string
	}{
		Artist:         artist,
		DatesLocations: datesLocations,
	}

	if err := renderTemplate(w, "artist.html", data); err != nil {
		displayError(w, "Internal Server Error", "Error rendering the artist page.", http.StatusInternalServerError)
	}
}

// Utility function to validate the HTTP method of a request.
func validateMethod(r *http.Request, allowedMethod string) error {
	if r.Method != allowedMethod {
		return fmt.Errorf("The method %s is not allowed for the requested URL", r.Method)
	}
	return nil
}

// Extracts and validates the artist ID from the URL path.
func extractArtistID(path string) (string, error) {
	parts := strings.Split(path, "/")
	if parts[2] == "" {
		return "", fmt.Errorf("ID cannot be empty")
	}

	artistIDStr := parts[2]

	// Validate artist ID characters.
	for _, char := range artistIDStr {
		if !unicode.IsDigit(char) {
			return "", fmt.Errorf("ID contains invalid characters, only numerical values are allowed")
		}
	}

	// Additional formatting checks.
	if strings.HasPrefix(artistIDStr, "0") && artistIDStr != "0" {
		return "", fmt.Errorf("ID has an invalid format, leading zeros are not allowed")
	}

	if strings.Contains(artistIDStr, " ") || strings.Contains(artistIDStr, ".") || len(artistIDStr) > 10 {
		return "", fmt.Errorf("ID must be a positive integer and cannot contain spaces or decimal points")
	}

	return artistIDStr, nil
}

// Finds an artist by their ID.
func findArtistByID(artistID int) (Artist, error) {
	artists, err := FetchArtists()
	if err != nil {
		return Artist{}, fmt.Errorf("Error fetching artists: %v", err)
	}

	for _, artist := range artists {
		if artist.ID == artistID {
			return artist, nil
		}
	}

	return Artist{}, fmt.Errorf("Artist with ID %d not found", artistID)
}



// Helper function to render templates.
func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) error {
	err := tpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		return fmt.Errorf("Error rendering template %s: %v", templateName, err)
	}
	return nil
}

type ErrorData struct {
	ErrorTitle string
	ErrorMessage string
	ErrorCode int
}

// displayError renders an error page with a custom message.
func displayError(w http.ResponseWriter, errTitle string, errMsg string, errCode int) {
	data := ErrorData{
		ErrorTitle:   errTitle,
		ErrorMessage: errMsg,
		ErrorCode:    errCode,
	}
	w.WriteHeader(errCode)
	err := tpl.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		log.Printf("template execution error: %v", err)
	}
}