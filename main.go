package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const InstallDir string = "D:/Games/Riot Games/League of Legends/"

//const InstallDir string = "C:/Riot Games/League of Legends/"

var cli *http.Client
var values []string
var interval = 3
var debug = true

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
		file, err = os.Open(InstallDir + "lockfile")
		if err == nil {
			//fmt.Println("Lockfile found")
			break
		} else {
			fmt.Println("Waiting for League process to open")
		}
		time.Sleep(time.Duration(interval) * time.Second)
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

	if data.Timer.AdjustedTimeLeftInPhase <= interval {
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

func getAccInfo() int {
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

	return data.SummonerID
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

func getRunes(sumId int, champId int, queueId int) {
	command := "/lol-champions/v1/inventories/" + strconv.Itoa(sumId) + "/champions/" + strconv.Itoa(champId)
	var data Champion
	var gameType, url string

	err := json.NewDecoder(requestApi(&command)).Decode(&data)

	if err != nil {
		panic(err)
	}
	// TODO: use flexible region here
	// TODO: Add URF link here

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

func main() {
	// Initialize http client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli = &http.Client{Transport: tr}

	fmt.Println(InstallDir)

	// Read lockfile
	content := readLock()
	values = strings.Split(content, ":")

	if debug {
		for _, val := range values {
			fmt.Println(val)
		}
	}

	//command := "/lol-summoner/v1/current-summoner" // returns login information
	//command := "/lol-champ-select/v1/session" // champion select session information
	//command := "/riotclient/auth-token" // returns auth token
	//command := "/riotclient/region-locale" // returns region, language etc
	//command := "/lol-perks/v1/currentpage" // get current rune page

	sumId := getAccInfo()

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
		time.Sleep(time.Duration(interval) * time.Second)
	}

	// Loop until champion select phase is over
	for isInChampSelect() {
		fmt.Println("Checking if Champion ID was updated...")
		champId := getChampId(sumId)

		if champId != 0 && prevChampId != champId {
			getRunes(sumId, champId, queueId)
			prevChampId = champId
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}

}
