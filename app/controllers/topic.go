package controllers

import (
	"fmt"
	"github.com/guotie/captcha"
	"github.com/robfig/revel"
)

var _ = fmt.Printf

func (c App) NewTopic() revel.Result {
	captchaId := captcha.New()

	return c.Render(captchaId)
}

func (c App) PostNewTopic() revel.Result {
	return c.Render()
}
