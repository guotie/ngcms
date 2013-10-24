package controllers

import (
	"github.com/robfig/revel"
)

func (c App) NewTopic() revel.Result {
	return c.Render()
}

func (c App) PostNewTopic() revel.Result {
	return c.Render()
}
