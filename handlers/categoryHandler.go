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


// GetCategories handles retrieving all categories.
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := h.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	if len(categories) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No categories found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
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

// UpdateCategory handles updating a category by its ID.
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
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

    var updatedCategory models.Category
    if err := c.ShouldBindJSON(&updatedCategory); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
        return
    }

    file, _, err := c.Request.FormFile("image")
    if err == nil {
        imageURL,err := utils.UploadImageToCloudinary(file)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image", "details": err.Error()})
            return
        }
        updatedCategory.ImageURL = imageURL 
    }

    if updatedCategory.Name != "" {
        category.Name = updatedCategory.Name
    }
    if updatedCategory.Description != "" {
        category.Description = updatedCategory.Description
    }
    if updatedCategory.ImageURL != "" {
        category.ImageURL = updatedCategory.ImageURL
    }

    if err := h.DB.Save(&category).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully", "category": category})
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

	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Only update fields that are provided (non-zero values)
	if updatedCategory.Name != "" {
		category.Name = updatedCategory.Name
	}
	if updatedCategory.Description != "" {
		category.Description = updatedCategory.Description
	}

	if err := h.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully", "category": category})
}