package handler

import (
	"fmt"
	"net/http"

	"github.com/arata-nvm/offblast/domain"
	"github.com/labstack/echo/v4"
)

type DetectJob struct {
	Text string `json:"text"`
}

type DetectResult struct {
	Haikus []string `json:"haikus"`
}

type RandomResult struct {
	Name   string   `json:"name"`
	Url    string   `json:"url"`
	Haikus []string `json:"haikus"`
}

func Detect(ctx echo.Context) error {
	job := new(DetectJob)
	if err := ctx.Bind(job); err != nil {
		return err
	}

	haikus, err := domain.FindHaikus(job.Text)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, DetectResult{
		Haikus: haikus,
	})
}

func Random(ctx echo.Context) error {
	law, err := domain.GetRandomLaw()
	if err != nil {
		return err
	}

	haikus, err := domain.FindHaikus(law.Body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://elaws.e-gov.go.jp/document?lawid=%s", law.Info.ID)

	result := RandomResult{
		Name:   law.Info.Name,
		Url:    url,
		Haikus: haikus,
	}
	domain.PostToSlack(result)
	return ctx.JSON(http.StatusOK, result)
}
