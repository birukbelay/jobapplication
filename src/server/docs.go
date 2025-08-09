package server

import (
	"io/fs"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

var EmbeddedAssets fs.FS

func ServFiberDoc(c *fiber.Ctx) error {
	c.Response().Header.Set("Content-Type", "text/html")
	c.Write([]byte(`<!doctype html>
			<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="referrer" content="same-origin" />
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
    <title>Docs Example reference</title>
    <!-- Embed elements Elements via Web Component -->
    <link href="/static/styles.min.css" rel="stylesheet" />
    <script src="/static/web-components.min.js"
            
            crossorigin="anonymous"></script>
  </head>
  <body style="height: 100vh;">
    <elements-api
      apiDescriptionUrl="/openapi.yaml"
      router="hash"
      layout="sidebar"
      tryItCredentialsPolicy="same-origin"
    />
  </body>
			</html>`))
	return nil
}

func ServGinDoc(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.Write([]byte(`<!doctype html>
			<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="referrer" content="same-origin" />
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
    <title>Docs Example reference</title>
    <!-- Embed elements Elements via Web Component -->
    <link href="/static/styles.min.css" rel="stylesheet" />
    <script src="/static/web-components.min.js"
            
            crossorigin="anonymous"></script>
  </head>
  <body style="height: 100vh;">
    <elements-api
      apiDescriptionUrl="/openapi.yaml"
      router="hash"
      layout="sidebar"
      tryItCredentialsPolicy="same-origin"
    />
  </body>
			</html>`))
}
