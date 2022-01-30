package core

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/anaskhan96/soup"
	"github.com/jaeha-choi/DFF/internal/cache"
	"github.com/jaeha-choi/DFF/internal/datatype"
	"github.com/jaeha-choi/DFF/pkg/log"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const ProjectName string = "DFF!"
const Version string = "v0.5.5-beta"
const IssueUrl string = "https://github.com/jaeha-choi/DFF/issues"

type DFFClient struct {
	apiPort     string
	apiPass     string
	apiProtocol string
	Log         *log.Logger
	gameClient  *http.Client
	account     *datatype.AccountInfo
	cache       *cache.Cache
	window      fyne.Window
	gameVersion string

	Debug       bool    `json:"debug"`
	Interval    float64 `json:"interval"`
	ClientDir   string  `json:"client_dir"`
	EnableRune  bool    `json:"enable_rune"`
	EnableItem  bool    `json:"enable_item"`
	EnableSpell bool    `json:"enable_spell"`
	DFlash      bool    `json:"d_flash"`
	Language    string  `json:"language"`
}

type SpellFile struct {
	Type    string                 `json:"type"`
	Version string                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}

type FileVersion struct {
	Version string `json:"version"`
}

// Initialize creates DFFClient structure and initialize files/variables
func Initialize(outTo io.Writer) (client *DFFClient) {
	var err error
	client = createDFFClient(outTo)

	if err = client.readConfig("config.json"); err != nil {
		client.Log.Warning(ProjectName + " may not be initialized properly")
	}

	if err = client.WriteConfig(); err != nil {
		client.Log.Error("Could not write config file")
	}

	if err = os.MkdirAll("cache", 0700); err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while creating cache folder")
	}

	if client.cache, err = cache.RestoreCache(filepath.Join("cache", "cache.bin")); err != nil {
		client.Log.Debug(err)
		client.Log.Warning("Could not restore cache")
		client.cache = cache.NewCache()
	}

	// Read/Download/Sync mandatory files if necessary
	if err = client.checkFiles(); err != nil {
		client.Log.Error("At least one mandatory file is missing")
	}

	return client
}

// createDFFClient initializes the DFF client and variables used by it
func createDFFClient(outTo io.Writer) *DFFClient {
	return &DFFClient{
		apiPort:     "",
		apiPass:     "",
		apiProtocol: "",
		Log:         log.NewLogger(outTo, log.INFO, ""),
		gameClient: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}},
		account:     nil,
		cache:       nil, // must be initialized later
		window:      nil,
		gameVersion: "",
		Debug:       false,
		Interval:    2,
		ClientDir:   "C:/Riot Games/League of Legends/",
		EnableRune:  true,
		EnableItem:  true,
		EnableSpell: true,
		DFlash:      true,
		Language:    "en_US",
	}
}

// readConfig reads configuration file if it exist and update variables, or use default config otherwise
func (client *DFFClient) readConfig(filename string) (err error) {
	if _, err = os.Stat(filename); err == nil {
		var fileBytes []byte

		if fileBytes, err = ioutil.ReadFile(filename); err != nil {
			client.Log.Debug(err)
			client.Log.Error("Error while reading ", filename)
			return err
		}

		if err = json.Unmarshal(fileBytes, &client); err != nil {
			client.Log.Debug(err)
			client.Log.Error("Error while loading a configuration file")
			return err
		}

		if client.Interval < 1 {
			client.Interval = 1
		} else if client.Interval > 5 {
			client.Interval = 5
		}

		if client.Debug {
			client.Log.Mode = log.DEBUG
		} else {
			client.Log.Mode = log.INFO
		}
	} else {
		client.Log.Debug(err)
		client.Log.Warning("Cannot open config file. Default settings will be used.")
	}
	return err
}

// WriteConfig writes configuration file
func (client *DFFClient) WriteConfig() (err error) {
	var jsonConf []byte

	if jsonConf, err = json.MarshalIndent(client, "", "\t"); err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while encoding client to bytes")
		return err
	}

	if err = ioutil.WriteFile("config.json", jsonConf, 0644); err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while writing a configuration file")
		return err
	}

	return err
}

