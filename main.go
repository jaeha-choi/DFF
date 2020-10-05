package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/anaskhan96/soup"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const ProjectName string = "AutoRunes"
const Version string = "v0.2"

var cli *http.Client
var config Config

var values []string

//var interval = 3
//var debug = false

type Config struct {
	Debug      bool
	Interval   float64
	ClientDir  string
	EnableRune bool
	EnableItem bool
}

type Item struct {
	Count int    `json:"count"`
	ID    string `json:"id"`
}

type ItemBlock struct {
	HideIfSummonerSpell string `json:"hideIfSummonerSpell"`
	Items               []Item `json:"items"`
	ShowIfSummonerSpell string `json:"showIfSummonerSpell"`
	Type                string `json:"type"`
}

type ItemSet struct {
	AssociatedChampions []int         `json:"associatedChampions"`
	AssociatedMaps      []int         `json:"associatedMaps"`
	Blocks              []ItemBlock   `json:"blocks"`
	Map                 string        `json:"map"`
	Mode                string        `json:"mode"`
	PreferredItemSlots  []interface{} `json:"preferredItemSlots"`
	Sortrank            int           `json:"sortrank"`
	StartedFrom         string        `json:"startedFrom"`
	Title               string        `json:"title"`
	Type                string        `json:"type"`
	UID                 string        `json:"uid"`
}

type ItemPage struct {
	AccountID int64     `json:"accountId"`
	ItemSets  []ItemSet `json:"itemSets"`
	Timestamp int64     `json:"timestamp"`
}

type RunePage struct {
	AutoModifiedSelections []interface{} `json:"autoModifiedSelections"`
	Current                bool          `json:"current"`
	ID                     int           `json:"id"`
	IsActive               bool          `json:"isActive"`
	IsDeletable            bool          `json:"isDeletable"`
	IsEditable             bool          `json:"isEditable"`
	IsValid                bool          `json:"isValid"`
	LastModified           int64         `json:"lastModified"`
	Name                   string        `json:"name"`
	Order                  int           `json:"order"`
	PrimaryStyleID         int           `json:"primaryStyleId"`
	SelectedPerkIds        []int         `json:"selectedPerkIds"`
	SubStyleID             int           `json:"subStyleId"`
}

type RunePages []struct {
	AutoModifiedSelections []interface{} `json:"autoModifiedSelections"`
	Current                bool          `json:"current"`
	ID                     int           `json:"id"`
	IsActive               bool          `json:"isActive"`
	IsDeletable            bool          `json:"isDeletable"`
	IsEditable             bool          `json:"isEditable"`
	IsValid                bool          `json:"isValid"`
	LastModified           int64         `json:"lastModified"`
	Name                   string        `json:"name"`
	Order                  int           `json:"order"`
	PrimaryStyleID         int           `json:"primaryStyleId"`
	SelectedPerkIds        []int         `json:"selectedPerkIds"`
	SubStyleID             int           `json:"subStyleId"`
}

type QueueInfo struct {
	CanInviteOthersAtEog bool `json:"canInviteOthersAtEog"`
	CurrentLobbyStatus   struct {
		AllowedPlayAgain      bool          `json:"allowedPlayAgain"`
		CustomSpectatorPolicy string        `json:"customSpectatorPolicy"`
		InvitedSummonerIds    []interface{} `json:"invitedSummonerIds"`
		IsCustom              bool          `json:"isCustom"`
		IsLeader              bool          `json:"isLeader"`
		IsPracticeTool        bool          `json:"isPracticeTool"`
		IsSpectator           bool          `json:"isSpectator"`
		LobbyID               string        `json:"lobbyId"`
		MemberSummonerIds     []int         `json:"memberSummonerIds"`
		QueueID               int           `json:"queueId"`
	} `json:"currentLobbyStatus"`
	LastQueuedLobbyStatus struct {
		AllowedPlayAgain      bool          `json:"allowedPlayAgain"`
		CustomSpectatorPolicy string        `json:"customSpectatorPolicy"`
		InvitedSummonerIds    []interface{} `json:"invitedSummonerIds"`
		IsCustom              bool          `json:"isCustom"`
		IsLeader              bool          `json:"isLeader"`
		IsPracticeTool        bool          `json:"isPracticeTool"`
		IsSpectator           bool          `json:"isSpectator"`
		LobbyID               string        `json:"lobbyId"`
		MemberSummonerIds     []int         `json:"memberSummonerIds"`
		QueueID               int           `json:"queueId"`
	} `json:"lastQueuedLobbyStatus"`
}

