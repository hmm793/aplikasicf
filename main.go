package main

import (
	"aplikasicf/auth"
	"aplikasicf/campaign"
	"aplikasicf/handler"
	"aplikasicf/helper"
	"aplikasicf/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/aplikasicf?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connection To Database is success")

	////////////////////////////////////////////////////////////
	
	// User 
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewServive()
	userHandler := handler.NewUserHandler(userService, authService)
	

	// Campaign 
	campaignRepository := campaign.NewRepository(db)

	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewHandler(campaignService)

	// input := campaign.CreateCampaignInput{}
	// input.Name = "Penggalangan Dana Startup"
	// input.ShortDescription = "Short"
	// input.Description = "Long"
	// input.GoalAmount = 1000000
	// input.Perks = "hadiah satu, dua, tiga"
	// inputUser, _ := userService.GetUserByID(1)
	// input.User =  inputUser

	// newCampaign,err := campaignService.CreateCampaign(input)
	// if err!= nil {
	// 	log.Fatal(err.Error())
	// }
	
	// fmt.Println(newCampaign)




	// campaigns,_ := campaignService.GetCampaigns(200)
	// fmt.Println(campaigns)

	// campaigns, err := campaignRepository.FindByUserID(1)

	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// for _, campaign := range campaigns {
	// 	// fmt.Println(campaign.Name)
	// 	fmt.Println("campaign.CampaignImages :: ",campaign.CampaignImages)
	// }


	// token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2fQ.36dr8oCj8qz8Y0Xf_bfRxXjX54zNq0LCk5bYhKZtwiw")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// if token.Valid {
	// 	fmt.Println("Token Valid")
	// } else {
	// 	fmt.Println("Token alid")
			
	// }







	// fmt.Println(authService.GenerateToken(1001))

	
	// userService.SaveAvatar(1, "images/1-profile.png")
	
	
	// var inputData user.LoginInput
	// inputData.Email = "eric@gmail.com"
	// inputData.Password = "1password"
	
	// userValid, err := userService.LoginUser(inputData)

	// if err!=nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(userValid)
	
	
	
	// userByEmail, err := userRepository.FindByEmail("indra@gmail.com")

	// if err!=nil {
	// 	fmt.Println(err.Error())
	// }

	// if userByEmail.ID == 0 {
	// 	fmt.Println("User Tidak Di Temukan")
	// } else {
	// 	fmt.Println(userByEmail.Name)
	// }
	
	router := gin.Default()
	router.Static("/images", "./images")
	api := router.Group("api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/campaigns",  campaignHandler.GetCampaigns)
	api.POST("/campaigns", authMiddleware(authService, userService),campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService),campaignHandler.Update)
	api.GET("/campaigns/:id",  campaignHandler.GetCampaign)
	router.Run(":8080")
}


func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func (c *gin.Context){
		authHeader := c.GetHeader("Authorization")
	
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized,"error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return 
		}
	
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		
		if len(arrayToken) == 2{
			tokenString = arrayToken[1]
		}
	
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized,"error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return 
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized,"error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return 
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized,"error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return 
		}

		c.Set("currentUser", user)

	
	}
}

