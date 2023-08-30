/*
*

	author: jerry
*/
package web

import (
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric/eep/web/controller"
)

func  WebStart(app *controller.Application)  {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/login.html", app.LoginView)
	http.HandleFunc("/createAccount.html", app.CreateAccountView)
	http.HandleFunc("/mainPage.html", app.MainPageView)
	http.HandleFunc("/request1.html", app.Request1View)
	http.HandleFunc("/request2.html", app.Request2View)
	http.HandleFunc("/request3.html", app.Request3View)
	http.HandleFunc("/request4.html", app.Request4View)
	http.HandleFunc("/requestView", app.RequestView)
	http.HandleFunc("/trackView", app.TrackView)
	http.HandleFunc("/historyView", app.HistoryView)

    http.HandleFunc("/register", app.Register)
    http.HandleFunc("/login", app.Login)
    http.HandleFunc("/offer", app.Offer)
    http.HandleFunc("/choice", app.Choice)
    http.HandleFunc("/match", app.Match)

	fmt.Println("Start test WEB, port: 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Start WEB Error")
	}
	go app.Schedule() // 啟動定時任務的 Goroutine

}
