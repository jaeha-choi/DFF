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
const Version string = "v0.5.4"
const IssueUrl string = "https://github.com/jaeha-choi/DFF/issues"

type DFFClient struct {
	apiPort     string
	apiPass     string
	apiProtocol string
	lastRole    string
	Log         *log.Logger
	gameClient  *http.Client
	account     *datatype.AccountInfo
	cache       *cache.Cache

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

type RuneNamePage struct {
	Name string
	Page datatype.RunePage
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
		lastRole:    "",
		Log:         log.NewLogger(outTo, log.INFO, ""),
		gameClient: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}},
		cache:       nil, // must be initialized later
		account:     nil,
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

// requestApi is a function interface for game client API (GET methods)
func (client *DFFClient) requestApi(command *string) io.ReadCloser {
	req, err := http.NewRequest("GET", client.apiProtocol+"://127.0.0.1:"+client.apiPort+(*command), nil)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error encountered while requesting information")
		return nil
	}
	req.SetBasicAuth("riot", client.apiPass)

	resp, err := client.gameClient.Do(req)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while requesting API")
		return nil
	}

	return resp.Body
}

// isInChampSelect returns true if the user is currently in a champion select phase, false otherwise
func (client *DFFClient) isInChampSelect() (bool, error) {
	command := "/lol-champ-select/v1/session"
	var champSelect datatype.ChampSelect

	err := json.NewDecoder(client.requestApi(&command)).Decode(&champSelect)
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

	// Repeat until API is functional
	reader := client.requestApi(&command)
	for reader == nil {
		reader = client.requestApi(&command)
		time.Sleep(time.Duration(client.Interval) * time.Second)
	}

	if err = json.NewDecoder(reader).Decode(&client.account); err != nil {
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

// checkIsInGame returns true if a user is currently in a game, false otherwise
func (client *DFFClient) checkIsInGame() (bool, error) {
	command := "/riotclient/ux-state"

	bodyBytes, err := ioutil.ReadAll(client.requestApi(&command))
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

	err := json.NewDecoder(client.requestApi(&command)).Decode(&queueInfo)
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

// TODO: Add caching
// setItems sets an item page
func (client *DFFClient) setItems(doc *soup.Root, champId int, gameType *string) (bool, error) {
	if !client.EnableItem {
		return false, nil
	}
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

	itemPage := datatype.ItemPage{
		AccountID: client.account.AccountID,
		ItemSets: []datatype.ItemSet{
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

	command := "/lol-item-sets/v1/item-sets/" + strconv.Itoa(client.account.SummonerID) + "/sets"

	//j, _ := json.Marshal(newRune)
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(itemPage)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while setting items")
		return false, err
	}

	req, err := http.NewRequest("PUT", client.apiProtocol+"://127.0.0.1:"+client.apiPort+command, b)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while setting items")
		return false, err
	}

	req.SetBasicAuth("riot", client.apiPass)

	resp, err := client.gameClient.Do(req)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while setting items")
		return false, err
	}

	return resp.StatusCode == http.StatusCreated, err
}

// TODO: Add caching
// setSpells sets spells
func (client *DFFClient) setSpells(doc *soup.Root) (bool, error) {
	if !client.EnableSpell {
		return false, nil
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
		client.Log.Debug(err)
		client.Log.Error("Error while setting spells")
		return false, err
	}

	err = json.Unmarshal(f, &spellFile)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while setting spells")
		return false, err
	}

	client.Log.Debug("summoner.json version: ", spellFile.Version)
	client.Log.Debug("Total spell count: ", len(spellFile.Data))

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
	if client.DFlash && spellKeyList[1] == 4 {
		spellKeyList[1] = spellKeyList[0]
		spellKeyList[0] = 4
	} else if !client.DFlash && spellKeyList[0] == 4 {
		spellKeyList[0] = spellKeyList[1]
		spellKeyList[1] = 4
	}

	spells := datatype.Spells{
		Spell1ID: spellKeyList[0],
		Spell2ID: spellKeyList[1],
	}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(spells)

	command := "/lol-champ-select/v1/session/my-selection"
	req, err := http.NewRequest("PATCH", client.apiProtocol+"://127.0.0.1:"+client.apiPort+command, b)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while setting spells")
		return false, err
	}

	req.SetBasicAuth("riot", client.apiPass)

	resp, err := client.gameClient.Do(req)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while setting spells")
		return false, err
	}

	return resp.StatusCode == http.StatusNoContent, nil
}

