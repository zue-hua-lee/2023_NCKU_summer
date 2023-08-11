/**
  author: jerry
 */
package controller

import (
	"net/http"
	"path/filepath"
	"html/template"
	"fmt"
)

func showView(w http.ResponseWriter, r *http.Request, templateName string, data interface{})  {
	page := filepath.Join("web", "tpl", templateName)

	// 創建模板
	resultTemplate, err := template.ParseFiles(page)
	if err != nil {
		fmt.Println("創建模板錯誤: ", err)
		return
	}

	// 資料融合
	err = resultTemplate.Execute(w, data)

	if err != nil {
		fmt.Println("資料融合錯誤: ", err)
		return
	}
}