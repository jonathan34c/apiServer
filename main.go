package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fmt"
	"strings"
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

// .net.http
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
	// router := gin.Default()
	// router.GET("/apps", getAppss)
	// router.POST("/apps", postApps)

	// router.Run(":8080")
	test()
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

func test() {
	// Example string with multiple matches
	text := "D8sV3 D16sV3 D32sV3 D8sV4 D16sV4 D32sV4 D8sV5 D16sV5 D32sV5 D8asV4 D16asV4 D32asV4 D8asV5 D16asV5 D32asV5 E8sV3 E16sV3 E32sV3 E8sV4 E16sV4 E20sV4 E32sV4 E48sV4 E64sV4 E8sV5 E16sV5 E20sV5 E32sV5 E48sV5 E64sV5 E96sV5 E4asV4 E8asV4 E16asV4 E20asV4 E32asV4 E48asV4 E64asV4 E96asV4 E8asV5 E16asV5 E20asV5 E32asV5 E48asV5 E64asV5 E96asV5 E64isV3 E80isV4 E80idsV4 E104isV5 E104idsV5 F72sV2 M128ms D4sV3 D8sV3 D16sV3 D32sV3 D4sV4 D8sV4 D16sV4 D32sV4 D64sV4 D96sV4 D4sV5 D8sV5 D16sV5 D32sV5 D64sV5 D96sV5 D4asV4 D8asV4 D16asV4 D32asV4 D64asV4 D96asV4 D4asV5 D8asV5 D16asV5 D32asV5 D64asV5 D96asV5 E4sV3 E8sV3 E16sV3 E32sV3 E2sV4 E4sV4 E8sV4 E16sV4 E20sV4 E32sV4 E48sV4 E64sV4 E96sV4 E2sV5 E4sV5 E8sV5 E16sV5 E20sV5 E32sV5 E48sV5 E64sV5 E96sV5 E4asV4 E8asV4 E16asV4 E20asV4 E32asV4 E48asV4 E64asV4 E96asV4 E8asV5 E16asV5 E20asV5 E32asV5 E48asV5 E64asV5 E96asV5 E64isV3 E80isV4 E80idsV4 E104isV5 E104idsV5 F4sV2 F8sV2 F16sV2 F32sV2 F72sV2 M128ms L4s L8s L16s L32s L8sV2 L16sV2 L32sV2 L48sV2 L64sV2 L8sV3 L16sV3 L32sV3 L48sV3 L64sV3 NC4asT4V3 NC8asT4V3 NC16asT4V3 NC64asT4V3 NC6sV3 NC12sV3 NC24sV3 NC24rsV3 "
	text = strings.ToLower(text)
	// Regular expression

	fmt.Println(text)
}
