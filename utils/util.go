package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)
// DumpText writes text data to a file.
func DumpText(data string, filePath string) error {
    // Create the directory if it doesn't exist
    err := os.MkdirAll(filepath.Dir(filePath), 0755)
    if err != nil {
        return err
    }

    // Open the file for writing
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write the data to the file
    _, err = file.WriteString(data)
    if err != nil {
        return err
    }

    return nil
}
// ListDir lists the contents of a directory
func ListDir(dirPath string) ([]os.FileInfo, error) {
	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the directory contents
	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	return files, nil
}
// LoadJSON loads the content of a JSON file into an interface{}.
func LoadJSON(filePaths ...string) (interface{}, error) {
    filePath := filepath.Join(filePaths...)
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var data interface{}
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&data)
    if err != nil {
        return nil, err
    }

    return data, nil
}

// ParseJSON parses a JSON string into an interface{}.
func ParseJSON(jsonStr string) (interface{}, error) {
    var data interface{}
    err := json.Unmarshal([]byte(jsonStr), &data)
    if err != nil {
        return nil, err
    }
    return data, nil
}
// Mkdir creates a directory at the specified path if it does not already exist.
func Mkdir(path string) error {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        err = os.MkdirAll(path, 0755) // 0755 is the default permission mode
        if err != nil {
            return err
        }
    }
    return nil
}
// DumpJSON dumps data into a JSON file.
func DumpJSON(data interface{}, filePath string) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    err = encoder.Encode(data)
    if err != nil {
        return err
    }
    return nil
}

// StringifyJSON returns a string representation of data.
func StringifyJSON(data interface{}) (string, error) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return "", err
    }
    return string(jsonData), nil
}

// CorrectJSON corrects common JSON errors.
func CorrectJSON(jsonStr string) (string, error) {
    // Remove tabs
    jsonStr = strings.ReplaceAll(jsonStr, "\t", "")

    // Parse JSON
    var parsedData interface{}
    err := json.Unmarshal([]byte(jsonStr), &parsedData)
    if err == nil {
        // JSON is valid, return as is
        return jsonStr, nil
    }

    // JSON is invalid, try to correct
    // Your correction logic here...

    return jsonStr, nil
}
