package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"shop-account/models"
	"shop-account/dtos"
	"shop-account/utils"
	"github.com/jinzhu/gorm"
	"fmt"
	
)

type CategoryHandler struct {
	DB *gorm.DB
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
    // var categoryRequest struct {
    //     Name        string `form:"name"`
    //     Description string `form:"description"`
    // }
    var categoryRequest dtos.CategoryRequest

    if err := c.ShouldBind(&categoryRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
        return
    }

    var categoryData models.Category
    categoryData.Name = categoryRequest.Name
    categoryData.Description = categoryRequest.Description

    file, err := c.FormFile("image")

    if err == nil {
        fileContent, err := file.Open()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file", "details": err.Error()})
            return
        }
        defer fileContent.Close() 

        imageURL, err := utils.UploadImageToCloudinary(fileContent)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image", "details": err.Error()})
            return
        }
		fmt.Printf("ImageUrl %s \n", imageURL)
        categoryData.ImageURL = imageURL
    } else if err != nil && err.Error() != "http: no such file" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image", "details": err.Error()})
        return
    }

    if err := h.DB.Create(&categoryData).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Category created successfully", "category": categoryData})
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category

	totalItems, page, totalPages, err := utils.PaginateAndSearch(c, h.DB, &models.Category{}, &categories, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"current_page":   page,
		"total_pages":    totalPages,
		"total_items":    totalItems,
		"items_per_page": c.DefaultQuery("limit", "10"),
		"categories":     categories,
	})
}


// GetCategory handles retrieving a category by its ID.
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := h.DB.First(&category, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve category"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID is required"})
		return
	}

	var categoryRequest dtos.CategoryRequest
	if err := c.ShouldBind(&categoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	var categoryData models.Category
	if err := h.DB.First(&categoryData, categoryID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve category", "details": err.Error()})
		return
	}

	categoryData.Name = categoryRequest.Name
	categoryData.Description = categoryRequest.Description

	file, err := c.FormFile("image")
	if err == nil {
		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file", "details": err.Error()})
			return
		}
		defer fileContent.Close() 

		imageURL, err := utils.UploadImageToCloudinary(fileContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image", "details": err.Error()})
			return
		}
		fmt.Printf("ImageUrl %s \n", imageURL)
		categoryData.ImageURL = imageURL
	}

	if err := h.DB.Save(&categoryData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully", "category": categoryData})
}


func (h *CategoryHandler) PatchCategory(c *gin.Context) {
    id := c.Param("id")
    var category models.Category

    if err := h.DB.First(&category, id).Error; err != nil {
        if gorm.IsRecordNotFoundError(err) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve category"})
        }
        return
    }

	var categoryRequest dtos.CategoryRequest
	if err := c.ShouldBind(&categoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	var categoryData models.Category
	if err := h.DB.First(&categoryData, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve category", "details": err.Error()})
		return
	}


    if categoryRequest.Name != "" {
        categoryData.Name = categoryRequest.Name
    }
    if categoryRequest.Description != "" {
        categoryData.Description = categoryRequest.Description
    }

    file, err := c.FormFile("image")
	if err == nil {
		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file", "details": err.Error()})
			return
		}
		defer fileContent.Close() 

		imageURL, err := utils.UploadImageToCloudinary(fileContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image", "details": err.Error()})
			return
		}
		categoryData.ImageURL = imageURL
	}

    if err := h.DB.Save(&categoryData).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully", "category": categoryData})
}



// DeleteCategory handles deleting a category by its ID.
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := h.DB.First(&category, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve category"})
		}
		return
	}

	if err := h.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
