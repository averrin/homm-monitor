//go:generate goversioninfo -icon=../../resources/magic-ball.ico -manifest=./player.exe.manifest
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/html/charset"
	"golang.org/x/net/websocket"
	"golang.org/x/text/encoding/charmap"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"

	. "github.com/averrin/homm-monitor/pkg/homm"
	. "github.com/averrin/homm-monitor/pkg/structs"
)

var PREV_STATE *Report
var STATE *Report

var heroes map[int]*HeroStore

var mw *walk.MainWindow
var outTE *walk.TextEdit
var updateLabel *walk.Label
var clientLabel *walk.Label
var matchNameLabel *walk.Label
var matchCreatedLabel *walk.Label
var matchNameLabel2 *walk.Label
var matchCreatedLabel2 *walk.Label
var pluginAliveLabel *walk.Label
var matchIdInput *walk.LineEdit
var playerNameInput *walk.LineEdit
var connectButton *walk.PushButton
var session *r.Session
var connected = false
var versionMismatch = false
var isPluginAlive = false
var wasPluginAlive = false
var heartbeatTimer = time.NewTimer(10 * time.Second)

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
				STATE.IsPluginAlive = isPluginAlive
				state, err := json.Marshal(STATE)
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

func processReport() {
	if !versionMismatch && STATE.Version != PLUGIN_VERSION {
		walk.MsgBox(mw, "Plugin version mismatch",
			fmt.Sprintf("Please update plugin!\r\nCurrent version: %d\r\nDesired version: %d",
				STATE.Version, PLUGIN_VERSION), walk.MsgBoxOK)
		versionMismatch = true
	}

	STATE.Resources.GoldFound = PREV_STATE.Resources.GoldFound
	if STATE.Resources.Gold > PREV_STATE.Resources.Gold {
		STATE.Resources.GoldFound += STATE.Resources.Gold - PREV_STATE.Resources.Gold
	}
	if STATE.Resources.Gold < PREV_STATE.Resources.Gold {
		STATE.Resources.GoldSpent += PREV_STATE.Resources.Gold - STATE.Resources.Gold
	}

	STATE.Resources.WoodFound = PREV_STATE.Resources.WoodFound
	if STATE.Resources.Wood > PREV_STATE.Resources.Wood {
		STATE.Resources.WoodFound += STATE.Resources.Wood - PREV_STATE.Resources.Wood
	}
	if STATE.Resources.Wood < PREV_STATE.Resources.Wood {
		STATE.Resources.WoodSpent += PREV_STATE.Resources.Wood - STATE.Resources.Wood
	}

	STATE.Resources.MercuryFound = PREV_STATE.Resources.MercuryFound
	if STATE.Resources.Mercury > PREV_STATE.Resources.Mercury {
		STATE.Resources.MercuryFound += STATE.Resources.Mercury - PREV_STATE.Resources.Mercury
	}
	if STATE.Resources.Mercury < PREV_STATE.Resources.Mercury {
		STATE.Resources.MercurySpent += PREV_STATE.Resources.Mercury - STATE.Resources.Mercury
	}

	STATE.Resources.OreFound = PREV_STATE.Resources.OreFound
	if STATE.Resources.Ore > PREV_STATE.Resources.Ore {
		STATE.Resources.OreFound += STATE.Resources.Ore - PREV_STATE.Resources.Ore
	}
	if STATE.Resources.Ore < PREV_STATE.Resources.Ore {
		STATE.Resources.OreSpent += PREV_STATE.Resources.Ore - STATE.Resources.Ore
	}

	STATE.Resources.SulfurFound = PREV_STATE.Resources.SulfurFound
	if STATE.Resources.Sulfur > PREV_STATE.Resources.Sulfur {
		STATE.Resources.SulfurFound += STATE.Resources.Sulfur - PREV_STATE.Resources.Sulfur
	}
	if STATE.Resources.Sulfur < PREV_STATE.Resources.Sulfur {
		STATE.Resources.SulfurSpent += PREV_STATE.Resources.Sulfur - STATE.Resources.Sulfur
	}

	STATE.Resources.CrystalFound = PREV_STATE.Resources.CrystalFound
	if STATE.Resources.Crystal > PREV_STATE.Resources.Crystal {
		STATE.Resources.CrystalFound += STATE.Resources.Crystal - PREV_STATE.Resources.Crystal
	}
	if STATE.Resources.Crystal < PREV_STATE.Resources.Crystal {
		STATE.Resources.CrystalSpent += PREV_STATE.Resources.Crystal - STATE.Resources.Crystal
	}

	STATE.Resources.GemsFound = PREV_STATE.Resources.GemsFound
	if STATE.Resources.Gems > PREV_STATE.Resources.Gems {
		STATE.Resources.GemsFound += STATE.Resources.Gems - PREV_STATE.Resources.Gems
	}
	if STATE.Resources.Gems < PREV_STATE.Resources.Gems {
		STATE.Resources.GemsSpent += PREV_STATE.Resources.Gems - STATE.Resources.Gems
	}

	STATE.TotalExp = PREV_STATE.TotalExp
	STATE.TotalMPSpent = PREV_STATE.TotalMPSpent
	for _, h := range STATE.Heroes {
		println(h.Name)
		dec := charmap.Windows1251.NewDecoder()
		out, _ := dec.Bytes([]byte(h.Name))
		println(string(out))

		_, ok := heroes[h.Id]
		if !ok {
			heroes[h.Id] = &HeroStore{h.Movement, 0, h.Experience, h.Experience}
		}
		hs := heroes[h.Id]
		if STATE.Day != PREV_STATE.Day {
			hs.LastMP = h.Movement
		}
		if h.Movement < hs.LastMP {
			STATE.TotalMPSpent += hs.LastMP - h.Movement
			hs.MPSpent += hs.LastMP - h.Movement
			hs.LastMP = h.Movement
		}
		if h.Experience > hs.LastExp {
			STATE.TotalExp += h.Experience - hs.LastExp
			hs.LastExp = h.Experience
		}

		for _, spell := range h.LearnedSpells {
			if spell == FLY {
				STATE.HUD.HasFly = true
			}
			if spell == DIMENSION_DOOR {
				STATE.HUD.HasDD = true
			}
			if spell == TOWN_PORTAL {
				STATE.HUD.HasTP = true
			}
			if spell == RESURRECTION {
				STATE.HUD.HasResurrect = true
			}
			if spell == ARMAGEDDON {
				STATE.HUD.HasArmageddon = true
			}
		}

		for _, item := range h.Weared {
			if item.Id == ANGEL_WINGS {
				STATE.HUD.HasWings = true
			}
			if item.Id == SPELLBINDERS_HAT {
				STATE.HUD.HasSpellbindersHat = true
				STATE.HUD.HasDD = true
				STATE.HUD.HasFly = true
			}
			if item.Id == TOME_OF_AIR_MAGIC {
				STATE.HUD.HasTomeOfAir = true
				STATE.HUD.HasFly = true
				STATE.HUD.HasDD = true
			}
			if item.Id == TOME_OF_FIRE_MAGIC {
				STATE.HUD.HasTomeOfFire = true
				STATE.HUD.HasArmageddon = true
			}
			if item.Id == TOME_OF_EARTH_MAGIC {
				STATE.HUD.HasTomeOfEarth = true
				STATE.HUD.HasTP = true
				STATE.HUD.HasResurrect = true
			}
			if item.Id == ANGELIC_ALLIANCE {
				STATE.HUD.HasAlliance = true
			}
			if item.Id == SHACKLES_OF_WAR {
				STATE.HUD.HasShackles = true
			}
			if item.Id == ARMOR_OF_THE_DAMNED {
				STATE.HUD.HasAOTD = true
			}
			if item.Id == IRONFIST_OF_THE_OGRE {
				STATE.HUD.HasIronFist = true
			}

			if item.Id == SPELL_SCROLL {
				if item.Type == FLY {
					STATE.HUD.HasFly = true
				}
				if item.Type == DIMENSION_DOOR {
					STATE.HUD.HasDD = true
				}
				if item.Type == ARMAGEDDON {
					STATE.HUD.HasArmageddon = true
				}
				if item.Type == RESURRECTION {
					STATE.HUD.HasResurrect = true
				}
				if item.Type == TOWN_PORTAL {
					STATE.HUD.HasTP = true
				}
			}
		}
		for _, item := range h.Backpack {
			if item.Id == ANGEL_WINGS {
				STATE.HUD.HasWings = true
			}
			if item.Id == SPELLBINDERS_HAT {
				STATE.HUD.HasSpellbindersHat = true
				STATE.HUD.HasDD = true
				STATE.HUD.HasFly = true
			}
			if item.Id == TOME_OF_AIR_MAGIC {
				STATE.HUD.HasTomeOfAir = true
				STATE.HUD.HasFly = true
				STATE.HUD.HasDD = true
			}
			if item.Id == TOME_OF_FIRE_MAGIC {
				STATE.HUD.HasTomeOfFire = true
				STATE.HUD.HasArmageddon = true
			}
			if item.Id == TOME_OF_EARTH_MAGIC {
				STATE.HUD.HasTomeOfEarth = true
				STATE.HUD.HasTP = true
				STATE.HUD.HasResurrect = true
			}
			if item.Id == ANGELIC_ALLIANCE {
				STATE.HUD.HasAlliance = true
			}
			if item.Id == SHACKLES_OF_WAR {
				STATE.HUD.HasShackles = true
			}
			if item.Id == ARMOR_OF_THE_DAMNED {
				STATE.HUD.HasAOTD = true
			}
			if item.Id == IRONFIST_OF_THE_OGRE {
				STATE.HUD.HasIronFist = true
			}
			if item.Id == SPELL_SCROLL {
				if item.Type == FLY {
					STATE.HUD.HasFly = true
				}
				if item.Type == DIMENSION_DOOR {
					STATE.HUD.HasDD = true
				}
				if item.Type == ARMAGEDDON {
					STATE.HUD.HasArmageddon = true
				}
				if item.Type == RESURRECTION {
					STATE.HUD.HasResurrect = true
				}
				if item.Type == TOWN_PORTAL {
					STATE.HUD.HasTP = true
				}
			}
		}

		if h.SecSkills[EARTH_MAGIC] == 3 {
			STATE.HUD.HasExpertEarth = true
		}

		STATE.ArmyValue += h.ArmyValue
		for _, t := range STATE.Towns {
			fmt.Printf("%d == %d\n", t.GarrisonHero, h.Id)
			if t.GarrisonHero == h.Id {
				fmt.Println("in garr")
				h.InGarrison = true
			}
		}
	}

	STATE.TownsCountMax = int(math.Max(float64(STATE.TownsCount), float64(PREV_STATE.TownsCountMax)))

	for _, t := range STATE.Towns {
		if t.Grail {
			STATE.HUD.HasGrail = true
		}
		STATE.ArmyValue += t.GuardsValue
	}

	if PREV_STATE.StartHero == -1 && len(STATE.Heroes) > 0 {
		STATE.StartHero = STATE.Heroes[0].Id
	} else {
		STATE.StartHero = PREV_STATE.StartHero
	}
	if PREV_STATE.StartTown.Type == -1 && len(STATE.Towns) > 0 {
		STATE.StartTown = *STATE.Towns[0]
	} else {
		STATE.StartTown = PREV_STATE.StartTown
	}

	STATE.MaxAPM = int(math.Max(float64(PREV_STATE.MaxAPM), float64(STATE.APM)))
	STATE.IsPluginAlive = isPluginAlive

	STATE.UpdateTime = time.Now()
	STATE.ClientVersion = VERSION
	state, _ := json.MarshalIndent(STATE, "\r\n", "    ")
	if connected {
		STATE.MatchId = matchIdInput.Text()
		STATE.PlayerName = playerNameInput.Text()
		_, err := r.Table(STATE.MatchId).Insert(STATE).RunWrite(session)
		if err != nil {
			outTE.SetText(fmt.Sprintf("%v", err))
		} else {
			outTE.SetText(string(state))
		}
	} else {
		outTE.SetText(string(state))
	}

	updateLabel.SetText(fmt.Sprintf("Last update: %s", STATE.UpdateTime))
}

