package regions

import (
	"api/controllers/handlers"
	"api/logic"
	"api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegionResponse struct {
	RegionID uint   `json:"region_id"`
	Name     string `json:"name"`
}

func GetRegions(c *gin.Context) {
	regions, err := logic.GetRegions()

	if err != nil {
		handlers.HandleUnknownError(c, err)
		return
	}

	response := CreateRegionsResponse(regions)
	c.JSON(http.StatusOK, response)
}

func CreateRegionsResponse(regions []models.Region) []RegionResponse {
	var response []RegionResponse

	for _, region := range regions {
		response = append(response, RegionResponse{
			RegionID: region.RegionID,
			Name:     region.Name,
		})
	}

	return response
}
