package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/reizt/ebra/middlewares"
	"github.com/reizt/ebra/renderings"
)

func SendTestMail(c echo.Context) error {
	rootDir := os.Getenv("WORKDIR")
	b, err := ioutil.ReadFile(rootDir + "/mails/test.html")
	if err != nil {
		return err
	}
	mailFormat := string(b)
	tmpl, err := template.New("testEmail").Parse(mailFormat)
	if err != nil {
		return err
	}
	type TestInput struct {
		Title     string
		Subtitle  string
		Href      string
		LinkTitle string
	}
	testInput := TestInput{
		Title:     "Foo",
		Subtitle:  "Bar",
		Href:      "golang.org",
		LinkTitle: "Let's start go",
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, testInput)
	if err != nil {
		return err
	}
	fmt.Println(tpl.String())
	m := new(middlewares.Mailer)
	input := &middlewares.SendMailInput{
		From:    os.Getenv("ADMIN_EMAIL"),
		To:      os.Getenv("ADMIN_EMAIL_SUB"),
		Subject: "Test Mail",
		Body:    tpl.String(),
	}
	messageId, sendErr := m.SendMail(input)
	if sendErr != nil {
		fmt.Println(sendErr.Error())
		return sendErr
	} else {
		return c.JSON(http.StatusCreated, &renderings.MessageIdResponse{
			MessasgeId: messageId,
		})
	}
}