func resetState() {
	STATE = &Report{}
	STATE.StartHero = -1
	STATE.StartTown = Town{}
	STATE.StartTown.Type = -1
	PREV_STATE = STATE
	heroes = map[int]*HeroStore{}
}

func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

func runServer() {
	e := echo.New()

	resetState()

	go func() {
		for {
			<-heartbeatTimer.C
			isPluginAlive = false
			heartbeatTimer.Reset(10 * time.Second)
			if wasPluginAlive {
				pluginAliveLabel.SetText("Plugin not found")
				pluginAliveLabel.SetTextColor(walk.RGB(0xee, 0, 0))
				font := pluginAliveLabel.Font()
				newFont, _ := walk.NewFont(font.Family(), font.PointSize(), walk.FontBold)
				pluginAliveLabel.SetFont(newFont)
			}
		}

	}()

	e.POST("/heartbeat", func(c echo.Context) error {
		isPluginAlive = true
		wasPluginAlive = true
		heartbeatTimer.Reset(10 * time.Second)
		pluginAliveLabel.SetText("Plugin detected")
		pluginAliveLabel.SetTextColor(walk.RGB(0, 0xbb, 0))
		font := pluginAliveLabel.Font()
		newFont, _ := walk.NewFont(font.Family(), font.PointSize(), walk.FontBold)
		pluginAliveLabel.SetFont(newFont)
		return nil
	})

	e.POST("/report", func(c echo.Context) error {
		PREV_STATE = STATE
		STATE = new(Report)
		if err := c.Bind(STATE); err != nil {
			return c.String(http.StatusOK, "error")
		}
		c.JSON(http.StatusOK, STATE)
		isPluginAlive = true
		wasPluginAlive = true
		processReport()
		return nil
	})
	e.POST("/reset", func(c echo.Context) error {
		resetState()
		c.JSON(http.StatusOK, STATE)
		outTE.SetText("State reseted")
		return nil
	})
	e.Static("/", "./assets")
	e.GET("/ws", wsHandler)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HideBanner = true
	e.Logger.Fatal(e.Start("localhost:8989"))
}