type Champion struct {
	Active             bool          `json:"active"`
	Alias              string        `json:"alias"`
	BanVoPath          string        `json:"banVoPath"`
	BaseLoadScreenPath string        `json:"baseLoadScreenPath"`
	BotEnabled         bool          `json:"botEnabled"`
	ChooseVoPath       string        `json:"chooseVoPath"`
	DisabledQueues     []interface{} `json:"disabledQueues"`
	FreeToPlay         bool          `json:"freeToPlay"`
	ID                 int           `json:"id"`
	Name               string        `json:"name"`
	Ownership          struct {
		FreeToPlayReward bool `json:"freeToPlayReward"`
		Owned            bool `json:"owned"`
		Rental           struct {
			EndDate           int  `json:"endDate"`
			PurchaseDate      int  `json:"purchaseDate"`
			Rented            bool `json:"rented"`
			WinCountRemaining int  `json:"winCountRemaining"`
		} `json:"rental"`
	} `json:"ownership"`
	Passive struct {
		Description string `json:"description"`
		Name        string `json:"name"`
	} `json:"passive"`
	Purchased         int      `json:"purchased"`
	RankedPlayEnabled bool     `json:"rankedPlayEnabled"`
	Roles             []string `json:"roles"`
	Skins             []struct {
		ChampionID int    `json:"championId"`
		ChromaPath string `json:"chromaPath"`
		Chromas    []struct {
			ChampionID   int      `json:"championId"`
			ChromaPath   string   `json:"chromaPath"`
			Colors       []string `json:"colors"`
			Disabled     bool     `json:"disabled"`
			ID           int      `json:"id"`
			LastSelected bool     `json:"lastSelected"`
			Name         string   `json:"name"`
			Ownership    struct {
				FreeToPlayReward bool `json:"freeToPlayReward"`
				Owned            bool `json:"owned"`
				Rental           struct {
					EndDate           int  `json:"endDate"`
					PurchaseDate      int  `json:"purchaseDate"`
					Rented            bool `json:"rented"`
					WinCountRemaining int  `json:"winCountRemaining"`
				} `json:"rental"`
			} `json:"ownership"`
			StillObtainable bool `json:"stillObtainable"`
		} `json:"chromas"`
		CollectionSplashVideoPath interface{}   `json:"collectionSplashVideoPath"`
		Disabled                  bool          `json:"disabled"`
		Emblems                   []interface{} `json:"emblems"`
		FeaturesText              interface{}   `json:"featuresText"`
		ID                        int           `json:"id"`
		IsBase                    bool          `json:"isBase"`
		LastSelected              bool          `json:"lastSelected"`
		LoadScreenPath            string        `json:"loadScreenPath"`
		Name                      string        `json:"name"`
		Ownership                 struct {
			FreeToPlayReward bool `json:"freeToPlayReward"`
			Owned            bool `json:"owned"`
			Rental           struct {
				EndDate           int  `json:"endDate"`
				PurchaseDate      int  `json:"purchaseDate"`
				Rented            bool `json:"rented"`
				WinCountRemaining int  `json:"winCountRemaining"`
			} `json:"rental"`
		} `json:"ownership"`
		QuestSkinInfo struct {
			CollectionsCardPath    string        `json:"collectionsCardPath"`
			CollectionsDescription string        `json:"collectionsDescription"`
			DescriptionInfo        []interface{} `json:"descriptionInfo"`
			Name                   string        `json:"name"`
			SplashPath             string        `json:"splashPath"`
			Tiers                  []interface{} `json:"tiers"`
			TilePath               string        `json:"tilePath"`
			UncenteredSplashPath   string        `json:"uncenteredSplashPath"`
		} `json:"questSkinInfo"`
		RarityGemPath        string      `json:"rarityGemPath"`
		SkinType             string      `json:"skinType"`
		SplashPath           string      `json:"splashPath"`
		SplashVideoPath      interface{} `json:"splashVideoPath"`
		StillObtainable      bool        `json:"stillObtainable"`
		TilePath             string      `json:"tilePath"`
		UncenteredSplashPath string      `json:"uncenteredSplashPath"`
	} `json:"skins"`
	Spells []struct {
		Description string `json:"description"`
		Name        string `json:"name"`
	} `json:"spells"`
	SquarePortraitPath string `json:"squarePortraitPath"`
	StingerSfxPath     string `json:"stingerSfxPath"`
	TacticalInfo       struct {
		DamageType string `json:"damageType"`
		Difficulty int    `json:"difficulty"`
		Style      int    `json:"style"`
	} `json:"tacticalInfo"`
	Title string `json:"title"`
}

