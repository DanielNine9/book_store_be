package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Cấu hình Cloudinary với thông tin từ biến môi trường
func ConfigureCloudinary() (*cloudinary.Cloudinary, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to configure Cloudinary: %v", err)
	}

	return cld, nil
}

// Hàm upload hình ảnh lên Cloudinary và trả về URL của hình ảnh
func UploadImageToCloudinary(filePath string) (string, error) {
	cld, err := ConfigureCloudinary()
	if err != nil {
		return "", err
	}

	// Sử dụng context.Background() và truyền vào hàm Upload
	uploadResult, err := cld.Upload.Upload(
		context.Background(), // Thêm context vào đây
		filePath,
		uploader.UploadParams{
			Folder: "categories", // Chọn folder nếu cần
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %v", err)
	}

	return uploadResult.SecureURL, nil
}