// readLockFile wait for lockfile to be generated and reads "lockfile", which provides a token to access
// the game client. Returns the content of lockfile as string or err if failed
func (client *DFFClient) readLockFile() (err error) {
	var file *os.File
	// Loop until file is available
	for {
		if file, err = os.Open(client.ClientDir + "lockfile"); err == nil {
			client.Log.Debug("lockfile found")
			break
		} else {
			client.Log.Debug(err)
			client.Log.Info("Waiting for League process to open")
		}
		time.Sleep(time.Duration(client.Interval) * time.Second)
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while reading a lockfile")
		return err
	}

	// lockfileValues[0] = not used
	// lockfileValues[1] = not used
	// lockfileValues[2] = Port
	// lockfileValues[3] = API auth password (username is always "riot")
	// lockfileValues[4] = Protocol (https)
	lockfileValues := strings.Split(string(b), ":")

	//for _, val := range lockfileValues {
	//	client.Log.Debug(val)
	//}

	client.apiPort = lockfileValues[2]
	client.apiPass = lockfileValues[3]
	client.apiProtocol = lockfileValues[4]

	return err
}

// requestApi is a function interface for game client API
func (client *DFFClient) requestApi(method string, command string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, client.apiProtocol+"://127.0.0.1:"+client.apiPort+command, body)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error encountered while requesting information")
		return nil
	}
	req.SetBasicAuth("riot", client.apiPass)

	resp, err := client.gameClient.Do(req)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while requesting API:", command)
		return nil
	}

	return resp
}

// isInChampSelect returns true if the user is currently in a champion select phase, false otherwise
func (client *DFFClient) isInChampSelect() (bool, error) {
	command := "/lol-champ-select/v1/session"
	var champSelect datatype.ChampSelect

	err := json.NewDecoder(client.requestApi("GET", command, nil).Body).Decode(&champSelect)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while decoding API response")
		return false, err
	}

	return float64(champSelect.Timer.AdjustedTimeLeftInPhase) > client.Interval, err
}

// getAccInfo returns login information
func (client *DFFClient) getAccInfo() (err error) {
	command := "/lol-summoner/v1/current-summoner"

	req, err := http.NewRequest("GET", client.apiProtocol+"://127.0.0.1:"+client.apiPort+(command), nil)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error encountered while requesting information")
		return err
	}
	req.SetBasicAuth("riot", client.apiPass)

	var resp *http.Response
	// Repeat until API is functional
	for {
		if resp, err = client.gameClient.Do(req); err == nil && resp.StatusCode == 200 {
			break
		}
		client.Log.Debug(err)
		if resp != nil {
			client.Log.Debug("Account Info Status code: ", resp.StatusCode)
		}
		time.Sleep(1 * time.Second)
	}

	if err = json.NewDecoder(resp.Body).Decode(&client.account); err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while decoding account information")
		return err
	}

	//client.Log.Debug("Logged in as...")
	//client.Log.Debug("Account ID: ", client.account.AccountID)
	//client.Log.Debug("Display Name: ", client.account.DisplayName)
	//client.Log.Debug("Internal Name: ", client.account.InternalName)
	//client.Log.Debug("Player UUID: ", client.account.Puuid)
	//client.Log.Debug("Summoner ID: ", client.account.SummonerID)
	client.Log.Info("Client API connection functional.")

	return nil
}

// TODO: Need to find a better API for this operation
// checkIsInGame returns true if a user is currently in a game, false otherwise
func (client *DFFClient) checkIsInGame() (bool, error) {
	command := "/riotclient/ux-state"

	bodyBytes, err := ioutil.ReadAll(client.requestApi("GET", command, nil).Body)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while checking if the user is in a game")
		return false, err
	}
	return string(bodyBytes) != "\"ShowMain\"", nil
}

