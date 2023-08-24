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
	go app.Schedule() // 啟動定時任務的 Goroutine
	
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", app.IndexView)

    http.HandleFunc("/login", app.Login)
    http.HandleFunc("/register", app.Register)
	http.HandleFunc("/index.html", app.IndexView)
	http.HandleFunc("/createAccount.html", app.CreateAccountView)

	http.HandleFunc("/mainPage.html", app.MainPageView)
	http.HandleFunc("/request1.html", app.Request1View)
    http.HandleFunc("/offer", app.Offer)
	http.HandleFunc("/track", app.TrackView)
	http.HandleFunc("/history", app.HistoryListView)

	
    // http.HandleFunc("/showAllUser", app.ShowAllUser)
    // http.HandleFunc("/showNowUser", app.ShowNowUser)
	// http.HandleFunc("/showAllOffer", app.ShowAllOffer)
    // http.HandleFunc("/match", app.Match)
	// http.HandleFunc("/showAllMatch", app.ShowAllMatch)
    // http.HandleFunc("/power", app.Power)
	// http.HandleFunc("/showAllPower", app.ShowAllPower)
	// http.HandleFunc("/showPowerbyCharger", app.ShowPowerbyCharger)


	fmt.Println("Start test WEB, port: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Start WEB Error")
	}

}
