//go:generate goversioninfo -icon=../../resources/crystal-ball.ico -manifest=./commentator.exe.manifest
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"

	. "github.com/averrin/homm-monitor/pkg/etc"
	. "github.com/averrin/homm-monitor/pkg/structs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var session *r.Session
var matchIdInput *walk.LineEdit
var loadMatchIdInput *walk.LineEdit
var matchNameInput *walk.LineEdit
var overlayMatchDescInput *walk.LineEdit
var overlayMatchNameInput *walk.LineEdit
var pb *walk.PushButton
var loadMatchPB *walk.PushButton
var outTE *walk.TextEdit
var clientLabel *walk.Label
var lastUpdateLabel *walk.Label
var started = false

var REPORTS map[int]Report

func listenChanges(id string) {
	res, err := r.Table(id).Changes().Run(session)
	if err != nil {
		outTE.SetText("Something went wrong")
		return
	}
	var value Update
	for res.Next(&value) {
		value.NewValue.RecTime = time.Now()
		REPORTS[value.NewValue.Color] = value.NewValue
		update := ""
		for i := 0; i < 8; i++ {
			if report, ok := REPORTS[i]; ok {
				color := []string{"Red", "Blue"}[report.Color]
				update += fmt.Sprintf("Last update: %s from %s (%s)\r\n", report.RecTime, report.PlayerName, color)
				if report.ClientVersion != VERSION {
					update += fmt.Sprintf("%s version mismatch detected! %s != %s\r\n", report.PlayerName, report.ClientVersion, VERSION)
				}
			}
		}
		outTE.SetText(update)
	}
}

func LoadMatch(id string) {
	cursor, err := r.Table("matches").Get(id).Run(session)
	if err != nil {
		log.Fatalln(err)
	}
	match := Match{}
	cursor.One(&match)
	cursor.Close()
	started = true
	overlayMatchNameInput.SetText(match.Name)

	//.OrderBy(r.Desc("updateTime"))
	report := Report{}
	err = r.Table(match.Id).Limit(1).ReadOne(&report, session)
	if err != nil {
		log.Fatalln(err)
	}

	cursor.One(&report)
	cursor.Close()
	REPORTS[report.Color] = report

	//state, _ := json.MarshalIndent(report, "\r\n", "    ")
	//outTE.SetText(string(state))
	color := []string{"Red", "Blue"}[report.Color]
	update := fmt.Sprintf("Last update: %s from %s (%s)\r\n", report.UpdateTime, report.PlayerName, color)
	lastUpdateLabel.SetText(update)
	outTE.SetText(update)

	go listenChanges(match.Id)
}

func CreateMatch(name string) string {
	gen := rand.New(rand.New(rand.NewSource(99)))
	gen.Seed(time.Now().UTC().UnixNano())
	table := fmt.Sprintf("%v-%v-%v", ADJECTIVES[gen.Intn(len(ADJECTIVES))], ADJECTIVES[gen.Intn(len(ADJECTIVES))], NOUNS[gen.Intn(len(NOUNS))])

	r.DB("homm_monitor_games").TableDrop(table).Exec(session)
	err := r.DB("homm_monitor_games").TableCreate(table).Exec(session)
	if err != nil {
		log.Fatalln(err)
		return "Something went wrong"
	}
	match := Match{
		table,
		name,
		time.Now(),
		-1,
	}
	r.Table("matches").Insert(match).RunWrite(session)
	started = true
	go listenChanges(table)

	return table
}

func wsHandler(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Write
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				//c.Logger().Error(err)
			}
			if msg == "update" {
				msg := CommenterState{REPORTS, overlayMatchNameInput.Text(), overlayMatchDescInput.Text()}
				state, err := json.Marshal(msg)
				if err != nil {
					panic(err)
				}
				err = websocket.Message.Send(ws, string(state))
				if err != nil {
					c.Logger().Error(err)
					panic(err)
				}
				clientLabel.SetText(fmt.Sprintf("Last overlay request: %s", time.Now().Format("15:04:05")))
			}
			time.Sleep(500 * time.Millisecond)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func runServer() {
	e := echo.New()
	e.Static("/", "./public")
	e.GET("/ws", wsHandler)

	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.HideBanner = true
	e.Logger.Fatal(e.Start("localhost:8988"))
}

func main() {
	REPORTS = map[int]Report{}
	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address:  URL,
		Database: "homm_monitor_games",
	})
	if err != nil {
		log.Fatalln(err)
	}

	var appIcon, _ = walk.NewIconFromResourceId(2)
	mw := MainWindow{
		Icon:   appIcon,
		Title:  fmt.Sprintf("HoMM Monitor v%s (Commentator)", VERSION),
		Size:   Size{500, 400},
		Layout: VBox{},
		Children: []Widget{
			TabWidget{
				Pages: []TabPage{
					TabPage{
						Title:  "New Match",
						Layout: VBox{},
						Children: []Widget{
							Label{Text: "Match Name"},
							LineEdit{AssignTo: &matchNameInput, OnTextChanged: func() {
								pb.SetEnabled(matchNameInput.Text() != "")
							}},
							PushButton{
								AssignTo: &pb,
								Text:     "Create Match",
								OnClicked: func() {
									matchIdInput.SetText(CreateMatch(matchNameInput.Text()))
									overlayMatchNameInput.SetText(matchNameInput.Text())
									pb.SetEnabled(!started)
								},
								Enabled: false,
							},
							Label{Text: "Match ID. Send it to players."},
							LineEdit{AssignTo: &matchIdInput, ReadOnly: true},
							VSpacer{},
						},
					},
					TabPage{
						Title:  "Load Match",
						Layout: VBox{},
						Children: []Widget{
							Label{Text: "Existing Match ID."},
							LineEdit{AssignTo: &loadMatchIdInput, OnTextChanged: func() {
								loadMatchPB.SetEnabled(loadMatchIdInput.Text() != "")
							}},
							PushButton{
								AssignTo: &loadMatchPB,
								Text:     "Load Match",
								OnClicked: func() {
									LoadMatch(loadMatchIdInput.Text())
									pb.SetEnabled(!started)
									loadMatchPB.SetEnabled(!started)
								},
								Enabled: false,
							},
							Label{AssignTo: &lastUpdateLabel},
							VSpacer{},
						},
					},
					TabPage{
						Title:  "Overlay settings",
						Layout: VBox{},
						Children: []Widget{
							Label{Text: "Displayed Match Name"},
							LineEdit{AssignTo: &overlayMatchNameInput},
							Label{Text: "Displayed Match Description"},
							LineEdit{AssignTo: &overlayMatchDescInput},
							VSpacer{},
						},
					},
					TabPage{
						Title:  "Status",
						Layout: VBox{},
						Children: []Widget{
							Label{AssignTo: &clientLabel},
							ScrollView{
								Layout: VBox{MarginsZero: true},
								Children: []Widget{
									TextEdit{AssignTo: &outTE, ReadOnly: true, VScroll: true},
								},
							},
						},
					},
				},
			},
		},
	}
	go runServer()
	mw.Run()
}
