package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
)

type Form struct {
	Title   string `form:"title" json:"title" xml:"title" binding:"required"`
	Content string `form:"content" json:"content" xml:"content" binding:"required"`
}

func main() {
	PORT := os.Getenv("PORT")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/message", func(c *gin.Context) {
		var form Form
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		emailFrom := os.Getenv("EMAIL_FROM")
		emailPassword := os.Getenv("EMAIL_PASSWORD")
		emailTo := []string{}
		emailTo = append(emailTo, os.Getenv("EMAIL_TO"))
		smtpHost := os.Getenv("SMTP_HOST")
		smtpPort := os.Getenv("SMTP_PORT")

		auth := smtp.PlainAuth("", emailFrom, emailPassword, smtpHost)

		t, _ := template.ParseFiles("template.html")

		var body bytes.Buffer

		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", form.Title, mimeHeaders)))
		t.Execute(&body, struct {
			Content string
		}{
			Content: form.Content,
		})

		fmt.Println("got upto just before email sending")
		fmt.Println(body.Bytes())
		address := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
		err := smtp.SendMail(address, auth, emailFrom, emailTo, body.Bytes())
		//err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", emailFrom, emailPassword, "smtp.gmail.com"), emailFrom, emailTo, body.Bytes())
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		fmt.Println("Got upto here too?")
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
	r.Run(":" + PORT)
	fmt.Printf("App is running on port %v\n", PORT)
}
