package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type app struct {
	Title       string       `json:"title"`
	Version     string       `json:"version"`
	Maintainers []maintainer `json:"maintainers"`
	Company     string       `json:"company"`
	Website     string       `json:"website"`
	Source      string       `json:"source"`
	License     string       `json:"license"`
	Description string       `json:"description"`
}

var apps = []app{
	{Title: "Valid App 1",
		Version: "0.0.1",
		Maintainers: []maintainer{
			{Name: "Jonathan",
				Email: "Jonathan@Jonathan.com"},
		},
		Company:     "Company",
		Website:     "jonathan.com",
		Source:      "Jonathan Source",
		License:     "MIT",
		Description: "Jonathan Description",
	},
}

func main() {
	router := gin.Default()
	router.GET("/apps", getAppss)
	router.POST("/apps", postApps)

	router.Run("localhost:8080")
}

func getAppss(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, apps)
}

// title: Valid App 1
// version: 0.0.1
// maintainers:
// - name: firstmaintainer app1
//   email: firstmaintainer@hotmail.com
// - name: secondmaintainer app1
//   email: secondmaintainer@gmail.com
// company: Random Inc.
// website: https://website.com
// source: https://github.com/random/repo
// license: Apache-2.0
// description: |
//  ### Interesting Title
//  Some application content, and description

// curl --location --request POST 'http://localhost:8080/apps' \
// --header 'Content-Type: application/x-yaml' \
// --data-raw '
// title: My Awesome App
// version: 1.0.0
// maintainers:
//   - name: John Doe
//     email: johndoe@example.com
//   - name: Jane Doe
//     email: janedoe@example.com
//
// company: Acme Inc.
// website: https://www.acme.com/
// source: https://github.com/acme/my-awesome-app
// license: MIT
// description: My awesome app is awesome!
// '
func postApps(c *gin.Context) {
	var newApp app

	if err := c.ShouldBindYAML(&newApp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	apps = append(apps, newApp)

	c.IndentedJSON(http.StatusCreated, newApp)
}