// TODO: Add caching
// setRunePage set a rune page
func (client *DFFClient) setRunePage(page datatype.RunePage) (bool, error) {
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

	req, err := http.NewRequest("POST", client.apiProtocol+"://127.0.0.1:"+client.apiPort+command, b)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error in rune helper")
		return false, err
	}

	req.SetBasicAuth("riot", client.apiPass)

	resp, errs := client.gameClient.Do(req)
	if errs != nil {
		client.Log.Debug(err)
		client.Log.Error("Error in rune helper")
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

// delRunePage deletes a rune page created by DFF, or the first rune page
func (client *DFFClient) delRunePage() (deleted bool, err error) {
	var runePages datatype.RunePages
	var runePageCnt datatype.RunePageCount

	command := "/lol-perks/v1/pages"
	if err = json.NewDecoder(client.requestApi(&command)).Decode(&runePages); err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while getting rune pages from the client")
		return false, err
	}

	command = "/lol-perks/v1/inventory"
	if err = json.NewDecoder(client.requestApi(&command)).Decode(&runePageCnt); err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while getting total rune pages count")
		return false, err
	}

	client.Log.Debug("Total Rune pages: ", len(runePages))

	// Look for rune page starting with "DFF"
	for _, page := range runePages {
		client.Log.Debug("Current rune page: " + page.Name)
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

// TODO: Add caching
// setRunePageHelper will parse runes and make a RuneNamePage structure
func (client *DFFClient) setRunePageHelper(doc *soup.Root, gameType *string) ([]RuneNamePage, [][]string) {
	if !client.EnableRune {
		return nil, nil
	}

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

	links := (*doc).FindAll("div", "class", "perk-page-wrap")
	for x, link := range links {
		// Category
		imgs := link.FindAll("div", "class", "perk-page__item--mark")

		runeCategoryList := make([]int, len(imgs))

		if len(imgs) != 2 {
			client.Log.Error("Rune category updated? Please submit a new issue at " + IssueUrl)
			return nil, nil
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
			return nil, nil
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

		runeInfo[x].Page = datatype.RunePage{
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

	if ok, err := client.setRunePage(runeInfo[0].Page); !ok || err != nil {
		client.Log.Debug(err)
		client.Log.Error("Unable to set a rune page")
		return nil, nil
	}

	return runeInfo, runeDetails
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

	if err = json.NewDecoder(client.requestApi(&command)).Decode(&champSelect); err != nil {
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

// TODO: remove panic, code review, add caching
func (client *DFFClient) getRunes(queueId int, champId int, champLabel *widget.Label) ([][4]string, []RuneNamePage, [][]string) {
	command := "/lol-champions/v1/inventories/" + strconv.Itoa(client.account.SummonerID) + "/champions/" + strconv.Itoa(champId)
	var champion datatype.Champion
	var gameType, url string
	var posUrlList [][4]string = nil
	var runeNamePages []RuneNamePage = nil
	var runeDetails [][]string

	err := json.NewDecoder(client.requestApi(&command)).Decode(&champion)
	if err != nil {
		client.Log.Debug(err)
		client.Log.Error("Error while getting runes")
	}
	champLabel.SetText(champion.Alias)
	client.Log.Debug("Selected Champion: ", champion.Alias)

	// URF, ARURF
	if queueId == 900 {
		gameType = "URF"
		client.Log.Info("ULTRA RAPID FIRE MODE IS ON!!!")
		url = "https://op.gg/urf/" + champion.Alias + "/statistics"
		soup.Cookie("customLocale", client.Language)
		resp, err := soup.Get(url)
		if err != nil {
			panic(err)
		}

		doc := soup.HTMLParse(resp)

		runeNamePages, runeDetails = client.setRunePageHelper(&doc, &gameType)
		client.setItems(&doc, champId, &gameType)
		client.setSpells(&doc)

	} else if queueId == 450 {
		gameType = "ARAM"
		client.Log.Info("ARAM MODE IS ON!!!")
		url = "https://op.gg/aram/" + champion.Alias + "/statistics"
		soup.Cookie("customLocale", client.Language)
		resp, err := soup.Get(url)
		if err != nil {
			panic(err)
		}

		doc := soup.HTMLParse(resp)

		runeNamePages, runeDetails = client.setRunePageHelper(&doc, &gameType)
		client.setItems(&doc, champId, &gameType)
		client.setSpells(&doc)
	} else {
		// Can add region here
		url = "https://op.gg/champion/" + champion.Alias
		soup.Cookie("customLocale", client.Language)
		resp, err := soup.Get(url)
		if err != nil {
			panic(err)
		}

		doc := soup.HTMLParse(resp)

		runeNamePages, runeDetails = client.setRunePageHelper(&doc, &gameType)
		client.setItems(&doc, champId, &gameType)
		client.setSpells(&doc)

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

			client.Log.Debug(i, ". "+role+": ", rate)
			posUrlList[i][0] = strings.TrimSpace(role)
			posUrlList[i][1] = rate
			posUrlList[i][2] = link
			posUrlList[i][3] = gameType

		}
		client.lastRole = posUrlList[0][0]

	}
	return posUrlList, runeNamePages, runeDetails
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

// TODO: remove panic, code review, add caching
func (client *DFFClient) Run(window fyne.Window, status *widget.Label, p *widget.Select, champLabel *widget.Label, runeSelect *widget.Select) {
	defer func() {
		err := client.cache.SaveCache(filepath.Join("cache", "cache.bin"))
		if err != nil {
			client.Log.Debug(err)
			client.Log.Error("Error while saving cache")
			return
		}
	}()

	var err error

	status.SetText("Starting...")

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

	// For debug purpose; replace command to test the api
	//command := "/lol-perks/v1/show-auto-modified-pages-notification"
	//resp := requestApi(&command, false)
	//body, _ := ioutil.ReadAll(resp)
	//fmt.Println(string(body))

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

		if champId != 0 && prevChampId != champId {
			status.SetText("Setting...")
			result, runeNamePages, runeDetails := client.getRunes(queueId, champId, champLabel)
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
							client.setRunePage(elem.Page)
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
					if client.lastRole != sel {
						for _, res := range result {
							if res[0] == sel {
								status.SetText("Setting...")
								resp, err := soup.Get(res[2])
								if err != nil {
									panic(err)
								}
								doc := soup.HTMLParse(resp)
								runeNamePages, runeDetails := client.setRunePageHelper(&doc, &(res[3]))

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
												client.setRunePage(elem.Page)
											}
										}
									}
									runeSelect.Refresh()
								}

								client.setItems(&doc, champId, &(res[3]))
								client.setSpells(&doc)
								client.lastRole = strings.TrimSpace(res[0])
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
		client.Log.Debug("Checking if Champion ID was updated...")
		time.Sleep(time.Duration(client.Interval) * time.Second)
	}
	p.Options = nil
	client.lastRole = ""
	p.Refresh()

	status.SetText("Idle...")
	client.Log.Debug("User is in the game")

	var isInGame = true
	for isInGame {
		time.Sleep(1 * time.Minute)
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
