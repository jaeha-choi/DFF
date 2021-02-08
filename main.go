package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/anaskhan96/soup"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const ProjectName string = "DFF!"
const Version string = "v0.5.2"

var cli *http.Client
var config Config

var values []string
var lastRole = ""

//var interval = 3
//var debug = false

//type Spell struct {
//	ID           string `json:"id"`
//	Name         string `json:"name"`
//	Description  string `json:"description"`
//	Tooltip      string `json:"tooltip"`
//	Maxrank      int    `json:"maxrank"`
//	Cooldown     []int  `json:"cooldown"`
//	CooldownBurn string `json:"cooldownBurn"`
//	Cost         []int  `json:"cost"`
//	CostBurn     string `json:"costBurn"`
//	Datavalues   struct {
//	} `json:"datavalues"`
//	Effect        []interface{} `json:"effect"`
//	EffectBurn    []interface{} `json:"effectBurn"`
//	Vars          []interface{} `json:"vars"`
//	Key           string        `json:"key"`
//	SummonerLevel int           `json:"summonerLevel"`
//	Modes         []string      `json:"modes"`
//	CostType      string        `json:"costType"`
//	Maxammo       string        `json:"maxammo"`
//	Range         []int         `json:"range"`
//	RangeBurn     string        `json:"rangeBurn"`
//	Image         struct {
//		Full   string `json:"full"`
//		Sprite string `json:"sprite"`
//		Group  string `json:"group"`
//		X      int    `json:"x"`
//		Y      int    `json:"y"`
//		W      int    `json:"w"`
//		H      int    `json:"h"`
//	} `json:"image"`
//	Resource string `json:"resource"`
//}
type GitHubUpdate struct {
	URL             string `json:"url"`
	AssetsURL       string `json:"assets_url"`
	UploadURL       string `json:"upload_url"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	NodeID          string `json:"node_id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		URL      string      `json:"url"`
		ID       int         `json:"id"`
		NodeID   string      `json:"node_id"`
		Name     string      `json:"name"`
		Label    interface{} `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}

type SpellFile struct {
	Type    string                 `json:"type"`
	Version string                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}

type FileVersion struct {
	Version string `json:"version"`
}

type MySelect struct {
	//SelectedSkinID int32 `json:"selectedSkinId"`
	Spell1ID int64 `json:"spell1Id"`
	Spell2ID int64 `json:"spell2Id"`
	//WardSkinID     int64 `json:"wardSkinId"`
}

type Config struct {
	Debug       bool    `json:"debug"`
	Interval    float64 `json:"interval"`
	ClientDir   string  `json:"client_dir"`
	EnableRune  bool    `json:"enable_rune"`
	EnableItem  bool    `json:"enable_item"`
	EnableSpell bool    `json:"enable_spell"`
	DFlash      bool    `json:"d_flash"`
	Language    string  `json:"language"`
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

type RunePageCount struct {
	OwnedPageCount int `json:"ownedPageCount"`
}

type RuneNamePage struct {
	Name string
	Page RunePage
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

func requestApi(command *string, loopFlag bool) io.ReadCloser {
	req, err := http.NewRequest("GET", values[4]+"://127.0.0.1:"+values[2]+(*command), nil)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("riot", values[3])

	resp, err := cli.Do(req)

	if loopFlag {
		for err != nil || resp.StatusCode != 200 {
			resp, err = cli.Do(req)
			time.Sleep(time.Duration(config.Interval) * time.Second)
		}
	}

	if err != nil {
		log.Print("API no longer available:", err)
		os.Exit(1)
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

	err := json.NewDecoder(requestApi(&command, false)).Decode(&data)

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
	err := json.NewDecoder(requestApi(&command, true)).Decode(&data)

	if err != nil {
		panic(err)
	}

	if config.Debug {
		fmt.Println("Logged in as...")
		fmt.Println("Account ID: ", data.AccountID)
		fmt.Println("Display Name: ", data.DisplayName)
		fmt.Println("Internal Name: ", data.InternalName)
		fmt.Println("Player UUID: ", data.Puuid)
		fmt.Println("Summoner ID: ", data.SummonerID)
	}
	fmt.Println("Client API connection functional.")
	return data.AccountID, data.SummonerID
}

func getUxStatus() bool {
	command := "/riotclient/ux-state"

	bodyBytes, err := ioutil.ReadAll(requestApi(&command, false))
	if err != nil {
		panic(err)
	}
	if string(bodyBytes) == "\"ShowMain\"" {
		return false
	} else {
		return true
	}
}

func getQueueId() (int, bool) {
	command := "/lol-gameflow/v1/gameflow-metadata/player-status"
	var queueInfo QueueInfo
	err := json.NewDecoder(requestApi(&command, false)).Decode(&queueInfo)
	if err != nil {
		log.Print("API no longer available:", err)
		os.Exit(1)
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
		fmt.Println("Old DFF page deleted")
		return true
	}
	return false
}

func setItems(doc *soup.Root, accId int64, sumId int, champId int, gameType *string) {
	if !config.EnableItem {
		return
	}
	builds := (*doc).FindAll("tr", "class", "champion-overview__row")
	blockCnt := len((*doc).FindAll("tr", "class", "champion-overview__row--first"))
	blockCnt++

	blockList := make([]ItemBlock, blockCnt)
	otherItemSet := make(map[string]bool)
	willBeAdded := 0
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
					willBeAdded++
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

	if config.Debug {
		fmt.Println("Total number of items:", len(otherItemSet))
		fmt.Println("Count of items that will be added:", willBeAdded)
	}

	itemList := make([]Item, willBeAdded)
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

func setSpells(doc *soup.Root) {
	if !config.EnableSpell {
		return
	}

	imgs := (*doc).Find("td", "class", "champion-overview__data")

	spellImgs := imgs.FindAll("img", "class", "tip")

	spellNameList := make([]string, len(spellImgs))
	for i, img := range spellImgs {
		str := img.Attrs()["src"]
		spellNameList[i] = str[strings.LastIndex(str, "/")+1 : strings.Index(str, ".png")]
	}

	var spellFile SpellFile

	f, err := ioutil.ReadFile("./data/" + "summoner.json")
	if err != nil {
		panic(err)
	}
	//var data interface{}
	err = json.Unmarshal(f, &spellFile)
	if err != nil {
		panic(err)
	}

	if config.Debug {
		fmt.Println("summoner.json version: ", spellFile.Version)
		fmt.Println("Total spell count: ", len(spellFile.Data))
	}

	spellKeyList := make([]int64, len(spellImgs))
	for x := 0; x < 10; x++ {

		i := 0
		for _, spells := range spellFile.Data {
			/* This works too, but takes 10 times longer
			var spell Spell
			jsonStr, err := json.Marshal(spells)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(jsonStr, &spell)
			//fmt.Println(spell.Name)
			//fmt.Println(spell.Key)
			*/

			news := spells.(map[string]interface{})
			tempSpellName := news["id"]
			if tempSpellName == spellNameList[0] || tempSpellName == spellNameList[1] {
				spellKeyList[i], _ = strconv.ParseInt(news["key"].(string), 10, 64)
				i++
			}
		}
	}

	// If user is using D key as flash, set flash for D
	// If user is using F key as flash, set flash for F
	// Otherwise, no change.
	if config.DFlash && spellKeyList[1] == 4 {
		spellKeyList[1] = spellKeyList[0]
		spellKeyList[0] = 4
	} else if !config.DFlash && spellKeyList[0] == 4 {
		spellKeyList[0] = spellKeyList[1]
		spellKeyList[1] = 4
	}

	spells := MySelect{
		Spell1ID: spellKeyList[0],
		Spell2ID: spellKeyList[1],
	}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(spells)

	command := "/lol-champ-select/v1/session/my-selection"
	req, err := http.NewRequest("PATCH", values[4]+"://127.0.0.1:"+values[2]+command, b)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("riot", values[3])

	resp, err := cli.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 204 {
		fmt.Println("Spell update successful")
	} else {
		fmt.Println("Spell update/skin update fail")
	}
}

func setRuneHelper(page RunePage) {
	delRunes()

	command := "/lol-perks/v1/pages"

	//j, _ := json.Marshal(newRune)
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(page)

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

func delRunes() {
	var runePages RunePages
	var runePageCnt RunePageCount
	deleted := false

	// Delete "DFF" page
	command := "/lol-perks/v1/pages"
	err := json.NewDecoder(requestApi(&command, false)).Decode(&runePages)
	if err != nil {
		panic(err)
	}

	command = "/lol-perks/v1/inventory"
	err = json.NewDecoder(requestApi(&command, false)).Decode(&runePageCnt)
	if err != nil {
		panic(err)
	}

	//fmt.Println("Total Rune pages:", len(runePages))
	for _, page := range runePages {
		//fmt.Println(page.Name)
		if strings.HasPrefix(page.Name, ProjectName) {
			deleteRunePage(page.ID)
			deleted = true
		}
	}
	if len(runePages)+5 >= runePageCnt.OwnedPageCount && !deleted {
		// Delete the first rune page if all pages are used
		deleteRunePage(runePages[0].ID)
	}
}

func setRunes(doc *soup.Root, gameType *string) ([]RuneNamePage, [][]string) {
	if !config.EnableRune {
		return nil, nil
	}

	//delRunes()

	runeDetailsDoc := (*doc).FindAll("span", "class", "pick-ratio__text")
	runeDetails := make([][]string, len(runeDetailsDoc)*2)

	var pr, wr, sample string
	for idx, runeDetailDoc := range runeDetailsDoc {
		next := runeDetailDoc.FindNextElementSibling()
		pr = next.Text()
		next = next.FindNextElementSibling()
		sample = next.Text()
		next = next.FindNextElementSibling()
		next = next.FindNextElementSibling()
		wr = next.Text()
		runeDetail := []string{pr, wr, sample}
		runeDetails[idx] = runeDetail
	}

	runeNames := (*doc).FindAll("div", "class", "champion-stats-summary-rune__name")
	runeInfo := make([]RuneNamePage, len(runeNames)*2)

	i := 0
	for _, runeName := range runeNames {
		names := strings.Split(runeName.Text(), "+")
		//fmt.Println(runeName.Text())
		for x := 0; x < 2; x++ {
			runeInfo[i] = RuneNamePage{
				Name: string([]rune(strings.TrimSpace(names[0]))[0]) + "+" + string([]rune(strings.TrimSpace(names[1]))[0]) + " (" + strconv.Itoa(x+1) + ")",
			}
			i++
			//fmt.Println(runeInfo[x].Name)
		}
	}
	// Could be converted to FindAll
	//links := (*doc).Find("div", "class", "perk-page-wrap")

	links := (*doc).FindAll("div", "class", "perk-page-wrap")
	for x, link := range links {
		//do work here
		// Category
		imgs := link.FindAll("div", "class", "perk-page__item--mark")

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
		imgs = link.FindAll("div", "class", "perk-page__item--active")
		// Fragments
		fragImgs := link.FindAll("div", "class", "fragment__row")

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

		runeInfo[x].Page = RunePage{
			AutoModifiedSelections: make([]interface{}, 0),
			Current:                true,
			ID:                     0,
			IsActive:               true,
			IsDeletable:            true,
			IsEditable:             true,
			IsValid:                true,
			LastModified:           0,
			Name:                   ProjectName + " " + runeInfo[x].Name + " " + (*gameType),
			Order:                  0,
			PrimaryStyleID:         runeCategoryList[0],
			SelectedPerkIds:        runeList,
			SubStyleID:             runeCategoryList[1],
		}
	}

	setRuneHelper(runeInfo[0].Page)

	return runeInfo, runeDetails
}

func downloadFile(url string) {
	fileName := strings.Split(url, "/")
	out, err := os.Create("./data/" + fileName[len(fileName)-1])
	if err != nil {
		panic(err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(fileName[len(fileName)-1] + " downloaded.")
}

func getRunes(accId int64, sumId int, champId int, queueId int, champLabel *widget.Label) ([][4]string, []RuneNamePage, [][]string) {
	command := "/lol-champions/v1/inventories/" + strconv.Itoa(sumId) + "/champions/" + strconv.Itoa(champId)
	var data Champion
	var gameType, url string
	var posUrlList [][4]string = nil
	var runeNamePages []RuneNamePage = nil
	var runeDetails [][]string

	err := json.NewDecoder(requestApi(&command, false)).Decode(&data)

	if err != nil {
		panic(err)
	}
	champLabel.SetText(data.Alias)
	fmt.Println("Selected Champion: ", data.Alias)

	if queueId == 900 {
		gameType = "URF"
		fmt.Println("ULTRA RAPID FIRE MODE IS ON!!!")
		url = "https://op.gg/urf/" + data.Alias + "/statistics"
		soup.Cookie("customLocale", config.Language)
		resp, err := soup.Get(url)
		if err != nil {
			panic(err)
		}

		doc := soup.HTMLParse(resp)

		runeNamePages, runeDetails = setRunes(&doc, &gameType)
		setItems(&doc, accId, sumId, champId, &gameType)
		setSpells(&doc)

	} else {
		// Can add region here
		url = "https://op.gg/champion/" + data.Alias
		soup.Cookie("customLocale", config.Language)
		resp, err := soup.Get(url)
		if err != nil {
			panic(err)
		}

		doc := soup.HTMLParse(resp)

		runeNamePages, runeDetails = setRunes(&doc, &gameType)
		setItems(&doc, accId, sumId, champId, &gameType)
		setSpells(&doc)

		// Find champion positions
		positions := doc.FindAll("li", "class", "champion-stats-header__position")

		//if len(positions) == 1 {
		//	fmt.Println("No alternative positions available.")
		//} else if len(positions) > 1 {
		posUrlList = make([][4]string, len(positions))

		for i, pos := range positions {
			link := "https://op.gg" + pos.Find("a").Attrs()["href"]
			role := pos.Find("span", "class", "champion-stats-header__position__role").Text()
			rate := pos.Find("span", "class", "champion-stats-header__position__rate").Text()

			fmt.Println(i, ". "+role+": ", rate)
			posUrlList[i][0] = strings.TrimSpace(role)
			posUrlList[i][1] = rate
			posUrlList[i][2] = link
			posUrlList[i][3] = gameType

		}
		lastRole = posUrlList[0][0]

	}
	return posUrlList, runeNamePages, runeDetails
}

func getChampId(sumId int) int {
	command := "/lol-champ-select/v1/session"
	var data ChampSelect
	var champId = 0

	err := json.NewDecoder(requestApi(&command, false)).Decode(&data)

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

func writeConfig() {
	jsonConf, _ := json.MarshalIndent(config, "", "\t")
	err := ioutil.WriteFile("config.json", jsonConf, 0644)
	if err != nil {
		panic(err)
	}
}

func readConfig() {
	_, err := os.Stat("config.json")

	config = Config{
		Debug:       false,
		Interval:    2,
		ClientDir:   "C:/Riot Games/League of Legends/",
		EnableRune:  true,
		EnableItem:  true,
		EnableSpell: true,
		DFlash:      true,
		Language:    "en_US",
	}

	if os.IsNotExist(err) {
		fmt.Println("Cannot open config file. Default settings will be used.")
		writeConfig()
	} else {
		jsonConf, _ := ioutil.ReadFile("config.json")
		err := json.Unmarshal(jsonConf, &config)
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

func checkFiles() {
	resp, err := http.Get("https://ddragon.leagueoflegends.com/api/versions.json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var version []string
	err = json.NewDecoder(resp.Body).Decode(&version)

	if err != nil {
		panic(err)
	}

	leagueVersion := version[0]

	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		_ = os.Mkdir("./data", 0700)
		// Summoner spells
		downloadFile("http://ddragon.leagueoflegends.com/cdn/" + leagueVersion + "/data/en_US/" + "summoner.json")
		// Items
		//downloadFile("http://ddragon.leagueoflegends.com/cdn/"+ leagueVersion +"/data/en_US/"+"item.json")
		// Maps
		//downloadFile("http://ddragon.leagueoflegends.com/cdn/"+ leagueVersion +"/data/en_US/"+"map.json.json")
		// Runes
		//downloadFile("http://ddragon.leagueoflegends.com/cdn/"+ leagueVersion +"/data/en_US/"+"runesReforged.json")
		// Champions
		//downloadFile("http://ddragon.leagueoflegends.com/cdn/"+ leagueVersion +"/data/en_US/"+"champion.json")
	} else {
		files, err := ioutil.ReadDir("./data")
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			var ver FileVersion
			f, err := ioutil.ReadFile("./data/" + file.Name())
			if err != nil {
				panic(err)
			}
			//var data interface{}
			err = json.Unmarshal(f, &ver)
			if err != nil {
				panic(err)
			}
			if leagueVersion != ver.Version {
				downloadFile("http://ddragon.leagueoflegends.com/cdn/" + leagueVersion + "/data/en_US/" + file.Name())
			}
		}
	}
}

func run(status *widget.Label, p *widget.Select, champLabel *widget.Label, wait *sync.WaitGroup, runeSelect *widget.Select) {
	defer wait.Done()
	status.SetText("Starting...")
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
		status.SetText("Waiting...")
		fmt.Println("Waiting for a champion to be selected...")
		queueId, _ = getQueueId()
		champId = getChampId(sumId)
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}

	// For debug purpose; replace command to test the api
	//command := "/lol-perks/v1/show-auto-modified-pages-notification"
	//resp := requestApi(&command, false)
	//body, _ := ioutil.ReadAll(resp)
	//fmt.Println(string(body))

	// Loop until champion select phase is over
	for isInChampSelect() {
		champId := getChampId(sumId)

		if champId != 0 && prevChampId != champId {
			status.SetText("Setting...")
			result, runeNamePages, runeDetails := getRunes(accId, sumId, champId, queueId, champLabel)
			status.SetText("Updated...")

			if len(runeNamePages) > 0 {
				options := make([]string, len(runeNamePages))
				for x, elem := range runeNamePages {
					runeDetail := runeDetails[x]
					options[x] = elem.Name + " PR:" + runeDetail[0] + " WR:" + runeDetail[1] + " Sample:" + runeDetail[2]
					//options[x] = elem.Name + "  PR:" + runeDetail[0] + " WR:" + runeDetail[1]
				}
				runeSelect.Options = options
				runeSelect.Selected = options[0]
				runeSelect.OnChanged = func(s string) {
					for _, elem := range runeNamePages {
						name := strings.Fields(s)
						if name[0]+" "+name[1] == elem.Name {
							setRuneHelper(elem.Page)
						}
					}
				}
				runeSelect.Refresh()
			}

			if len(result) > 0 {
				options := make([]string, len(result))

				for x, elem := range result {
					options[x] = elem[0] + " - Pick rate: " + elem[1]
				}

				p.Options = options
				p.Selected = options[0]
				p.OnChanged = func(s string) {
					sel := strings.TrimSpace(strings.Split(s, "-")[0])
					if lastRole != sel {
						for _, res := range result {
							if res[0] == sel {
								status.SetText("Setting...")
								resp, err := soup.Get(res[2])
								if err != nil {
									panic(err)
								}
								doc := soup.HTMLParse(resp)
								runeNamePages, runeDetails := setRunes(&doc, &(res[3]))

								if len(runeNamePages) > 0 {
									options := make([]string, len(runeNamePages))
									for x, elem := range runeNamePages {
										runeDetail := runeDetails[x]
										options[x] = elem.Name + " PR:" + runeDetail[0] + " WR:" + runeDetail[1] + " Sample:" + runeDetail[2]
										//options[x] = elem.Name + "  PR:" + runeDetail[0] + " WR:" + runeDetail[1]
									}
									runeSelect.Options = options
									runeSelect.Selected = options[0]
									runeSelect.OnChanged = func(s string) {
										for _, elem := range runeNamePages {
											name := strings.Fields(s)
											if name[0]+" "+name[1] == elem.Name {
												setRuneHelper(elem.Page)
											}
										}
									}
									runeSelect.Refresh()
								}

								setItems(&doc, accId, sumId, champId, &(res[3]))
								setSpells(&doc)
								lastRole = strings.TrimSpace(res[0])
								status.SetText("Updated...")
							}
						}
					}
				}
				p.Refresh()
			} else {
				p.PlaceHolder = "No alternative role available."
				p.Refresh()
			}
			prevChampId = champId
		}
		fmt.Println("Checking if Champion ID was updated...")
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
	p.Options = nil
	p.Refresh()

	status.SetText("Idle...")

	time.Sleep(30 * time.Second)
	inGame := getUxStatus()
	// Perhaps also check available?
	for inGame {
		fmt.Println("In game: ", inGame)
		time.Sleep(1 * time.Minute)
		inGame = getUxStatus()
	}
}

func runLoop(status *widget.Label, roleSelect *widget.Select, selectedChamp *widget.Label, runeSelect *widget.Select) {
	var wait sync.WaitGroup
	for {
		go run(status, roleSelect, selectedChamp, &wait, runeSelect)
		wait.Add(1)
		wait.Wait()
	}
}

func main() {
	// Initialize http client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli = &http.Client{Transport: tr}

	readConfig()
	checkFiles()

	a := app.New()
	w := a.NewWindow(ProjectName + " " + Version)

	w.SetOnClosed(func() {
		writeConfig()
		os.Exit(0)
	})

	// TODO: add write config, debug flag
	// TODO: add check if in game
	// TODO: Add infinite execution

	sl := widget.NewSlider(1, 5)
	sl.Value = config.Interval
	sl.Step = 0.5
	sl.OnChanged = func(f float64) {
		//fmt.Println("Poll interval: ", f)
		config.Interval = f
	}

	infoTextStyle := fyne.TextStyle{
		Bold:      true,
		Italic:    false,
		Monospace: true,
	}
	status := widget.NewLabelWithStyle("Not running", fyne.TextAlignCenter, infoTextStyle)
	selectedChamp := widget.NewLabelWithStyle("Not selected", fyne.TextAlignCenter, infoTextStyle)

	enableRunesCheck := widget.NewCheck("", func(b bool) {
		config.EnableRune = b
	})
	enableRunesCheck.SetChecked(config.EnableRune)

	enableItemsCheck := widget.NewCheck("", func(b bool) {
		config.EnableItem = b
	})
	enableItemsCheck.SetChecked(config.EnableItem)

	roleSelect := widget.NewSelect(nil, nil)
	roleSelect.PlaceHolder = "No champion selected"

	runeSelect := widget.NewSelect(nil, nil)
	runeSelect.PlaceHolder = "No rune selected"

	enableSpellCheck := widget.NewCheck("", func(b bool) {
		config.EnableSpell = b
	})
	enableSpellCheck.SetChecked(config.EnableSpell)

	enableDFlash := widget.NewCheck("", func(b bool) {
		config.DFlash = b
	})
	enableDFlash.SetChecked(config.DFlash)

	enableDebugging := widget.NewCheck("", func(b bool) {
		config.Debug = b
	})
	enableDebugging.SetChecked(config.Debug)

	go runLoop(status, roleSelect, selectedChamp, runeSelect)

	checkUpdateButton := widget.NewButton("Check Update", func() {
		req, err := http.Get("https://api.github.com/repos/jaeha-choi/DFF/releases/latest")
		if err != nil {
			panic(err)
		}

		if req.StatusCode != 200 {
			widget.ShowPopUpAtPosition(widget.NewLabel("Could not connect to github repository."),
				w.Canvas(), fyne.NewPos(50, 50))
		} else {
			var update GitHubUpdate
			err = json.NewDecoder(req.Body).Decode(&update)
			if err != nil {
				panic(err)
			}

			if update.TagName != Version {
				downloadUrl := ""
				for _, asset := range update.Assets {
					if asset.Name == "DFF_windows.zip" {
						downloadUrl = asset.BrowserDownloadURL
					}
				}

				if downloadUrl == "" {
					popup := widget.NewLabel("Update Error. File not found.")
					widget.ShowPopUpAtPosition(popup,
						w.Canvas(), fyne.NewPos(50, 50))
					time.Sleep(3 * time.Second)
				}

				name := strings.Split(downloadUrl, "/")
				out, err := os.Create(name[len(name)-1])
				if err != nil {
					panic(err)
				}
				defer out.Close()

				resp, err := http.Get(downloadUrl)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				_, err = io.Copy(out, resp.Body)
				if err != nil {
					panic(err)
				}

				fmt.Println(name[len(name)-1] + " downloaded.")
				popup := widget.NewVBox(widget.NewLabel("Version "+update.TagName+" downloaded."),
					widget.NewLabel("Program will now exit."))
				widget.ShowPopUpAtPosition(popup,
					w.Canvas(), fyne.NewPos(50, 50))
				time.Sleep(3 * time.Second)
				os.Exit(0)
			} else {
				widget.ShowPopUpAtPosition(widget.NewLabel("No update found"),
					w.Canvas(), fyne.NewPos(50, 50))
			}
		}
	})

	//ss := "123456789012345678901234"
	//output := widget.NewTextGridFromString(ss[:23]+"\n"+ss[23:])

	///lol-lobby/v1/lobby/availability
	///lol-lobby/v1/lobby/countdown
	///riotclient/get_region_locale

	w.SetContent(
		widget.NewVBox(
			widget.NewHBox(
				widget.NewVBox(
					widget.NewLabel(ProjectName+" "+Version),
					widget.NewLabel("Program Status:"),
					status,
					widget.NewHBox(widget.NewLabel("Current Champion:")),
					selectedChamp),
				widget.NewVBox(
					checkUpdateButton,
					widget.NewHBox(widget.NewLabel("Debug"), enableDebugging),
					widget.NewHBox(widget.NewLabel("Auto runes"), enableRunesCheck),
					widget.NewHBox(widget.NewLabel("Auto items"), enableItemsCheck),
					widget.NewHBox(widget.NewLabel("Auto spells"), enableSpellCheck),
					widget.NewHBox(widget.NewLabel("Left Flash"), enableDFlash),
					widget.NewLabel("Polling interval"),
					sl,
				)),
			roleSelect,
			runeSelect))

	//output))
	w.SetFixedSize(true)
	w.ShowAndRun()
}
