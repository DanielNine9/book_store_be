package utils

import (
	"fmt"
	// "log"
	"reflect"
	// "strconv"
	"github.com/jinzhu/gorm"
	// "strings"
	// "shop-account/models" 
)
var modelPrefixes = map[string]string{
	"Category": "CA",
	"Book":     "BO",
	"Transaction": "TST",
	"Author":   "AU", 
	"Purchase":     "PC",
}
func GenerateCode(db *gorm.DB, model interface{}) (string, error) {
	// Get the actual model type name (e.g., "Author")
	modelType := reflect.TypeOf(model).Elem().Name()

	// Convert the model name to lowercase to match the table name in the database
	// (If needed, you can use modelType for more complex mappings to the database table name)
	// modelName := strings.ToLower(modelType)

	prefix, exists := modelPrefixes[modelType]
	if !exists {
		return "", fmt.Errorf("unknown model type: %s", modelType)
	}

	var count int64
	err := db.Unscoped().Model(model).Where("id IS NOT NULL").Count(&count).Error
	if err != nil {
		return "", fmt.Errorf("failed to count records for model %s: %v", modelType, err)
	}

	code := fmt.Sprintf("%s%02d", prefix, count+1)
	return code, nil
}
