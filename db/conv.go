package db

func GetString(r map[string]interface{}, key string) string {
	if v, exists := r[key]; exists {
		if _v, ok := v.(string); ok {
			return _v
		}
	}
	return ""
}

func GetInt(r map[string]interface{}, key string) int {
	if v, exists := r[key]; exists {
		if _v, ok := v.(int); ok {
			return _v
		}
	}
	return 0
}

func GetID(r map[string]interface{}) string {
	return GetString(r, ID)
}
