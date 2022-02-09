package core

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/jaeha-choi/DFF/internal/cache"
	"os"
	"time"
)

type Meta struct {
	CreationTime      time.Time
	CacheVersion      uint16
	GameClientVersion string // Must be updated once game client API is accessible
	Existing          map[int]*MetaChampion
}

type MetaChampion struct {
	IsRip     bool
	Positions []MetaPosition
}

type MetaPosition struct {
	Position cache.Position
	RoleRate string
}

// ChampListDataVersion is used to keep track of cache file versions.
// If cache structure is edited in any way, this value must be incremented.
const ChampListDataVersion uint16 = 1

// ChampListDataExpiration data expiration time in days
const ChampListDataExpiration = 7

var incompatibleDataError = errors.New("existing data is incompatible")
var expiredDataError = errors.New("existing data expired")

func (client *DFFClient) saveChampionList(filename string) (err error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(&client.metaInfo)
}

func (client *DFFClient) restoreChampionList(filename string, gameVer string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	if err = gob.NewDecoder(file).Decode(&client.metaInfo); err != nil {
		client.Log.Debug(err)
		return
	}

	if t := time.Now().Sub(client.metaInfo.CreationTime); t >= time.Hour*24*ChampListDataExpiration {
		return expiredDataError
	}

	if client.metaInfo.CacheVersion != ChampListDataVersion || client.metaInfo.GameClientVersion != gameVer {
		return incompatibleDataError
	}

	return
}

func (client *DFFClient) createChampionList(gameVer string) (ok bool) {
	url := "https://na.op.gg/champions"

	data, ok := client.getFromJson(url)
	if !ok {
		client.Log.Fatal("Could not create champion list")
		return
	}

	champList := data.Props.PageProps.ChampionMetaList

	client.metaInfo = &Meta{
		CreationTime:      time.Now(),
		CacheVersion:      ChampListDataVersion,
		GameClientVersion: gameVer,
		Existing:          make(map[int]*MetaChampion, len(champList)),
	}

	var role cache.Position

	for _, champ := range champList {
		client.metaInfo.Existing[champ.ID] = &MetaChampion{}
		client.metaInfo.Existing[champ.ID].IsRip = champ.IsRip
		if champ.IsRip || len(champ.Positions) == 0 {
			client.metaInfo.Existing[champ.ID].Positions = make([]MetaPosition, len(cache.PositionList))
			for i, pos := range cache.PositionList {
				client.metaInfo.Existing[champ.ID].Positions[i].Position = pos
				client.metaInfo.Existing[champ.ID].Positions[i].RoleRate = "Not enough sample count"
			}
		} else {
			client.metaInfo.Existing[champ.ID].Positions = make([]MetaPosition, len(champ.Positions))
			for i, position := range champ.Positions {
				switch position.Name {
				case "TOP":
					role = cache.Top
				case "JUNGLE":
					role = cache.Jungle
				case "MID":
					role = cache.Mid
				case "ADC":
					role = cache.Adc
				case "SUPPORT":
					role = cache.Support
				default:
					client.Log.Debug("Role not found: ", position.Name)
					client.Log.Error("Role changed? Please submit a new issue at " + IssueUrl)
					return false
				}
				client.metaInfo.Existing[champ.ID].Positions[i].Position = role
				client.metaInfo.Existing[champ.ID].Positions[i].RoleRate = fmt.Sprintf("Pick rate: %.1f%%", position.Stats.RoleRate*100)
			}
		}
	}

	return true
}
