package handler

import (
	"aplikasicf/campaign"
	"aplikasicf/helper"
	"aplikasicf/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// tangkap parameter di handler
// handler ke service
// service yang menentukan repository mana yang di call (get all & get by user id)
// repository akses database


type campaignHandler struct {
	service campaign.Service
}

func NewHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}


func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)

	if err != nil {
		response := helper.APIResponse("Error To Get Campaigns", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of Campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
	return
} 


func (h *campaignHandler) GetCampaign(c *gin.Context){
	// api/v1/campaigns/1
	// handler : mapping id yang di url ke struct input => service, call formatter
	// service : inputnya struct input => mengangkap id url 
	// repository : get campaign by id

	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed To Get Detail of Campaign", http.StatusBadRequest, "error",nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed To Get Detail of Campaign", http.StatusBadRequest, "error",nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign detail", http.StatusOK, "success",campaign.FormatDetailCampaign(campaignDetail))

	c.JSON(http.StatusOK, response)
}


// tangkap parameter dari user ke input struct
// ambil current user dari jwt/handler
// panggil service, parameter nya input dstruct dan juga buat slug
// panggil repository untuk suimpan data campaign baru


func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors":errors}
		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)

	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))

	c.JSON(http.StatusOK, response)
}

// Update Campaign
// Handler
// Mapping dari input ke input struct 
// input dari user, dan juga input yang ada di uri (passing ke service)
// service (find campaign by id, tangkap parameter)
// repository update data campaign


func (h *campaignHandler) Update(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to Update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors":errors}
		response := helper.APIResponse("Failed to Update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser


	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	
	if err != nil {
		response := helper.APIResponse("Failed to Update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))

	c.JSON(http.StatusOK, response)
}


// Upload Campaign Images
// handler
// tangkap input dan ubah ke struct input
// save image campaign ke suatu folder
// service (kondisi manggil point 2 di repo, panggil point 1)
// repository 
// 1. create image/save data image ke dalam tabel campaign_images
// 2. ubah is_primary true ke false
