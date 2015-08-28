package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/simonchong/linny/common"
	"github.com/simonchong/linny/constants"

	"github.com/zenazn/goji/web"
)

func (f *Factory) AssetHTML() func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {

		fmt.Println("Requested: ", r.URL.Path[1:])
		// fmt.Println(r.Proto)
		// fmt.Println(r.Host)
		// fmt.Println(r.URL)

		fileReq := c.URLParams["file"]

		fileAbs, err := common.ResolveSecure(f.Conf.ContentRoot+"/"+constants.AssetsRoute, fileReq)
		if err != nil {
			fmt.Println("Secure Resolve Failed: ", err)
			http.NotFound(w, r)
			return
		}
		exists, fileAbs := common.FileExistsHTML(fileAbs)
		if !exists {
			fmt.Println("File Does Not Exist: ", fileAbs)
			http.NotFound(w, r)
			return
		}

		content, err := common.GetWrappedContent(fileAbs, f.Conf.ContentRoot)
		if err != nil {
			fmt.Println("Content Error: ", err)
			http.NotFound(w, r)
			return
		}
		content = common.InjectLinks(content, r)

		w.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprint(w, content)
	}
}

func (f *Factory) AssetFiles() http.Handler {

	absBaseDir, _ := filepath.Abs(f.Conf.ContentRoot)
	fileServeDir := absBaseDir + "/" + constants.AssetsRoute
	fmt.Println(fileServeDir)
	return http.StripPrefix("/"+constants.AssetsRoute+"/", http.FileServer(http.Dir(fileServeDir)))
}