package main

import (
	"io"
	"net/http"
	"pypi/config"
	"pypi/storage"
	"pypi/storage/filesystem"

	"html/template"

	"github.com/gin-gonic/gin"
)

const indexHTML = `
<!DOCTYPE html>
<html>
	<body>
		{{range .}}
		<a href="{{.Url}}">{{.Name}}</a><br>
		{{end}}
	</body>
</html>
`

var indexTemplate = template.Must(template.New("index").Parse(indexHTML))

func main() {
	var store storage.Storage
	if config.StorageType == "filesystem" {
		store = filesystem.NewStorageFs()
	} else {
		panic("other storage types not supported yet")
	}

	r := gin.Default()
	r.SetHTMLTemplate(indexTemplate)

	r.GET("/simple/", func(c *gin.Context) {
		// Get link index
		links, err := store.GetIndex()
		if err != nil {
			c.String(500, "Error getting links")
			return
		}
		c.HTML(200, "index", links)
	})

	r.GET("/simple/:package/", func(c *gin.Context) {
		pkg := c.Param("package")
		// Get package links
		links, err := store.GetPackageLinks(pkg)
		if err != 200 {
			c.String(err, http.StatusText(err))
			return
		}
		c.HTML(200, "index", links)
	})

	r.GET("/simple/:package/:hash/:filename/", func(c *gin.Context) {
		packageName := c.Param("package")
		filename := c.Param("filename")
		hash := c.Param("hash")
		file, err := store.GetFile(packageName, filename, hash)
		if err != 200 {
			c.String(err, http.StatusText(err))
			return
		}
		c.Data(200, "application/octet-stream", file)
	})

	r.POST("/", func(c *gin.Context) {
		file, err := c.FormFile("content")
		if err != nil {
			c.String(500, "Error getting file")
			return
		}
		fileContent, err := file.Open()
		if err != nil {
			c.String(500, "Error opening file")
			return
		}
		fileData, err := io.ReadAll(fileContent)
		if err != nil {
			c.String(500, "Error reading file")
			return
		}
		sha256Digest := c.PostForm("sha256_digest")
		packageName := c.PostForm("name")
		// Save file
		err = store.PutFile(packageName,file.Filename, fileData, sha256Digest)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.String(200, "OK")
		
	})

	r.Run()
}
