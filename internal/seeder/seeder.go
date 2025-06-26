package seeder

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"hospital/internal/models"
)

func SeedUsers(db *gorm.DB) {

	var count int64
	db.Model(&models.User{}).Where("email IN ?", []string{"doc@example.com", "reception@example.com"}).Count(&count)

	if count >= 2 {
		log.Println(" Seed users already exist,,,Skipping")
		return
	}

	log.Println(" Seeding initial users")

	users := []models.User{
		{
			Name:         "Dr. Strange",
			Email:        "doc@example.com",
			Role:         "doctor",
			PasswordHash: hashPassword("doc123"),
		},
		{
			Name:         "Receptionist Amy",
			Email:        "reception@example.com",
			Role:         "receptionist",
			PasswordHash: hashPassword("recep123"),
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v\n", user.Email, err)
		}
	}

	log.Println("User seeding complete.")
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("bcrypt error:", err)
	}
	return string(bytes)
}
