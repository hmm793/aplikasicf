package handler

import (
	"aplikasicf/auth"
	"aplikasicf/helper"
	"aplikasicf/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

// Handler 
func (h *userHandler) RegisterUser(c *gin.Context){
	// tangkap input dari user 
	// map input dari user ke struct RegsiterUserInput
	// struct di atas kita passing sebagai parameter service

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessages := gin.H{
			"errors" : errors,
		}
		response := helper.APIResponse("Registered Account Failed", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	NewUser, err := h.userService.RegisterUser(input)
	
	if err != nil {
		response := helper.APIResponse("Registered Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	
	}


	token, err := h.authService.GenerateToken(NewUser.ID)
	if err != nil {
		response := helper.APIResponse("Registered Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// token, err := h.userService.RegisterUser(input)
	formatter := user.FormatUser(NewUser, token)
	response := helper.APIResponse("Account Has Been Registered", http.StatusOK, "succes", formatter)
	
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context){
	// User memasukkan input (email dan password)
	// input di tangkap handler
	// mapping dari input user ke input struct
	// input struct passing service
	// di service mencari dengan bantuan repository user dengan email x
	// mencocokkan password

	var input user.LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessages := gin.H{
			"errors" : errors,
		}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loggedinUser, err := h.userService.LoginUser(input)

	if err != nil {
		errorMessages := gin.H{
			"errors" : err.Error(),
		}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusBadRequest, response)
		return
	}


	token, err := auth.NewServive().GenerateToken(loggedinUser.ID)
	
	if err != nil {
		errorMessages := gin.H{
			"errors" : err.Error(),
		}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	formatter := user.FormatUser(loggedinUser, token)
	response := helper.APIResponse("Successfully Loggedin", http.StatusOK, "succes", formatter)
	
	c.JSON(http.StatusOK, response)
	
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context){
	// ada input email dari user 
	// input email di mapping ke struct input
	// struct input di passing ke service
	// service akan memanggil repository untuk menentukan apakah email sudah ada
	// repository akan mengakukan query ke database

	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessages := gin.H{
			"errors" : errors,
		}
		response := helper.APIResponse("Email Checking Failed", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)

	if err != nil {
		errorMessages := gin.H{
			"errors" : "Server Error",
		}
		response := helper.APIResponse("Email Checking Failed", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_available" : isEmailAvailable,
	}

	metaMessage := "Email Address has been registered"

	if isEmailAvailable {
		metaMessage = "Email Is Available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "succes", data)
	
	c.JSON(http.StatusOK, response)
}


func (h *userHandler) UploadAvatar(c *gin.Context) {
	// tangkap input dari user
	// simpan gambar di folder "images/"
	// di service kita panggil repo
	// JWT ( sementara hardcore, seakan2 user yang login ID = 1)
	// repo ambil data user dengan id = 1
	// repo update data user simpan lokasi file
	file, err := c.FormFile("avatar")

	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}


	// // sementara hardcore
	// userID := 2

	// Sudah dinamik 
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	// path := "images/" + string(userID) + file.Filename
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)

	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_uploaded" : true,
	}
	response := helper.APIResponse("Avatar successfully uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
	return
	
}