// getQueueId returns the type of the game (normal, urf, aram, etc)
func (client *DFFClient) getQueueId() (int, error) {
	command := "/lol-gameflow/v1/gameflow-metadata/player-status"
	var queueInfo datatype.QueueInfo

	err := json.NewDecoder(client.requestApi("GET", command, nil).Body).Decode(&queueInfo)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while getting the queue type")
		return -1, err
	}
	return queueInfo.CurrentLobbyStatus.QueueID, err
}

// deleteRunePageWithId deletes old rune page and return true if deleted, false otherwise
func (client *DFFClient) deleteRunePageWithId(runePageId int) (bool, error) {
	command := "/lol-perks/v1/pages/"

	req, err := http.NewRequest("DELETE", client.apiProtocol+"://127.0.0.1:"+
		client.apiPort+command+strconv.Itoa(runePageId), nil)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while deleting an old DFF rune page")
		return false, err
	}
	req.SetBasicAuth("riot", client.apiPass)

	resp, err := client.gameClient.Do(req)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while requesting to delete the old DFF rune page")
		return false, err
	}

	return resp.StatusCode == http.StatusNoContent, nil
}

// retrieveItems sets an item page
func (client *DFFClient) retrieveItems(doc *soup.Root, cachedData *cache.CachedData, champId int, gameType string) (isSet bool) {
	builds := (*doc).FindAll("tr", "class", "champion-overview__row")
	blockCnt := len((*doc).FindAll("tr", "class", "champion-overview__row--first")) + 1

	blockList := make([]datatype.ItemBlock, blockCnt)
	otherItemSet := make(map[string]bool)
	willBeAdded := 0
	i := 0
	for _, build := range builds {
		if strings.HasSuffix(build.Attrs()["class"], "champion-overview__row--first") {
			items := build.FindAll("li", "class", "champion-stats__list__item")
			itemList := make([]datatype.Item, len(items))
			for j, img := range items {
				str := img.Find("img").Attrs()["src"]
				str = str[strings.LastIndex(str, "/")+1 : strings.Index(str, ".png")]
				newItem := datatype.Item{
					Count: 1,
					ID:    str,
				}
				otherItemSet[str] = false
				itemList[j] = newItem
			}
			newItemBlock := datatype.ItemBlock{
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

	ward := datatype.Item{
		Count: 1,
		ID:    "3340",
	}

	blockList[0].Items = append(blockList[0].Items, ward)

	consumable := datatype.ItemBlock{
		HideIfSummonerSpell: "",
		Items: []datatype.Item{
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

	client.Log.Debug("Total number of items:", len(otherItemSet))
	client.Log.Debug("Count of items that will be added:", willBeAdded)

	itemList := make([]datatype.Item, willBeAdded)

	idx := 0
	for otherItem := range otherItemSet {
		if otherItemSet[otherItem] {
			newItem := datatype.Item{
				Count: 1,
				ID:    otherItem,
			}
			itemList[idx] = newItem
			idx++
		}
	}

	newItemBlock := datatype.ItemBlock{
		HideIfSummonerSpell: "",
		Items:               itemList,
		ShowIfSummonerSpell: "",
		Type:                "Other items to consider",
	}

	blockList[i] = newItemBlock
	i++

	cachedData.ItemPages.AccountID = client.account.AccountID
	cachedData.ItemPages.ItemSets = []datatype.ItemSet{
		{
			AssociatedChampions: []int{champId},
			AssociatedMaps:      []int{11, 12},
			Blocks:              blockList,
			Map:                 "any",
			Mode:                "any",
			PreferredItemSlots:  []interface{}{},
			Sortrank:            0,
			StartedFrom:         "blank",
			Title:               ProjectName + " Item Page " + gameType,
			Type:                "custom",
			UID:                 "",
		},
	}
	cachedData.ItemPages.Timestamp = 0

	return true
}

// retrieveSpells sets spells
func (client *DFFClient) retrieveSpells(doc *soup.Root, cachedData *cache.CachedData) (isSet bool) {
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
		client.Log.Debug(err)
		client.Log.Error("Error while setting spells")
		return false
	}

	err = json.Unmarshal(f, &spellFile)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while setting spells")
		return false
	}

	client.Log.Debug("summoner.json version: ", spellFile.Version)
	client.Log.Debug("Total spell count: ", len(spellFile.Data))

	spellKeyList := make([]int64, len(spellImgs))
	for x := 0; x < 10; x++ {

		i := 0
		for _, spells := range spellFile.Data {
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
	if client.DFlash && spellKeyList[1] == 4 {
		spellKeyList[1] = spellKeyList[0]
		spellKeyList[0] = 4
	} else if !client.DFlash && spellKeyList[0] == 4 {
		spellKeyList[0] = spellKeyList[1]
		spellKeyList[1] = 4
	}

	cachedData.Spells.Spell1ID = spellKeyList[0]
	cachedData.Spells.Spell2ID = spellKeyList[1]

	return true
}

// setRunePage set a rune page
func (client *DFFClient) setRunePage(page *datatype.RunePage) (bool, error) {
	command := "/lol-perks/v1/pages"

	if ok, err := client.delRunePage(); !ok || err != nil {
		client.Log.Debug(err)
		client.Log.Warning("Old rune page not deleted")
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(page)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error in rune helper")
		return false, err
	}

	req := client.requestApi("POST", command, b)
	if req == nil || req.StatusCode != http.StatusOK {
		client.Log.Debug(err)
		client.Log.Error("Error while setting items")
		return false, nil
	}

	return req.StatusCode == http.StatusOK, nil
}

// delRunePage deletes a rune page created by DFF, or the first rune page
func (client *DFFClient) delRunePage() (deleted bool, err error) {
	var runePages datatype.RunePages
	var runePageCnt datatype.RunePageCount

	command := "/lol-perks/v1/pages"
	if err = json.NewDecoder(client.requestApi("GET", command, nil).Body).Decode(&runePages); err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while getting rune pages from the client")
		return false, err
	}

	command = "/lol-perks/v1/inventory"
	if err = json.NewDecoder(client.requestApi("GET", command, nil).Body).Decode(&runePageCnt); err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while getting total rune pages count")
		return false, err
	}

	client.Log.Debug("Total Rune pages: ", len(runePages))

	// Look for rune page starting with "DFF"
	for _, page := range runePages {
		//client.Log.Debug("Iterating rune page: " + page.Name)
		if strings.HasPrefix(page.Name, ProjectName) {
			if ok, err := client.deleteRunePageWithId(page.ID); ok && err == nil {
				deleted = true
			}
		}
	}

	// Delete the first rune page if all pages are used (excluding 5 default rune pages)
	if !deleted && len(runePages)+5 >= runePageCnt.OwnedPageCount {
		if ok, err := client.deleteRunePageWithId(runePages[0].ID); ok && err == nil {
			deleted = true
		}
	}

	return deleted, nil
}

// retrieveRunes will parse runes and make a RuneNamePage structure
func (client *DFFClient) retrieveRunes(doc *soup.Root, cachedData *cache.CachedData, gameType string) (isSet bool) {
	runeDetailsDoc := (*doc).FindAll("span", "class", "pick-ratio__text")

	// Getting Pick rate/Win rate/Sample count
	var pr, wr, sample string
	for idx, runeDetailDoc := range runeDetailsDoc {
		next := runeDetailDoc.FindNextElementSibling()
		pr = next.Text()
		next = next.FindNextElementSibling()
		sample = next.Text()
		next = next.FindNextElementSibling()
		next = next.FindNextElementSibling()
		wr = next.Text()
		cachedData.RunePages[idx].PickRate = pr
		cachedData.RunePages[idx].WinRate = wr
		cachedData.RunePages[idx].SampleCnt = sample
	}

	// Creating rune page name
	runeNames := (*doc).FindAll("div", "class", "champion-stats-summary-rune__name")
	i := 0
	for _, runeName := range runeNames {
		names := strings.Split(runeName.Text(), "+")
		//fmt.Println(runeName.Text())
		for x := 0; x < 2; x++ {
			cachedData.RunePages[i].Name = string([]rune(strings.TrimSpace(names[0]))[0]) + "+" + string([]rune(strings.TrimSpace(names[1]))[0]) + " (" + strconv.Itoa(x+1) + ")"
			i++
			//fmt.Println(runeInfo[x].Name)
		}
	}

	// Creating rune page
	links := (*doc).FindAll("div", "class", "perk-page-wrap")
	for x, link := range links {
		// Category
		imgs := link.FindAll("div", "class", "perk-page__item--mark")

		runeCategoryList := make([]int, len(imgs))

		if len(imgs) != 2 {
			client.Log.Error("Rune category updated? Please submit a new issue at " + IssueUrl)
			return false
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
			client.Log.Error("Runes updated? Please submit a new issue at " + IssueUrl)
			return false
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

		cachedData.RunePages[x].Page = datatype.RunePage{
			AutoModifiedSelections: []interface{}{},
			Current:                true,
			ID:                     0,
			IsActive:               true,
			IsDeletable:            true,
			IsEditable:             true,
			IsValid:                true,
			LastModified:           0,
			Name:                   ProjectName + " " + cachedData.RunePages[x].Name + " " + gameType,
			Order:                  0,
			PrimaryStyleID:         runeCategoryList[0],
			SelectedPerkIds:        runeList,
			SubStyleID:             runeCategoryList[1],
		}
	}

	return true
}

func (client *DFFClient) downloadFile(url string) error {
	fileName := strings.Split(url, "/")
	out, err := os.Create("./data/" + fileName[len(fileName)-1])
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while creating a file")
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while downloading a file")
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while saving a file")
		return err
	}
	client.Log.Info(fileName[len(fileName)-1] + " downloaded.")

	return nil
}

func (client *DFFClient) getChampId() (champId int, err error) {
	command := "/lol-champ-select/v1/session"
	var champSelect datatype.ChampSelect

	if err = json.NewDecoder(client.requestApi("GET", command, nil).Body).Decode(&champSelect); err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while getting champion ID")
		return 0, err
	}

	// Find current user's champion ID
	for _, member := range champSelect.MyTeam {
		if member.SummonerID == client.account.SummonerID {
			champId = member.ChampionID
			break
		}
	}

	return champId, err
}

func (client *DFFClient) retrieveData(gameMode datatype.GameMode, champion *datatype.Champion, champLabel *widget.Label, position cache.Position) (cachedData *cache.CachedData, ok bool) {
	var gameType, url string
	var err error
	champLabel.SetText(champion.Alias)
	client.Log.Debug("Selected Champion: ", champion.Alias)

	cacheData, isCached := client.cache.GetPut(champion.Alias, gameMode, position)
	client.Log.Debug("Using cache: ", isCached)
	if !isCached {
		switch gameMode {
		case datatype.ARAM:
			gameType = "ARAM"
			client.Log.Info("ARAM MODE IS ON!!!")
			url = "https://op.gg/aram/" + champion.Alias + "/statistics"
		case datatype.URF:
			gameType = "URF"
			client.Log.Info("ULTRA RAPID FIRE MODE IS ON!!!")
			url = "https://op.gg/urf/" + champion.Alias + "/statistics"
		case datatype.DEFAULT:
			url = "https://op.gg/champion/" + champion.Alias
			if position != cache.NONE {
				node, _ := client.cache.GetPutNode(champion.Alias)
				url = node.Value.Default[position].URL
			}
		}
		soup.Cookie("customLocale", client.Language)

		var resp string
		retryCnt := 3
		for i := 0; i < retryCnt; i++ {
			resp, err = soup.Get(url)
			if err == nil {
				break
			} else if i == retryCnt-1 {
				client.Log.Debug(err)
				client.Log.Error("Couldn't connect to op.gg")
				return nil, false
			}
			time.Sleep(500 * time.Millisecond)
		}

		doc := soup.HTMLParse(resp)

		if gameMode == datatype.DEFAULT && position == cache.NONE {
			// Find champion positions
			positions := doc.FindAll("li", "class", "champion-stats-header__position")
			node, _ := client.cache.GetPutNode(champion.Alias)
			node.Value.AvailablePositions = make([]cache.Position, len(positions))

			for i, pos := range positions {
				link := "https://op.gg" + pos.Find("a").Attrs()["href"]
				roleStr := strings.TrimSpace(pos.Find("span", "class", "champion-stats-header__position__role").Text())
				rate := pos.Find("span", "class", "champion-stats-header__position__rate").Text()

				client.Log.Debug(i, ". "+roleStr+": ", rate)

				var role cache.Position

				switch roleStr {
				case "Top":
					role = cache.TOP
				case "Jungle":
					role = cache.JUNGLE
				case "Middle":
					role = cache.MID
				case "Bottom":
					role = cache.ADC
				case "Support":
					role = cache.SUPPORT
				default:
					client.Log.Error("Role changed? Please submit a new issue at " + IssueUrl)
					return nil, false
				}
				node.Value.AvailablePositions[i] = role

				cacheData, _ = client.cache.GetPut(champion.Alias, datatype.DEFAULT, role)
				cacheData.CreationTime = time.Now()
				cacheData.PositionPickRate = rate
				cacheData.URL = link
			}
			node.Value.DefaultPosition = node.Value.AvailablePositions[0]
			cacheData, _ = client.cache.GetPut(champion.Alias, datatype.DEFAULT, node.Value.DefaultPosition)
		}

		cacheData.Version = client.gameVersion
		cacheData.RunePages = make([]datatype.DFFRunePage, 4)

		isSet := client.retrieveRunes(&doc, cacheData, gameType)
		if !isSet {
			client.Log.Error("Error while caching rune page")
		}

		isSet = client.retrieveItems(&doc, cacheData, champion.ID, gameType)
		if !isSet {
			client.Log.Error("Error while caching item page")
		}

		isSet = client.retrieveSpells(&doc, cacheData)
		if !isSet {
			client.Log.Error("Error while caching spell page")
		}
	}

	if client.EnableRune {
		if ok, err := client.setRunePage(&cacheData.RunePages[0].Page); !ok || err != nil {
			client.Log.Debug(err)
			client.Log.Error("Unable to set a rune page")
			return nil, false
		}
	}

	if client.EnableItem {
		command := "/lol-item-sets/v1/item-sets/" + strconv.Itoa(client.account.SummonerID) + "/sets"

		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(cacheData.ItemPages)
		if err != nil {
			client.Log.Debug(err)
			client.Log.Error("Error while setting items")
			return nil, false
		}

		req := client.requestApi("PUT", command, b)
		if req == nil || req.StatusCode != http.StatusCreated {
			client.Log.Debug(err)
			client.Log.Error("Error while setting items")
			return nil, false
		}

	}

	if client.EnableSpell {
		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(cacheData.Spells)
		if err != nil {
			client.Log.Debug(err)
			client.Log.Error("Error while setting spells")
			return nil, false
		}

		command := "/lol-champ-select/v1/session/my-selection"
		req := client.requestApi("PATCH", command, b)
		if req == nil || req.StatusCode != http.StatusNoContent {
			client.Log.Debug(err)
			client.Log.Error("Error while setting spells")
			return nil, false
		}
	}

	return cacheData, true
}

// TODO: remove panic, code review
// checkFiles checks essential files and syncs if outdated
func (client *DFFClient) checkFiles() (err error) {
	var version []string

	// Get version list
	resp, err := http.Get("https://ddragon.leagueoflegends.com/api/versions.json")
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while checking the version")
		return err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&version); err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while decoding versions.json file")
		return err
	}

	// First index contains the latest version (e.g. "12.1.1")
	leagueVersion := version[0]

	if _, err = os.Stat("./data"); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Join(".", "data"), 0700); err != nil {
			client.Log.Debug(err)
			client.Log.Error("Error while downloading a file")
		}

		// Summoner spells
		client.downloadFile("https://ddragon.leagueoflegends.com/cdn/" + leagueVersion + "/data/en_US/" + "summoner.json")
		// Items
		//downloadFile("https://ddragon.leagueoflegends.com/cdn/"+ leagueVersion +"/data/en_US/"+"item.json")
		// Maps
		//downloadFile("https://ddragon.leagueoflegends.com/cdn/"+ leagueVersion +"/data/en_US/"+"map.json.json")
		// Runes
		//downloadFile("https://ddragon.leagueoflegends.com/cdn/"+ leagueVersion +"/data/en_US/"+"runesReforged.json")
		// Champions
		//downloadFile("https://ddragon.leagueoflegends.com/cdn/"+ leagueVersion +"/data/en_US/"+"champion.json")
	} else if err != nil {
		var files []fs.FileInfo
		var ver FileVersion
		var f []byte

		if files, err = ioutil.ReadDir("./data"); err != nil {
			client.Log.Debug(err)
			client.Log.Error("Error while reading directory")
			return err
		}

		for _, file := range files {
			if f, err = ioutil.ReadFile("./data/" + file.Name()); err != nil {
				client.Log.Debug(err)
				client.Log.Error("Error while reading file: " + file.Name())
				continue
			}

			if err = json.Unmarshal(f, &ver); err != nil {
				client.Log.Debug(err)
				client.Log.Error("Error while decoding json file: " + file.Name())
				continue
			}

			if leagueVersion != ver.Version {
				//if err = os.Remove("./data/" + file.Name()); err != nil {
				//	client.log.Debug(err)
				//	client.log.Error("Error while deleting outdated file: " + file.Name())
				//}
				client.downloadFile("https://ddragon.leagueoflegends.com/cdn/" + leagueVersion + "/data/en_US/" + file.Name())
			}
		}
	}
	return err
}