type AccountInfo struct {
	AccountID                   int64  `json:"accountId"`
	DisplayName                 string `json:"displayName"`
	InternalName                string `json:"internalName"`
	NameChangeFlag              bool   `json:"nameChangeFlag"`
	PercentCompleteForNextLevel int    `json:"percentCompleteForNextLevel"`
	ProfileIconID               int    `json:"profileIconId"`
	Puuid                       string `json:"puuid"`
	RerollPoints                struct {
		CurrentPoints    int `json:"currentPoints"`
		MaxRolls         int `json:"maxRolls"`
		NumberOfRolls    int `json:"numberOfRolls"`
		PointsCostToRoll int `json:"pointsCostToRoll"`
		PointsToReroll   int `json:"pointsToReroll"`
	} `json:"rerollPoints"`
	SummonerID       int  `json:"summonerId"`
	SummonerLevel    int  `json:"summonerLevel"`
	Unnamed          bool `json:"unnamed"`
	XpSinceLastLevel int  `json:"xpSinceLastLevel"`
	XpUntilNextLevel int  `json:"xpUntilNextLevel"`
}

type ChampSelect struct {
	Actions [][]struct {
		ActorCellID  int    `json:"actorCellId"`
		ChampionID   int    `json:"championId"`
		Completed    bool   `json:"completed"`
		ID           int    `json:"id"`
		IsAllyAction bool   `json:"isAllyAction"`
		IsInProgress bool   `json:"isInProgress"`
		Type         string `json:"type"`
	} `json:"actions"`
	AllowBattleBoost    bool `json:"allowBattleBoost"`
	AllowDuplicatePicks bool `json:"allowDuplicatePicks"`
	AllowLockedEvents   bool `json:"allowLockedEvents"`
	AllowRerolling      bool `json:"allowRerolling"`
	AllowSkinSelection  bool `json:"allowSkinSelection"`
	Bans                struct {
		MyTeamBans    []interface{} `json:"myTeamBans"`
		NumBans       int           `json:"numBans"`
		TheirTeamBans []interface{} `json:"theirTeamBans"`
	} `json:"bans"`
	BenchChampionIds   []interface{} `json:"benchChampionIds"`
	BenchEnabled       bool          `json:"benchEnabled"`
	BoostableSkinCount int           `json:"boostableSkinCount"`
	ChatDetails        struct {
		ChatRoomName     string      `json:"chatRoomName"`
		ChatRoomPassword interface{} `json:"chatRoomPassword"`
	} `json:"chatDetails"`
	Counter              int `json:"counter"`
	EntitledFeatureState struct {
		AdditionalRerolls int           `json:"additionalRerolls"`
		UnlockedSkinIds   []interface{} `json:"unlockedSkinIds"`
	} `json:"entitledFeatureState"`
	GameID               int64 `json:"gameId"`
	HasSimultaneousBans  bool  `json:"hasSimultaneousBans"`
	HasSimultaneousPicks bool  `json:"hasSimultaneousPicks"`
	IsCustomGame         bool  `json:"isCustomGame"`
	IsSpectating         bool  `json:"isSpectating"`
	LocalPlayerCellID    int   `json:"localPlayerCellId"`
	LockedEventIndex     int   `json:"lockedEventIndex"`
	MyTeam               []struct {
		AssignedPosition    string `json:"assignedPosition"`
		CellID              int    `json:"cellId"`
		ChampionID          int    `json:"championId"`
		ChampionPickIntent  int    `json:"championPickIntent"`
		EntitledFeatureType string `json:"entitledFeatureType"`
		SelectedSkinID      int    `json:"selectedSkinId"`
		Spell1ID            int    `json:"spell1Id"`
		Spell2ID            int    `json:"spell2Id"`
		SummonerID          int    `json:"summonerId"`
		Team                int    `json:"team"`
		WardSkinID          int    `json:"wardSkinId"`
	} `json:"myTeam"`
	RerollsRemaining   int  `json:"rerollsRemaining"`
	SkipChampionSelect bool `json:"skipChampionSelect"`
	TheirTeam          []struct {
		AssignedPosition    string `json:"assignedPosition"`
		CellID              int    `json:"cellId"`
		ChampionID          int    `json:"championId"`
		ChampionPickIntent  int    `json:"championPickIntent"`
		EntitledFeatureType string `json:"entitledFeatureType"`
		SelectedSkinID      int    `json:"selectedSkinId"`
		Spell1ID            int    `json:"spell1Id"`
		Spell2ID            int    `json:"spell2Id"`
		SummonerID          int    `json:"summonerId"`
		Team                int    `json:"team"`
		WardSkinID          int    `json:"wardSkinId"`
	} `json:"theirTeam"`
	Timer struct {
		AdjustedTimeLeftInPhase int    `json:"adjustedTimeLeftInPhase"`
		InternalNowInEpochMs    int64  `json:"internalNowInEpochMs"`
		IsInfinite              bool   `json:"isInfinite"`
		Phase                   string `json:"phase"`
		TotalTimeInPhase        int    `json:"totalTimeInPhase"`
	} `json:"timer"`
	Trades []interface{} `json:"trades"`
}

