package auth

import (
	"net/http"
	"server/db"
	"server/models"
	"path/filepath"
	"os"
	"io"
	"github.com/gin-gonic/gin"
	"server/utils"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := HashPassword(input.Password)
	user := models.User{Username: input.Username, Email: input.Email, PasswordHash: hashedPassword}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
		return

	}

	LogAudit(user.ID, "signup", "user", user.ID, "New user registered")

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !CheckPasswordHash(input.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, _ := GenerateToken(user.ID)

	LogAudit(user.ID, "login", "user", user.ID, "User logged in")


	c.JSON(http.StatusOK, gin.H{"token": token})
}


func UploadYAML(c *gin.Context) {
	// Parse the file from the form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid file upload"})
		return
	}
	defer file.Close()

	// Ensure it's a YAML file
	if filepath.Ext(header.Filename) != ".yaml" {
		c.JSON(400, gin.H{"error": "Only .yaml files are allowed"})
		return
	}

	// Define the save path
	savePath := "./config/jobs/" + header.Filename

	// Create or replace the file
	out, err := os.Create(savePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save the file"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to write file content"})
		return
	}

	userID := c.MustGet("userID").(uint)

	fileContent, _ := os.ReadFile(savePath)

	checksum := utils.ComputeChecksum(fileContent)

	yamlUpload := models.YAMLUpload{
		UserID:     userID,
		FileName:   header.Filename,
		Checksum:   checksum,
		Content:    string(fileContent),
	}

	db.DB.Create(&yamlUpload)

	LogAudit(userID, "upload", "yaml_file", yamlUpload.ID, "Uploaded YAML file "+header.Filename)	

	c.JSON(200, gin.H{"message": "File uploaded successfully", "file": header.Filename})
}

func Profile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	LogAudit(userID, "profile_access", "user", userID, "User accessed profile")
	c.JSON(http.StatusOK, gin.H{"message": "Access granted", "user_id": userID})
}

func LogAudit(userID uint, actionType, targetType string, targetID uint, details string) {
	audit := models.AuditLog{
		UserID:         userID,
		ActionType:     actionType,
		TargetType:     targetType,
		TargetID:       targetID,
		Details:        details,
	}
	db.DB.Create(&audit)
}