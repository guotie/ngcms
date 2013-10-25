package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/robfig/revel"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

const MaxUpFileSize = (5 * 1024 * 1024)

var validFiletypes map[string]bool

func init() {
	validFiletypes = map[string]bool{".jpeg": true,
		".jpg": true,
		".png": true,
	}
}

func uploadpath() (string, string) {
	tm := time.Now()
	ymd := fmt.Sprintf("%d%d", tm.Year(), tm.Month())
	return revel.BasePath + "/public/upload/" + ymd, ymd
}

func validfile(fn string) (string, string, bool) {
	pos := strings.LastIndex(fn, ".")
	if pos == -1 {
		return fn, "", false
	}
	basename := fn[0:pos]
	ext := strings.ToLower(fn[pos:])
	_, ok := validFiletypes[ext]
	if !ok {
		return basename, ext, false
	}

	return basename, ext, true
}

func tmstr() string {
	tm := time.Now()
	return fmt.Sprintf("02d%02d%02d%02d", tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
}

func savefile(basename, ext string, content []byte) (string, error) {
	p, relp := uploadpath()
	tmpname := basename + tmstr()
	fname := base64.StdEncoding.EncodeToString([]byte(tmpname))
	if len(fname) >= 100 {
		fname = fname[0:100]
	}

	fullname := path.Join(p, fname+ext)
	fd, err := os.OpenFile(fullname, os.O_RDWR|os.O_CREATE, 0666)

	if os.IsExist(err) {
		// 为每个文件后缀加上日期小时分钟秒，因此，不应该出现重名的问题，如果出现，认为两个文件相同
		revel.ERROR.Printf("filename conflict: %s\n", tmpname)
		return fullname, nil
	} else if os.IsNotExist(err) {
		err := os.MkdirAll(p, 0666)
		if err != nil {
			revel.ERROR.Printf("mkdir failed: %v\n", err)
			return "", err
		}
		fd, err = os.OpenFile(fullname, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			revel.ERROR.Printf("openfile failed: %v\n", err)
			return "", err
		}
	} else if err != nil {
		revel.ERROR.Printf("openfile failed: %v\n", err)
		return "", err
	}

	defer fd.Close()
	_, err = fd.Write(content)
	if err != nil {
		revel.ERROR.Printf("write file failed: %v\n", err)
		return "", err
	}

	return path.Join(relp, fname+ext), nil
}

func (c App) Upload(editorid string) revel.Result {
	fmt.Println(c.Request.URL, c.Request.RequestURI)
	upfile := c.Params.Files["upfile"]
	if len(upfile) != 1 {
		return c.RenderError(fmt.Errorf("Only 1 file is accepted per upload.\n"))
	}

	upFileHeader := upfile[0]
	basename, ext, valid := validfile(upFileHeader.Filename)
	if !valid {
		return c.RenderError(fmt.Errorf("invalid upload filename or type: %s.\n", upFileHeader.Filename))
	}

	input, err := upFileHeader.Open()
	if err != nil {
		return c.RenderError(err)
	}

	upBytes, err := ioutil.ReadAll(input)
	input.Close()
	if err != nil || len(upBytes) == 0 {
		return c.RenderError(err)
	}
	if len(upBytes) >= MaxUpFileSize {
		return c.RenderError(fmt.Errorf("upfile size too large, cannot exceed 5MB."))
	}

	savedFile, err := savefile(basename, ext, upBytes)
	if err != nil {
		return c.RenderError(err)
	}

	reqtype := c.Params.Query["type"]
	if editorid == "" || (len(reqtype) > 0 && reqtype[0] == "ajax") {
		return c.RenderText(savedFile)
	}

	res := fmt.Sprintf(`<script>parent.UM.getEditor('%s').getWidgetCallback('image')('%s','SUCCESS')</script>`,
		editorid, savedFile)

	fmt.Println(res)

	return c.RenderText(res)
}