// Run starts DFF
func (client *DFFClient) Run(window fyne.Window, status *widget.Label, p *widget.Select, champLabel *widget.Label, runeSelect *widget.Select) {
	defer func() {
		client.Log.Debug("Saving cache...")
		err := client.cache.SaveCache(filepath.Join("cache", "cache.bin"))
		client.Log.Debug("Cache saved")
		if err != nil {
			client.Log.Debug(err)
			client.Log.Error("Error while saving cache")
			return
		}
	}()

	var err error

	status.SetText("Starting...")
	champLabel.SetText("Not selected")

	if err = client.readLockFile(); err != nil {
		status.SetText("Error. Check log")
		window.RequestFocus()
		return
	}

	if err = client.getAccInfo(); err != nil {
		status.SetText("Error. Check log")
		window.RequestFocus()
		return
	}

	//var isCustomGame = false
	var prevChampId, champId int
	var queueId = -1

	// Check if in lobby
	for champId == 0 {
		status.SetText("Waiting...")
		client.Log.Debug("Waiting for a champion to be selected...")
		if queueId, err = client.getQueueId(); err != nil {
			status.SetText("Error. Check log")
			window.RequestFocus()
		}
		if champId, err = client.getChampId(); err != nil {
			status.SetText("Error. Check log")
			window.RequestFocus()
		}
		time.Sleep(time.Duration(client.Interval) * time.Second)
	}

	var gameMode datatype.GameMode
	if queueId == int(datatype.ARAM) || queueId == int(datatype.URF) {
		gameMode = datatype.GameMode(queueId)
	} else {
		gameMode = datatype.DEFAULT
	}

	lastRole := cache.NONE
	position := cache.NONE
	positionIdx := 0
	var isInChampSelect = true
	for isInChampSelect {
		if isInChampSelect, err = client.isInChampSelect(); err != nil {
			status.SetText("Error. Check log")
			window.RequestFocus()
		}

		if champId, err = client.getChampId(); err != nil {
			status.SetText("Error. Check log")
			window.RequestFocus()
		}

		if champId != 0 && prevChampId != champId || lastRole != position {
			if prevChampId != champId {
				position = cache.NONE
				positionIdx = 0
			}
			// Convert champ id to datatype.Champion
			var champion datatype.Champion
			command := "/lol-champions/v1/inventories/" + strconv.Itoa(client.account.SummonerID) + "/champions/" + strconv.Itoa(champId)
			err = json.NewDecoder(client.requestApi("GET", command, nil).Body).Decode(&champion)
			if err != nil {
				client.Log.Debug(err)
				client.Log.Error("Error while getting runes")
				champLabel.SetText("Error. Check log")
				if client.window != nil {
					client.window.RequestFocus()
				}
			}

			status.SetText("Setting...")
			cachedData, ok := client.retrieveData(gameMode, &champion, champLabel, position)
			if !ok {
				champLabel.SetText("Error. Check log")
				if client.window != nil {
					client.window.RequestFocus()
				}
			}
			lastRole = position

			status.SetText("Updated...")

			if len(cachedData.RunePages) > 0 {
				runeSelect.Options = make([]string, len(cachedData.RunePages))
				for x, elem := range cachedData.RunePages {
					runeSelect.Options[x] = elem.Name + " PR:" + elem.PickRate + " WR:" + elem.WinRate + " Sample:" + elem.SampleCnt
				}
				runeSelect.Selected = runeSelect.Options[0]
				runeSelect.OnChanged = func(s string) {
					for _, elem := range cachedData.RunePages {
						name := strings.Fields(s)
						if name[0]+" "+name[1] == elem.Name {
							ok, err := client.setRunePage(&elem.Page)
							if !ok || err != nil {
								champLabel.SetText("Error. Check log")
								if client.window != nil {
									client.window.RequestFocus()
								}
							}
						}
					}
				}
				runeSelect.Refresh()
			}

			if gameMode == datatype.DEFAULT {
				node, _ := client.cache.GetPutNode(champion.Alias)
				if lastRole == cache.NONE {
					position = node.Value.DefaultPosition
					lastRole = position
				}
				p.Options = make([]string, len(node.Value.AvailablePositions))
				for i := 0; i < len(node.Value.AvailablePositions); i++ {
					pos := node.Value.AvailablePositions[i]
					switch pos {
					case cache.TOP:
						p.Options[i] += "Top"
					case cache.JUNGLE:
						p.Options[i] += "Jungle"
					case cache.MID:
						p.Options[i] += "Middle"
					case cache.ADC:
						p.Options[i] += "Bottom"
					case cache.SUPPORT:
						p.Options[i] += "Support"
					}
					p.Options[i] += " - Pick rate: " + node.Value.Default[pos].PositionPickRate
				}
				p.Selected = p.Options[positionIdx]
				p.OnChanged = func(s string) {
					var res string
					for positionIdx, res = range p.Options {
						if s == res {
							position = node.Value.AvailablePositions[positionIdx]
							break
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
		client.Log.Debug("Checking if Champion ID was updated...")
		time.Sleep(time.Duration(client.Interval) * time.Second)
	}
	p.Options = nil
	p.Refresh()

	status.SetText("Idle...")
	client.Log.Debug("User is in the game")

	var isInGame = true
	for isInGame {
		time.Sleep(30 * time.Second)
		client.Log.Debug("In game: ", isInGame)
		if isInGame, err = client.checkIsInGame(); err != nil {
			status.SetText("Error. Check log")
		}
	}
}

//command := "/lol-summoner/v1/current-summoner" // returns login information
//command := "/lol-champ-select/v1/session" // champion select session information
//command := "/riotclient/auth-token" // returns auth token
//command := "/riotclient/region-locale" // returns region, language etc
//command := "/lol-perks/v1/currentpage" // get current rune page
///lol-lobby/v1/lobby/availability
///lol-lobby/v1/lobby/countdown
///riotclient/get_region_locale
//command := "/lol-perks/v1/show-auto-modified-pages-notification"
