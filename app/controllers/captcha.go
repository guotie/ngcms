package controllers

import (
	"fmt"
	"github.com/guotie/captcha"
	"github.com/robfig/revel"
	"image/png"
	"strings"
)

var _ = fmt.Printf

type CapImage struct {
	*captcha.Image
}

func (c App) Captcha(id string) revel.Result {
	pos := strings.LastIndex(id, ".")
	name := ""
	if pos == -1 {
		name = id
	} else {
		name = id[0:pos]
	}

	m := captcha.GetImage(name, captcha.StdWidth, captcha.StdHeight)
	return &CapImage{m}
}

func (img *CapImage) Apply(req *revel.Request, resp *revel.Response) {
	png.Encode(resp.Out, img.Paletted)
}