func readLock() string {
	var file *os.File
	var err error

	for {
		file, err = os.Open(config.ClientDir + "lockfile")
		if err == nil {
			//fmt.Println("Lockfile found")
			break
		} else {
			fmt.Println("Waiting for League process to open")
		}
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}

	b, err := ioutil.ReadAll(file)

	if err != nil {
		panic(err)
	}

	return string(b)
}

func requestApi(command *string) io.ReadCloser {
	req, err := http.NewRequest("GET", values[4]+"://127.0.0.1:"+values[2]+(*command), nil)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("riot", values[3])

	resp, err := cli.Do(req)

	if err != nil {
		panic(err)
	}

	//for err != nil || resp.StatusCode != 200 {
	//	resp, err = cli.Do(req)
	//	time.Sleep(time.Duration(interval) * time.Second)
	//}

	//err = req.Body.Close()
	//
	//if err != nil {
	//	panic(err)
	//}

	return resp.Body
}

func isInChampSelect() bool {
	command := "/lol-champ-select/v1/session"
	var data ChampSelect

	err := json.NewDecoder(requestApi(&command)).Decode(&data)

	if err != nil {
		panic(err)
	}

	if float64(data.Timer.AdjustedTimeLeftInPhase) <= config.Interval {
		return false
	}
	return true

	//command := "/lol-champ-select/v1/session"
	//req, err := http.NewRequest("GET", values[4]+"://127.0.0.1:"+values[2]+command, nil)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//req.SetBasicAuth("riot", values[3])
	//
	//resp, err := cli.Do(req)
	//
	//if err != nil {
	//	panic(err)
	//}

	//return resp.StatusCode == 200
}

