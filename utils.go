package kelly

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

// func JsonConfToStruct(path string, obj interface{}) error {
// 	file, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		return err
// 	}

// 	if err := json.Unmarshal(file, obj); err != nil {
// 		return err
// 	}
// 	return nil
// }
