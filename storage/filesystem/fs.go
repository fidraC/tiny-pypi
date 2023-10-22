package filesystem

import (
	"errors"
	"fmt"
	"log"
	"os"
	"pypi/storage"
	"pypi/utilities"
	"strings"
)

// Implements the Storage interface
type StorageFs struct {
}

func NewStorageFs() *StorageFs {
	return &StorageFs{}
}

func (s *StorageFs) GetIndex() ([]storage.Link, error) {
	index := make([]storage.Link, 0)
	// List all files in repo directory
	folders, err := os.ReadDir("repo")
	if err != nil {
		return nil, err
	}
	for _, folder := range folders {
		// Check if folder is a directory
		if folder.IsDir() {
			index = append(index, storage.Link{
				Name: folder.Name(),
				Url:  "/simple/" + folder.Name() + "/", // Assume normalized when saved
			})
		}
	}
	return index, nil
}

func (s *StorageFs) GetPackageLinks(packageName string) ([]storage.Link, int) {
	packageName = utilities.Normalize(packageName)
	links := make([]storage.Link, 0)
	// List all files in repo directory
	files, err := os.ReadDir("repo/" + packageName)
	if err != nil {
		return nil, 404
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// Get hash for file (Should be in repo/<package>/hashes/<filename>.sha256)
		hash, err := os.ReadFile("repo/" + packageName + "/hashes/" + file.Name() + ".sha256")
		if err != nil {
			log.Println("Error reading hash file for", file.Name(), ":", err)
			continue
		}
		links = append(links, storage.Link{
			Name: file.Name(),
			Url:  fmt.Sprintf("/simple/%s/%s/%s#sha256=%s", packageName, strings.Trim(string(hash), "\n"), file.Name(), strings.Trim(string(hash), "\n")),
		})
	}
	return links, 200
}

func (s *StorageFs) GetFile(packageName, filename, hash string) ([]byte, int) {
	packageName = utilities.Normalize(packageName)
	filename = utilities.Normalize(filename)
	// Check if file exists
	if _, err := os.Stat("repo/" + packageName + "/" + filename); os.IsNotExist(err) {
		log.Println(err)
		return nil, 404
	}
	// Check if hash exists
	if _, err := os.Stat("repo/" + packageName + "/hashes/" + filename + ".sha256"); os.IsNotExist(err) {
		return nil, 500
	}
	// Read hash
	OriginalHash, err := os.ReadFile("repo/" + packageName + "/hashes/" + filename + ".sha256")
	if err != nil {
		return nil, 500
	}
	// Check if hash matches
	if strings.Trim(string(OriginalHash), "\n") != hash {
		return nil, 412 // Precondition failed
	}
	// Read file
	file, err := os.ReadFile("repo/" + packageName + "/" + filename)
	if err != nil {
		return nil, 500
	}
	return file, 200
}

func (s *StorageFs) PutFile(packageName, filename string, content []byte) (string, error) {
	// Normalize package name and filename
	packageName = utilities.Normalize(packageName)
	filename = utilities.Normalize(filename)
	// Create package folder if it doesn't exist
	if _, err := os.Stat("repo/" + packageName); os.IsNotExist(err) {
		err = os.Mkdir("repo/"+packageName, 0755)
		if err != nil {
			return "", err
		}
	}
	// Create hashes folder if it doesn't exist
	if _, err := os.Stat("repo/" + packageName + "/hashes"); os.IsNotExist(err) {
		err = os.Mkdir("repo/"+packageName+"/hashes", 0755)
		if err != nil {
			return "", err
		}
	}
	// Check if the file already exists
	if _, err := os.Stat("repo/" + packageName + "/" + filename); !os.IsNotExist(err) {
		return "", errors.New("file already exists")
	}
	// Write file
	err := os.WriteFile("repo/"+packageName+"/"+filename, content, 0644)
	if err != nil {
		return "", err
	}
	// Write hash
	hash := utilities.Hash(content)
	err = os.WriteFile("repo/"+packageName+"/hashes/"+filename+".sha256", []byte(hash), 0644)
	if err != nil {
		return "", err
	}
	return hash, nil
}
