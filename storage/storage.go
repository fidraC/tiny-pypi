package storage
import (
	"log"
	"os"
)
func init() {
	// Create repo folder if it doesn't exist
	if _, err := os.Stat("repo"); os.IsNotExist(err) {
		err = os.Mkdir("repo", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type Storage interface {
	// GetPackageLinks returns a list of links for the given package name and the HTTP status code
	GetPackageLinks(packageName string) ([]Link, int)
	GetIndex() ([]Link, error)
	// GetFile returns the file content if the sha256 hash matches
	GetFile(packageName, filename, hash string) ([]byte, int)
	// PutFile stores the file content and returns the sha256 hash
	PutFile(packageName, filename string, content []byte, hash string) error
}

type Link struct {
	Name string
	Url  string
}