func getAccInfo() (int64, int) {
	command := "/lol-summoner/v1/current-summoner" // returns login information

	var data AccountInfo
	err := json.NewDecoder(requestApi(&command)).Decode(&data)

	if err != nil {
		panic(err)
	}

	fmt.Println("Logged in as...")
	fmt.Println("Account ID: ", data.AccountID)
	fmt.Println("Display Name: ", data.DisplayName)
	fmt.Println("Internal Name: ", data.InternalName)
	fmt.Println("Player UUID: ", data.Puuid)
	fmt.Println("Summoner ID: ", data.SummonerID)

	return data.AccountID, data.SummonerID
}

func getQueueId() (int, bool) {
	command := "/lol-gameflow/v1/gameflow-metadata/player-status"
	var queueInfo QueueInfo
	err := json.NewDecoder(requestApi(&command)).Decode(&queueInfo)
	if err != nil {
		panic(err)
	}

	return queueInfo.CurrentLobbyStatus.QueueID, queueInfo.CurrentLobbyStatus.IsCustom
}

func deleteRunePage(runePageId int) bool {

	req, err := http.NewRequest("DELETE", values[4]+"://127.0.0.1:"+values[2]+"/lol-perks/v1/pages/"+strconv.Itoa(runePageId), nil)

	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("riot", values[3])

	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 204 {
		fmt.Println("Old AutoRune page deleted")
		return true
	}
	return false
}

func setItems(doc *soup.Root, accId int64, sumId int, champId int, gameType *string) {
	builds := (*doc).FindAll("tr", "class", "champion-overview__row")
	blockCnt := len((*doc).FindAll("tr", "class", "champion-overview__row--first"))
	blockCnt++

	blockList := make([]ItemBlock, blockCnt)
	otherItemSet := make(map[string]bool)
	alreadyAdded := 0

	i := 0
	for _, build := range builds {
		if strings.HasSuffix(build.Attrs()["class"], "champion-overview__row--first") {
			items := build.FindAll("li", "class", "champion-stats__list__item")
			itemList := make([]Item, len(items))
			for j, img := range items {
				str := img.Find("img").Attrs()["src"]
				str = str[strings.LastIndex(str, "/")+1 : strings.Index(str, ".png")]
				newItem := Item{
					Count: 1,
					ID:    str,
				}
				otherItemSet[str] = false
				alreadyAdded++
				itemList[j] = newItem
			}
			newItemBlock := ItemBlock{
				HideIfSummonerSpell: "",
				Items:               itemList,
				ShowIfSummonerSpell: "",
				Type:                build.Find("th", "class", "champion-overview__sub-header").Text(),
			}
			blockList[i] = newItemBlock
			i++
		} else {
			items := build.FindAll("li", "class", "champion-stats__list__item")
			for _, img := range items {
				str := img.Find("img").Attrs()["src"]
				str = str[strings.LastIndex(str, "/")+1 : strings.Index(str, ".png")]
				if _, hasElem := otherItemSet[str]; !hasElem {
					otherItemSet[str] = true
				}
			}
		}
	}

	ward := Item{
		Count: 1,
		ID:    "3340",
	}

	blockList[0].Items = append(blockList[0].Items, ward)

	consumable := ItemBlock{
		HideIfSummonerSpell: "",
		Items: []Item{
			{
				Count: 1,
				ID:    "2055",
			},
			{
				Count: 1,
				ID:    "3340",
			},
			{
				Count: 1,
				ID:    "3363",
			},
			{
				Count: 1,
				ID:    "3364",
			},
			{
				Count: 1,
				ID:    "2047",
			},
			{
				Count: 1,
				ID:    "2138",
			},
			{
				Count: 1,
				ID:    "2139",
			},
			{
				Count: 1,
				ID:    "2140",
			},
		},
		ShowIfSummonerSpell: "",
		Type:                "Consumables",
	}

	i++
	blockList = append(blockList[:2], blockList[1:]...)
	blockList[1] = consumable

	//fmt.Println(len(otherItemSet))
	//fmt.Println(alreadyAdded)

	itemList := make([]Item, len(otherItemSet)-alreadyAdded)
	idx := 0

	for otherItem := range otherItemSet {
		if otherItemSet[otherItem] {
			newItem := Item{
				Count: 1,
				ID:    otherItem,
			}
			itemList[idx] = newItem
			idx++
		}
	}

	newItemBlock := ItemBlock{
		HideIfSummonerSpell: "",
		Items:               itemList,
		ShowIfSummonerSpell: "",
		Type:                "Other items to consider",
	}

	blockList[i] = newItemBlock
	i++

	itemPage := ItemPage{
		AccountID: accId,
		ItemSets: []ItemSet{
			{
				AssociatedChampions: []int{champId},
				AssociatedMaps:      []int{11, 12},
				Blocks:              blockList,
				Map:                 "any",
				Mode:                "any",
				PreferredItemSlots:  make([]interface{}, 0),
				Sortrank:            0,
				StartedFrom:         "blank",
				Title:               ProjectName + " Item Page " + (*gameType),
				Type:                "custom",
				UID:                 "",
			},
		},
		Timestamp: 0,
	}

	command := "/lol-item-sets/v1/item-sets/" + strconv.Itoa(sumId) + "/sets"

	//j, _ := json.Marshal(newRune)
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(itemPage)

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", values[4]+"://127.0.0.1:"+values[2]+command, b)

	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("riot", values[3])

	resp, errs := cli.Do(req)
	if errs != nil {
		panic(errs)
	}

	if resp.StatusCode == 201 {
		fmt.Println("Item page updated successfully")
	} else {
		fmt.Println("Item page update failed")
	}
}

