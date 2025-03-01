package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// GetBusinessPhotos fetches photos from the Google Maps Data API
func GetBusinessPhotos(c *gin.Context) {
	url := "https://maps-data.p.rapidapi.com/photos.php?business_id=0x47e66e2964e34e2d%3A0x8ddca9ee380ef7e0&lang=en&country=IN"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	rapidAPIKey := os.Getenv("RAPIDAPI_KEY")
	if rapidAPIKey == "" {
		fmt.Println("RAPIDAPI_KEY is not set in .env")
		return
	}

	// Add required headers
	req.Header.Add("x-rapidapi-key", rapidAPIKey)
	req.Header.Add("x-rapidapi-host", "maps-data.p.rapidapi.com")

	// Send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Send the response as JSON
	c.Data(http.StatusOK, "application/json", body)
}

// InitRoutes sets up the routes for the API
