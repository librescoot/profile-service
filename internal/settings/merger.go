package settings

func Merge(global, userOverrides map[string]string) map[string]string {
	result := make(map[string]string, len(global))
	for k, v := range global {
		result[k] = v
	}

	if userOverrides == nil {
		return result
	}

	for k, v := range userOverrides {
		if IsUserOverridable(k) {
			result[k] = v
		}
	}

	return result
}