func setRunes(doc *soup.Root, gameType *string) {
	// Delete "AutoRune" page
	command := "/lol-perks/v1/pages"
	var runePages RunePages
	err := json.NewDecoder(requestApi(&command)).Decode(&runePages)

	if err != nil {
		panic(err)
	}

	//fmt.Println("Total Rune pages:", len(runePages))
	for _, page := range runePages {
		if strings.HasPrefix(page.Name, "AutoRune") {
			deleteRunePage(page.ID)
		}
	}

	// Could be converted to FindAll
	links := (*doc).Find("div", "class", "perk-page-wrap")

	//links := (*doc).FindAll("div", "class", "perk-page-wrap")
	//for _, link := range links{
	//	do work here
	//}

	// Category
	imgs := links.FindAll("div", "class", "perk-page__item--mark")

	runeCategoryList := make([]int, len(imgs))

	if len(imgs) != 2 {
		panic("Rune category updated?")
	}

	for i, img := range imgs {
		str := img.Find("img").Attrs()["src"]
		str = str[strings.LastIndex(str, "/")+1 : strings.Index(str, ".png")]
		runeCategoryList[i], _ = strconv.Atoi(str)
	}

	// Runes
	imgs = links.FindAll("div", "class", "perk-page__item--active")
	// Fragments
	fragImgs := links.FindAll("div", "class", "fragment__row")

	runeList := make([]int, len(imgs)+len(fragImgs))

	if len(runeList) != 9 {
		panic("Runes updated?")
	}

	for i, img := range imgs {
		str := img.Find("img").Attrs()["src"]
		str = str[strings.LastIndex(str, "/")+1 : strings.Index(str, ".png")]
		runeList[i], _ = strconv.Atoi(str)
	}

	for i, img := range fragImgs {
		str := img.Find("img", "class", "active").Attrs()["src"]
		str = str[strings.LastIndex(str, "/")+1 : strings.Index(str, ".png")]
		runeList[len(imgs)+i], _ = strconv.Atoi(str)
	}

	newRune := RunePage{
		AutoModifiedSelections: make([]interface{}, 0),
		Current:                true,
		ID:                     0,
		IsActive:               true,
		IsDeletable:            true,
		IsEditable:             true,
		IsValid:                true,
		LastModified:           0,
		Name:                   "AutoRune " + (*gameType),
		Order:                  0,
		PrimaryStyleID:         runeCategoryList[0],
		SelectedPerkIds:        runeList,
		SubStyleID:             runeCategoryList[1],
	}

	command = "/lol-perks/v1/pages"

	//j, _ := json.Marshal(newRune)
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(newRune)

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", values[4]+"://127.0.0.1:"+values[2]+command, b)

	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("riot", values[3])

	resp, errs := cli.Do(req)
	if errs != nil {
		panic(errs)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Rune Update successful")
	} else {
		fmt.Println("Rune Update failed")
	}
}

