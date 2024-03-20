package agents

import (
	"os"
	"path/filepath"
)

func loadControlPrimitivesContext(primitiveNames []string, packagePath string) ([]string, error) {
	var err error
	if len(primitiveNames) == 0 {
		primitiveNames, err = listPrimitiveNames(packagePath)
		if err != nil {
			return nil, err
		}
	}

	primitives := make([]string, 0)
	for _, primitiveName := range primitiveNames {
		content, err := loadText(filepath.Join(packagePath, "control_primitives_context", primitiveName+".js"))
		if err != nil {
			return nil, err
		}
		primitives = append(primitives, content)
	}
	return primitives, nil
}

func listPrimitiveNames(packagePath string) ([]string, error) {
	var primitiveNames []string
	primitiveFiles, err := os.ReadDir(filepath.Join(packagePath, "control_primitives_context"))
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range primitiveFiles {
		if fileInfo.IsDir() || filepath.Ext(fileInfo.Name()) != ".js" {
			continue
		}
		primitiveNames = append(primitiveNames, fileInfo.Name()[:len(fileInfo.Name())-3])
	}
	return primitiveNames, nil
}

func loadText(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}