func ConnectToMatch(matchId string, playerName string) {
	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address:  URL,
		Database: "homm_monitor_games",
	})
	if err != nil {
		connectButton.SetEnabled(false)
		log.Fatalln(err)
		return
	}

	cursor, err := r.Table("matches").Get(matchId).Run(session)
	if err != nil {
		connectButton.SetEnabled(false)
		log.Fatalln(err)
		return
	}
	match := Match{}
	cursor.One(&match)
	cursor.Close()

	matchNameLabel.SetText(fmt.Sprintf("Match name: %s", match.Name))
	matchCreatedLabel.SetText(fmt.Sprintf("Match created: %s", match.Created.Format("02.01.2006 15:04:05")))
	matchNameLabel2.SetText(fmt.Sprintf("Match name: %s", match.Name))
	matchCreatedLabel2.SetText(fmt.Sprintf("Match created: %s", match.Created.Format("02.01.2006 15:04:05")))

	STATE.PlayerName = playerName
	_, err = r.Table(matchId).Insert(STATE).RunWrite(session)
	if err != nil {
		connectButton.SetEnabled(false)
		log.Fatalln(err)
	}

	connected = true
	connectButton.SetEnabled(false)
	connectButton.SetText("Connected")
}

func main() {
	var appIcon, _ = walk.NewIconFromResourceId(2)
	mw_ := MainWindow{
		AssignTo: &mw,
		Icon:     appIcon,
		Title:    fmt.Sprintf("HoMM Monitor v%s (Player)", VERSION),
		Size:     Size{500, 600},
		Layout:   VBox{},
		Children: []Widget{
			TabWidget{
				Pages: []TabPage{
					TabPage{
						Title:  "Main",
						Layout: VBox{},
						Children: []Widget{
							Label{AssignTo: &pluginAliveLabel, Text: "Waiting for plugin...", Background: SolidColorBrush{Color: walk.RGB(0xff, 0xff, 0xff)}},
							Label{AssignTo: &updateLabel},
							Label{AssignTo: &clientLabel},
							Label{AssignTo: &matchNameLabel2},
							Label{AssignTo: &matchCreatedLabel2},
							ScrollView{
								Layout: VBox{MarginsZero: true},
								Children: []Widget{
									TextEdit{AssignTo: &outTE, ReadOnly: true, VScroll: true},
								},
							},
							PushButton{
								Text: "Reset",
								OnClicked: func() {
									resetState()
									outTE.SetText("State reseted")
								},
							},
						},
					},
					TabPage{
						Title:  "Match",
						Layout: VBox{},
						Children: []Widget{
							Label{Text: "Player name"},
							LineEdit{AssignTo: &playerNameInput, OnTextChanged: func() {
								connectButton.SetEnabled(matchIdInput.Text() != "" && playerNameInput.Text() != "")
							}},
							Label{Text: "Match ID. Take it from commenter."},
							LineEdit{AssignTo: &matchIdInput, OnTextChanged: func() {
								connectButton.SetEnabled(matchIdInput.Text() != "" && playerNameInput.Text() != "")
							}},
							PushButton{
								AssignTo: &connectButton,
								Text:     "Connect",
								OnClicked: func() {
									ConnectToMatch(matchIdInput.Text(), playerNameInput.Text())
								},
								Enabled: false,
							},
							Label{AssignTo: &matchNameLabel},
							Label{AssignTo: &matchCreatedLabel},
							VSpacer{},
						},
					},
				},
			},
		},
	}
	//mw.Create()
	go runServer()
	mw_.Run()
}