func getRunes(accId int64, sumId int, champId int, queueId int) {
	command := "/lol-champions/v1/inventories/" + strconv.Itoa(sumId) + "/champions/" + strconv.Itoa(champId)
	var data Champion
	var gameType, url string

	err := json.NewDecoder(requestApi(&command)).Decode(&data)

	if err != nil {
		panic(err)
	}
	// TODO: use flexible region here

	if queueId == 900 {
		gameType = "URF"
		fmt.Println("ULTRA RAPID FIRE MODE IS ON!!!")
		url = "https://na.op.gg/urf/" + data.Alias + "/statistics"
	} else {
		url = "https://na.op.gg/champion/" + data.Alias
	}

	fmt.Println("Selected Champion: ", data.Alias)

	resp, err := soup.Get(url)
	if err != nil {
		panic(err)
	}

	doc := soup.HTMLParse(resp)

	setRunes(&doc, &gameType)
	setItems(&doc, accId, sumId, champId, &gameType)

	// Find champion positions
	positions := doc.FindAll("li", "class", "champion-stats-header__position")

	if len(positions) == 1 {
		fmt.Println("No alternative positions available.")
	} else if len(positions) > 1 {
		posUrlList := make([]string, len(positions))

		for i, pos := range positions {
			link := "https://na.op.gg" + pos.Find("a").Attrs()["href"]
			role := pos.Find("span", "class", "champion-stats-header__position__role").Text()
			rate := pos.Find("span", "class", "champion-stats-header__position__rate").Text()
			fmt.Println(i, ". "+role+": ", rate)
			posUrlList[i] = link
		}

		fmt.Println("Current role: 0")

		var i int
		for i != -1 {
			fmt.Print("Change role to... (-1 to exit): ")
			_, err = fmt.Scan(&i)

			if err != nil {
				panic(err)
			}

			if i != -1 {
				resp, err := soup.Get(posUrlList[i])
				if err != nil {
					panic(err)
				}
				doc := soup.HTMLParse(resp)
				setRunes(&doc, &gameType)
				setItems(&doc, accId, sumId, champId, &gameType)
				fmt.Println("Current role:", i)
			}
		}
	}
}

func getChampId(sumId int) int {
	command := "/lol-champ-select/v1/session"
	var data ChampSelect
	var champId = 0

	err := json.NewDecoder(requestApi(&command)).Decode(&data)

	if err != nil {
		panic(err)
	}

	// Find current user's champion ID
	for _, member := range data.MyTeam {
		if member.SummonerID == sumId {
			champId = member.ChampionID
		}
	}

	return champId
}

//func pageTest(){
//	resp, err := http.Get("https://na.op.gg/champion/Lucian")
//	if err != nil {
//		panic(err)
//	}
//	responseData,err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(string(responseData))
//}

