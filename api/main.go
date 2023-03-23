package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ShevonKuan/translate-server/module"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func translate(i string) string {
	// 调用翻译函数
	input := &module.InputObj{
		SourceText: i,
		SourceLang: "zh",
		TargetLang: "zh-TW",
	}
	output, _, err := module.Engine["google"](input)
	// 返回翻译结果
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return output.TransText
}

func modify(i *string, position string, content string) {
	*i, _ = sjson.Set(*i, position, content)
}

var (
	postion = []string{
		"editCommonBoxDetail.title",
		"editCommonBoxDetail.notesText",
	}
)

func Api(w http.ResponseWriter, r *http.Request) {
	posi := make(map[string]string) // position map[position]originText
	r.ParseForm()
	resp := `{
			"result": "success",
			"editCommonBoxDetail": ` + r.FormValue("editCommonBoxDetail") +
		`}`
	respJSON := gjson.Parse(resp)
	for _, v := range postion {
		posi[v] = respJSON.Get(v).String()
	}

	// add title
	posi["editCommonBoxDetail.title"] = respJSON.Get("editCommonBoxDetail.title").String()
	// add notesText
	posi["editCommonBoxDetail.notesText"] = respJSON.Get("editCommonBoxDetail.notesText").String()
	// add sizeMap
	respJSON.Get("editCommonBoxDetail.sizeMap").ForEach(func(key, value gjson.Result) bool {
		p := "editCommonBoxDetail.sizeMap." + key.String() + ".name"
		posi[p] = value.Get("name").String()
		return true
	})
	// add colorMap
	respJSON.Get("editCommonBoxDetail.colorMap").ForEach(func(key, value gjson.Result) bool {
		p := "editCommonBoxDetail.colorMap." + key.String() + ".name"
		posi[p] = value.Get("name").String()
		return true
	})
	// add sourceAttrs
	respJSON.Get("editCommonBoxDetail.sourceAttrs").ForEach(func(key, value gjson.Result) bool {
		p := "editCommonBoxDetail.sourceAttrs." + key.String() + ".name"
		posi[p] = value.Get("name").String()
		p = "editCommonBoxDetail.sourceAttrs." + key.String() + ".value"
		posi[p] = value.Get("value").String()
		return true
	})
	var wg sync.WaitGroup

	for k, v := range posi {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			if v != "" {
				modify(&resp, k, translate(v))
			}
		}(k, v)
	}
	wg.Wait()

	// return application/json
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resp))
}
