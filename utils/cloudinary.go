package utils

import (
	"context"
	"fmt"
	"os"
	"io"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// ConfigureCloudinary configures the Cloudinary instance using environment variables
func ConfigureCloudinary() (*cloudinary.Cloudinary, error) {
	// Load Cloudinary credentials from environment variables
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Check if any required environment variable is missing
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("missing required Cloudinary environment variables")
	}

	// Initialize Cloudinary client with credentials
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to configure Cloudinary: %v", err)
	}

	return cld, nil
}

// UploadImageToCloudinary uploads the image to Cloudinary and returns the secure URL
func UploadImageToCloudinary(file io.Reader) (string, error) {
	// Configure Cloudinary with environment variables
	cld, err := ConfigureCloudinary()
	if err != nil {
		return "", err
	}

	// Upload image to Cloudinary
	uploadResult, err := cld.Upload.Upload(
		context.Background(), // Background context for the request
		file,                 // File input (must implement io.Reader)
		uploader.UploadParams{
			Folder: "categories", // You can set a folder for the image
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %v", err)
	}

	// Return the secure URL of the uploaded image
	return uploadResult.SecureURL, nil
}