func ReadConfig() {
	file, err := os.Open("./config")

	config = Config{
		Debug:      false,
		Interval:   2,
		ClientDir:  "C:/Riot Games/League of Legends/",
		EnableRune: true,
		EnableItem: true,
	}

	if err != nil {
		fmt.Println("Cannot open config file. Default settings will be used.")
	} else {
		defer file.Close()
		configMap := make(map[string]string)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// Line with # is comments
			if !strings.HasPrefix(line, "#") {
				args := strings.SplitN(line, "=", 2)
				configMap[strings.TrimSpace(args[0])] = strings.TrimSpace(args[1])
			}
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}

		config.ClientDir = configMap["CLIENT_DIRECTORY"]
		config.Debug, err = strconv.ParseBool(configMap["DEBUG"])
		if err != nil {
			panic(err)
		}
		config.EnableRune, err = strconv.ParseBool(configMap["ENABLE_RUNE"])
		if err != nil {
			panic(err)
		}
		config.EnableItem, err = strconv.ParseBool(configMap["ENABLE_ITEM"])
		if err != nil {
			panic(err)
		}
		config.Interval, err = strconv.ParseFloat(configMap["INTERVAL"], 64)
		if err != nil {
			panic(err)
		}

		if config.Interval < 1 {
			config.Interval = 1
		} else if config.Interval > 5 {
			config.Interval = 5
		}
	}
}

func run(status *widget.Label) {
	status.SetText("Running")
	// Read lockfile
	content := readLock()
	values = strings.Split(content, ":")

	if config.Debug {
		for _, val := range values {
			fmt.Println(val)
		}
	}

	//command := "/lol-summoner/v1/current-summoner" // returns login information
	//command := "/lol-champ-select/v1/session" // champion select session information
	//command := "/riotclient/auth-token" // returns auth token
	//command := "/riotclient/region-locale" // returns region, language etc
	//command := "/lol-perks/v1/currentpage" // get current rune page

	accId, sumId := getAccInfo()

	//var isCustomGame = false
	var prevChampId, champId int
	var queueId = -1

	//// Check game type (norms, urf)
	//for queueId == -1 && !isCustomGame {
	//	queueId, isCustomGame = getQueueId()
	//	time.Sleep(time.Duration(interval) * time.Second)
	//}

	// Check if in lobby
	for champId == 0 {
		fmt.Println("Waiting for a champion to be selected...")
		queueId, _ = getQueueId()
		champId = getChampId(sumId)
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}

	// Loop until champion select phase is over
	for isInChampSelect() {
		fmt.Println("Checking if Champion ID was updated...")
		champId := getChampId(sumId)

		if champId != 0 && prevChampId != champId {
			getRunes(accId, sumId, champId, queueId)
			prevChampId = champId
		}
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
	status.SetText("Not running")
}

func main() {
	// Initialize http client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli = &http.Client{Transport: tr}

	ReadConfig()

	a := app.New()
	w := a.NewWindow(ProjectName + " " + Version)

	sl := widget.NewSlider(1, 3)
	sl.Value = config.Interval
	sl.Step = 0.5
	sl.OnChanged = func(f float64) {
		fmt.Println("Poll interval: ", f)
		config.Interval = f
	}

	status := widget.NewLabel("Not running")

	enableRunesCheck := widget.NewCheck("", func(b bool) {
		fmt.Println("Auto runes: ", b)
		config.EnableRune = b
	})
	enableRunesCheck.SetChecked(config.EnableRune)

	enableItemsCheck := widget.NewCheck("", func(b bool) {
		fmt.Println("Auto items: ", b)
		config.EnableItem = b
	})
	enableItemsCheck.SetChecked(config.EnableItem)

	w.SetContent(
		widget.NewVBox(
			widget.NewLabel(ProjectName+" "+Version),
			widget.NewHBox(
				widget.NewVBox(
					widget.NewButton("Manual Start", func() {
						go run(status)
					}),
					status),
				widget.NewVBox(
					widget.NewHBox(widget.NewLabel("Auto runes"), enableRunesCheck),
					widget.NewHBox(widget.NewLabel("Auto items"), enableItemsCheck),
					widget.NewLabel("Polling interval"),
					sl))))

	w.SetFixedSize(true)
	w.ShowAndRun()
}
