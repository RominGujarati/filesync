package file

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"path/filepath"
	"os"
)

func GetAllFiles(c *gin.Context) {
	var fileList []map[string]interface{}

	// find all the file with sub directory recursively
	var walkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".yaml" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return nil // skip files we can't read
			}
			relPath, _ := filepath.Rel("./config", path)
			fileList = append(fileList, map[string]interface{}{
				"fileName": relPath,
				"fileType": "YAML",
				"content":  string(content),
			})
		}
		return nil
	}

	filepath.Walk("./config", walkFunc)

	c.JSON(http.StatusOK, gin.H{"files": fileList})
}

func GetFileContent(c *gin.Context) {
	fileName := c.Param("fileName")
	content, err := ioutil.ReadFile(filepath.Join("./config", fileName))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": string(content)})
}