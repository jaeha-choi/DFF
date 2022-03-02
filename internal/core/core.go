package core

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/jaeha-choi/DFF/internal/cache"
	"github.com/jaeha-choi/DFF/internal/datatype"
	"github.com/jaeha-choi/DFF/pkg/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const ProjectName string = "DFF!"
const Version string = "v0.6.3 (beta 2)"
const IssueUrl string = "https://github.com/jaeha-choi/DFF/issues"

type DFFClient struct {
	apiPort     string
	apiPass     string
	apiProtocol string
	Log         *log.Logger
	gameClient  *http.Client
	account     *datatype.AccountInfo
	cache       *cache.Cache
	metaInfo    *Meta
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

	// Read/Download/Sync mandatory files if necessary
	if err = client.checkFiles(); err != nil {
		client.Log.Error("At least one mandatory file is missing")
	}

	if client.cache, err = cache.RestoreCache(filepath.Join("cache", "cache.bin"), client.gameVersion); err != nil {
		client.Log.Debug(err)
		client.Log.Warning("Could not restore cache, creating a new cache")
		client.cache = cache.NewCache(client.gameVersion)
	}

	if err = client.restoreChampionList(filepath.Join("cache", "positions.bin"), client.gameVersion); err != nil {
		client.Log.Debug(err)
		client.Log.Warning("Could not restore position data, attempting to download new position data")
		if !client.createChampionList(client.gameVersion) {
			client.Log.Error("Failed to get champion list")
			os.Exit(1)
		}
		if err = client.saveChampionList(filepath.Join("cache", "positions.bin")); err != nil {
			client.Log.Debug(err)
			client.Log.Error("Failed to save champion list")
		}
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
		Debug:       false,
		Interval:    2,
		ClientDir:   DefaultInstallPath,
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

// TODO: remove panic, code review
// checkFiles performs a version check and checks essential files and syncs if outdated
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
	client.gameVersion = version[0]

	//if _, err = os.Stat("./data"); os.IsNotExist(err) {
	//	if err = os.MkdirAll(filepath.Join(".", "data"), 0700); err != nil {
	//		client.Log.Debug(err)
	//		client.Log.Error("Error while downloading a file")
	//	}
	//
	//	// Summoner spells
	//	client.downloadFile("https://ddragon.leagueoflegends.com/cdn/" + client.gameVersion + "/data/en_US/" + "summoner.json")
	//	// Items
	//	//downloadFile("https://ddragon.leagueoflegends.com/cdn/"+ client.gameVersion +"/data/en_US/"+"item.json")
	//	// Maps
	//	//downloadFile("https://ddragon.leagueoflegends.com/cdn/"+ client.gameVersion +"/data/en_US/"+"map.json.json")
	//	// Runes
	//	//downloadFile("https://ddragon.leagueoflegends.com/cdn/"+ client.gameVersion +"/data/en_US/"+"runesReforged.json")
	//	// Champions
	//	//downloadFile("https://ddragon.leagueoflegends.com/cdn/"+ client.gameVersion +"/data/en_US/"+"champion.json")
	//} else if err != nil {
	//	var files []fs.FileInfo
	//	var ver FileVersion
	//	var f []byte
	//
	//	if files, err = ioutil.ReadDir("./data"); err != nil {
	//		client.Log.Debug(err)
	//		client.Log.Error("Error while reading directory")
	//		return err
	//	}
	//
	//	for _, file := range files {
	//		if f, err = ioutil.ReadFile("./data/" + file.Name()); err != nil {
	//			client.Log.Debug(err)
	//			client.Log.Error("Error while reading file: " + file.Name())
	//			continue
	//		}
	//
	//		if err = json.Unmarshal(f, &ver); err != nil {
	//			client.Log.Debug(err)
	//			client.Log.Error("Error while decoding json file: " + file.Name())
	//			continue
	//		}
	//
	//		if client.gameVersion != ver.Version {
	//			//if err = os.Remove("./data/" + file.Name()); err != nil {
	//			//	client.log.Debug(err)
	//			//	client.log.Error("Error while deleting outdated file: " + file.Name())
	//			//}
	//			client.downloadFile("https://ddragon.leagueoflegends.com/cdn/" + client.gameVersion + "/data/en_US/" + file.Name())
	//		}
	//	}
	//}
	return err
}

// readLockFile wait for lockfile to be generated and reads "lockfile", which provides a token to access
// the game client. Returns the content of lockfile as string or err if failed
func (client *DFFClient) readLockFile() (err error) {
	var file *os.File
	// Loop until file is available
	for {
		if file, err = os.Open(path.Join(client.ClientDir, "lockfile")); err == nil {
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

	if req.StatusCode == http.StatusOK {
		client.Log.Debug("Rune page set")
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

// retrieveItems sets an item page
func (client *DFFClient) retrieveItems(data *datatype.OPGGChampData, cachedData *cache.CachedData, champId int, gameType string) (isSet bool) {
	skillBuildStr := "Skill Tree: " +
		data.SkillMasteries[0].Ids[0] + " -> " +
		data.SkillMasteries[0].Ids[1] + " -> " +
		data.SkillMasteries[0].Ids[2]

	// First three skill tree
	firstThreeStr := "First 3 skills: " +
		data.SkillMasteries[0].Builds[0].Order[0] + " -> " +
		data.SkillMasteries[0].Builds[0].Order[1] + " -> " +
		data.SkillMasteries[0].Builds[0].Order[2]

	// Create blocks: "Starter", "Core", "Boots", "Other items" (4 blocks)
	blockList := make([]datatype.ItemBlock, 4)
	blockIdx := 0

	otherItemSet := make(map[int]bool)
	willBeAdded := 0

	// ---- Create Starter Items block
	if len(data.StarterItems) > 0 {
		title := "Starter Items (" + firstThreeStr + ")"
		// +1 to add a ward
		itemList := make([]datatype.Item, len(data.StarterItems[0].Ids)+1)
		for i, id := range data.StarterItems[0].Ids {
			itemList[i] = datatype.Item{
				Count: 1,
				ID:    strconv.Itoa(id),
			}
			otherItemSet[id] = false
		}
		// Add ward to starting item
		itemList[len(itemList)-1] = datatype.Item{
			Count: 1,
			ID:    "3340",
		}
		otherItemSet[3340] = false
		newItemBlock := datatype.ItemBlock{
			HideIfSummonerSpell: "",
			Items:               itemList,
			ShowIfSummonerSpell: "",
			Type:                title,
		}
		blockList[blockIdx] = newItemBlock
		blockIdx++
	}

	// ---- Create Core Items block
	if len(data.CoreItems) > 0 {
		title := "Core Items (" + skillBuildStr + ")"
		itemList := make([]datatype.Item, len(data.CoreItems[0].Ids))
		// Search up to max 5 core item blocks
		for j := 0; j < min(len(data.CoreItems), 5); j++ {
			for i, id := range data.CoreItems[j].Ids {
				// If added to "Core Items" tab, don't add it to "Other Core Items" tab
				if j == 0 {
					itemList[i] = datatype.Item{
						Count: 1,
						ID:    strconv.Itoa(id),
					}
					otherItemSet[id] = false
				} else if _, exist := otherItemSet[id]; !exist {
					otherItemSet[id] = true
					willBeAdded++
				}
			}
		}
		newItemBlock := datatype.ItemBlock{
			HideIfSummonerSpell: "",
			Items:               itemList,
			ShowIfSummonerSpell: "",
			Type:                title,
		}
		blockList[blockIdx] = newItemBlock
		blockIdx++
	}

	// Search up to 10 last item blocks
	for j := 0; j < min(len(data.LastItems), 10); j++ {
		for _, id := range data.LastItems[j].Ids {
			if _, exist := otherItemSet[id]; !exist {
				otherItemSet[id] = true
				willBeAdded++
			}
		}
	}

	// ---- Create Boots block
	if len(data.Boots) > 0 {
		title := "Boots"
		// Add 3 Boots
		itemList := make([]datatype.Item, 3)
		for j := 0; j < min(len(data.Boots), 3); j++ {
			itemList[j] = datatype.Item{
				Count: 1,
				ID:    strconv.Itoa(data.Boots[j].Ids[0]),
			}
		}
		newItemBlock := datatype.ItemBlock{
			HideIfSummonerSpell: "",
			Items:               itemList,
			ShowIfSummonerSpell: "",
			Type:                title,
		}
		blockList[blockIdx] = newItemBlock
		blockIdx++
	}

	// ---- Create other core Items block
	itemList := make([]datatype.Item, willBeAdded)
	idx := 0
	for key, val := range otherItemSet {
		if val {
			newItem := datatype.Item{
				Count: 1,
				ID:    strconv.Itoa(key),
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
	blockList[blockIdx] = newItemBlock
	blockIdx++

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
func (client *DFFClient) retrieveSpells(data *datatype.OPGGChampData, cachedData *cache.CachedData) (isSet bool) {
	if len(data.SummonerSpells) < 1 {
		return false
	}

	spellKeyList := data.SummonerSpells[0].Ids

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

	cachedData.Spells.Spell1ID = int64(spellKeyList[0])
	cachedData.Spells.Spell2ID = int64(spellKeyList[1])

	return true
}

// retrieveRunes will parse runes and make a RuneNamePage structure
func (client *DFFClient) retrieveRunes(data *datatype.OPGGChampData, cachedData *cache.CachedData, champName string, gameType string) (isSet bool) {
	// Create 4 or less pages
	cachedData.RunePages = make([]datatype.DFFRunePage, min(len(data.RunePages), 4))

	// Getting Pick rate/Win rate/Sample count
	for i := 0; i < len(cachedData.RunePages); i++ {
		cachedData.RunePages[i].PickRate = data.RunePages[i].PickRate * 100
		cachedData.RunePages[i].WinRate = float64(data.RunePages[i].Win) / float64(data.RunePages[i].Play) * 100
		cachedData.RunePages[i].SampleCnt = data.RunePages[i].Play
	}

	// Creating rune page name
	for i := 0; i < len(cachedData.RunePages); i++ {
		cachedData.RunePages[i].Name = champName + " (" + strconv.Itoa(i+1) + ")"
	}

	// Creating rune page
	for i := 0; i < len(cachedData.RunePages); i++ {
		runeList := make([]int, 9)

		if len(runeList) != 9 {
			client.Log.Error("Runes updated? Please submit a new issue at " + IssueUrl)
			return false
		}

		idx := 0
		currPage := data.RunePages[i].Builds[0]

		for _, id := range currPage.PrimaryRuneIds {
			runeList[idx] = id
			idx++
		}

		for _, id := range currPage.SecondaryRuneIds {
			runeList[idx] = id
			idx++
		}

		for _, id := range currPage.StatModIds {
			runeList[idx] = id
			idx++
		}

		cachedData.RunePages[i].Page = datatype.RunePage{
			AutoModifiedSelections: []interface{}{},
			Current:                true,
			ID:                     0,
			IsActive:               true,
			IsDeletable:            true,
			IsEditable:             true,
			IsValid:                true,
			LastModified:           0,
			Name:                   ProjectName + " " + cachedData.RunePages[i].Name + " " + gameType,
			Order:                  0,
			PrimaryStyleID:         data.RunePages[i].PrimaryPageID,
			SelectedPerkIds:        runeList,
			SubStyleID:             data.RunePages[i].SecondaryPageID,
		}
	}

	return true
}

func (client *DFFClient) retrieveData(gameMode datatype.GameMode, champion *datatype.Champion, champLabel *widget.Label, position cache.Position) (cacheData *cache.CachedData, pos cache.Position, ok bool) {
	var gameType string

	champLabel.SetText(champion.Alias)
	client.Log.Debug("Selected Champion: ", champion.Alias)

	// Normal mode, no specified position
	if gameMode == datatype.Default && position == cache.None {
		champMeta := client.metaInfo.Existing[champion.ID]
		if champMeta == nil {
			client.Log.Error("champMeta returns nil")
			return nil, 0, false
		}
		// "RIP" champions use Top as a default position
		if champMeta.IsRip || len(champMeta.Positions) == 0 {
			client.Log.Info(champion.Alias, " does not have enough sample count.")
			position = cache.Top
		} else {
			// Other champions use most frequently used position as a default position
			position = champMeta.Positions[0].Position
		}
	}

	cacheData, isCached := client.cache.GetPut(champion.ID, gameMode, position)
	client.Log.Debug("Using cache: ", isCached)
	if !isCached {
		switch gameMode {
		case datatype.Aram:
			gameType = "ARAM"
			client.Log.Info("ARAM MODE IS ON!!!")
			cacheData.URL = "https://na.op.gg/modes/aram/" + champion.Alias + "/build"
		case datatype.Urf:
			gameType = "URF"
			client.Log.Info("ULTRA RAPID FIRE MODE IS ON!!!")
			cacheData.URL = "https://na.op.gg/modes/urf/" + champion.Alias + "/build"
		case datatype.Default:
			cacheData.URL = "https://op.gg/champions/" + champion.Alias + "/" + position.String() + "/build"
		}

		cacheData.CreationTime = time.Now()
		data, ok := client.getFromJson(cacheData.URL)
		if !ok {
			client.Log.Debug("error while getting data from ", cacheData.URL)
			return nil, cache.None, false
		}

		champData := data.Props.PageProps.Data

		isSet := client.retrieveRunes(&champData, cacheData, champion.Alias, gameType)
		if !isSet {
			client.Log.Error("Error while retrieving rune page")
			return nil, cache.None, false
		}

		isSet = client.retrieveItems(&champData, cacheData, champion.ID, gameType)
		if !isSet {
			client.Log.Error("Error while retrieving item page")
			return nil, cache.None, false
		}

		isSet = client.retrieveSpells(&champData, cacheData)
		if !isSet {
			client.Log.Error("Error while retrieving spell page")
			return nil, cache.None, false
		}
	}

	if client.EnableRune {
		if ok, err := client.setRunePage(&cacheData.RunePages[0].Page); !ok || err != nil {
			client.Log.Debug(err)
			client.Log.Error("Unable to set a rune page")
			return nil, cache.None, false
		}
	}

	if client.EnableItem {
		command := "/lol-item-sets/v1/item-sets/" + strconv.Itoa(client.account.SummonerID) + "/sets"

		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(cacheData.ItemPages)
		if err != nil {
			client.Log.Debug(err)
			client.Log.Error("Error while setting items")
			return nil, cache.None, false
		}

		req := client.requestApi("PUT", command, b)
		if req == nil || req.StatusCode != http.StatusCreated {
			client.Log.Debug(err)
			client.Log.Error("Error while setting items")
			return nil, cache.None, false
		}
		client.Log.Debug("Item page set")
	}

	if client.EnableSpell {
		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(cacheData.Spells)
		if err != nil {
			client.Log.Debug(err)
			client.Log.Error("Error while setting spells")
			return nil, cache.None, false
		}

		command := "/lol-champ-select/v1/session/my-selection"
		req := client.requestApi("PATCH", command, b)
		if req == nil || req.StatusCode != http.StatusNoContent {
			client.Log.Debug(err)
			client.Log.Error("Error while setting spells")
			return nil, cache.None, false
		}
		client.Log.Debug("Spells set")
	}

	return cacheData, position, true
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
	if queueId == int(datatype.Aram) || queueId == int(datatype.Urf) {
		gameMode = datatype.GameMode(queueId)
	} else {
		gameMode = datatype.Default
	}

	var cachedData *cache.CachedData
	var ok bool

	lastRole := cache.None
	position := cache.None
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
				position = cache.None
				positionIdx = 0
			}

			// Convert champ id to datatype.Champion
			var champion datatype.Champion
			command := "/lol-champions/v1/inventories/" + strconv.Itoa(client.account.SummonerID) + "/champions/" + strconv.Itoa(champId)
			err = json.NewDecoder(client.requestApi("GET", command, nil).Body).Decode(&champion)
			if err != nil {
				client.Log.Debug(err)
				client.Log.Error("Error while getting runes")
				status.SetText("Error. Check log")
				if client.window != nil {
					client.window.RequestFocus()
				}
			}

			status.SetText("Setting...")
			cachedData, position, ok = client.retrieveData(gameMode, &champion, champLabel, position)
			if !ok {
				status.SetText("Error. Check log")
				if client.window != nil {
					client.window.RequestFocus()
				}
			} else {
				status.SetText("Updated...")
			}
			lastRole = position

			if ok && len(cachedData.RunePages) > 0 {
				runeSelect.Options = make([]string, len(cachedData.RunePages))
				for x, elem := range cachedData.RunePages {
					runeSelect.Options[x] = fmt.Sprintf("%d. PR:%.1f%% WR:%.1f%% Sample: %d", x+1, elem.PickRate, elem.WinRate, elem.SampleCnt)
				}
				runeSelect.Selected = runeSelect.Options[0]
				runeSelect.OnChanged = func(s string) {
					client.Log.Debug("Alternative rune selected")
					endI := strings.Index(s, ". ")
					i, _ := strconv.Atoi(s[:endI])
					ok, err := client.setRunePage(&cachedData.RunePages[i-1].Page)
					if !ok || err != nil {
						status.SetText("Error. Check log")
						if client.window != nil {
							client.window.RequestFocus()
						}
					}

				}
				runeSelect.Refresh()
			}

			if gameMode == datatype.Default {
				p.Options = make([]string, len(client.metaInfo.Existing[champion.ID].Positions))
				for i := 0; i < len(client.metaInfo.Existing[champion.ID].Positions); i++ {
					p.Options[i] += client.metaInfo.Existing[champion.ID].Positions[i].Position.String() + " - " +
						client.metaInfo.Existing[champion.ID].Positions[i].RoleRate
				}

				p.Selected = p.Options[positionIdx]
				p.OnChanged = func(s string) {
					var res string
					for positionIdx, res = range p.Options {
						if s == res {
							position = client.metaInfo.Existing[champion.ID].Positions[positionIdx].Position
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

	var isInGame bool
	for {
		if isInGame, err = client.checkIsInGame(); err != nil {
			status.SetText("Error. Check log")
		}
		if !isInGame {
			break
		}
		time.Sleep(30 * time.Second)
		client.Log.Debug("In game: ", isInGame)
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
