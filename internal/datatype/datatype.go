package datatype

import "time"

type GameMode int

const (
	Default GameMode = 0
	Aram    GameMode = 450
	Urf     GameMode = 900
)

type Spells struct {
	//SelectedSkinID int32 `json:"selectedSkinId"`
	Spell1ID int64 `json:"spell1Id"`
	Spell2ID int64 `json:"spell2Id"`
	//WardSkinID     int64 `json:"wardSkinId"`
}

type DFFRunePage struct {
	Name      string
	PickRate  float64
	WinRate   float64
	SampleCnt int
	Page      RunePage
}

type RunePageCount struct {
	OwnedPageCount int `json:"ownedPageCount"`
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

type OPGGChampList struct {
	ID           int  `json:"id"`
	IsRotation   bool `json:"is_rotation"`
	IsRip        bool `json:"is_rip"`
	AverageStats struct {
		WinRate  float64     `json:"win_rate"`
		PickRate float64     `json:"pick_rate"`
		BanRate  float64     `json:"ban_rate"`
		Kda      interface{} `json:"kda"`
		Tier     int         `json:"tier"`
		Rank     int         `json:"rank"`
	} `json:"average_stats"`
	Positions []struct {
		Name  string `json:"name"`
		Stats struct {
			WinRate  float64 `json:"win_rate"`
			PickRate float64 `json:"pick_rate"`
			BanRate  float64 `json:"ban_rate"`
			RoleRate float64 `json:"role_rate"`
			TierData struct {
				Tier     int `json:"tier"`
				Rank     int `json:"rank"`
				RankDiff int `json:"rank_diff"`
			} `json:"tier_data"`
		} `json:"stats"`
	} `json:"positions"`
	ImageURL         string  `json:"image_url"`
	Name             string  `json:"name"`
	Display          bool    `json:"display"`
	Key              string  `json:"key"`
	PositionWinRate  float64 `json:"positionWinRate,omitempty"`
	PositionPickRate float64 `json:"positionPickRate,omitempty"`
	PositionBanRate  float64 `json:"positionBanRate,omitempty"`
	PositionRoleRate float64 `json:"positionRoleRate,omitempty"`
	PositionTierData struct {
		Tier     int `json:"tier"`
		Rank     int `json:"rank"`
		RankDiff int `json:"rank_diff"`
	} `json:"positionTierData,omitempty"`
	PositionTier int `json:"positionTier,omitempty"`
	PositionRank int `json:"positionRank,omitempty"`
}

type OPGGChampData struct {
	Summary struct {
		Version struct {
			Version    string `json:"version"`
			PatchIndex int    `json:"patch_index"`
		} `json:"version"`
		Meta struct {
			ID        int      `json:"id"`
			Key       string   `json:"key"`
			Name      string   `json:"name"`
			ImageURL  string   `json:"image_url"`
			EnemyTips []string `json:"enemy_tips"`
			AllyTips  []string `json:"ally_tips"`
			Skins     []struct {
				Name         string `json:"name"`
				HasChromas   bool   `json:"has_chromas"`
				SplashImage  string `json:"splash_image"`
				LoadingImage string `json:"loading_image"`
				TilesImage   string `json:"tiles_image"`
			} `json:"skins"`
			Passive struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				ImageURL    string `json:"image_url"`
				VideoURL    string `json:"video_url"`
			} `json:"passive"`
			Spells []struct {
				Key          string `json:"key"`
				Name         string `json:"name"`
				Description  string `json:"description"`
				MaxRank      int    `json:"max_rank"`
				RangeBurn    []int  `json:"range_burn"`
				CooldownBurn []int  `json:"cooldown_burn"`
				CostBurn     []int  `json:"cost_burn"`
				Tooltip      string `json:"tooltip"`
				ImageURL     string `json:"image_url"`
				VideoURL     string `json:"video_url"`
			} `json:"spells"`
		} `json:"meta"`
		Summary struct {
			ID           int  `json:"id"`
			IsRotation   bool `json:"is_rotation"`
			IsRip        bool `json:"is_rip"`
			AverageStats struct {
				WinRate  float64 `json:"win_rate"`
				PickRate float64 `json:"pick_rate"`
				BanRate  float64 `json:"ban_rate"` // null if ARAM/URF
				Kda      float64 `json:"kda"`      // null if norm
				Tier     int     `json:"tier"`
				Rank     int     `json:"rank"`
			} `json:"average_stats"`
			Positions []struct { // null if ARAM/URF
				Name  string `json:"name"`
				Stats struct {
					WinRate  float64 `json:"win_rate"`
					PickRate float64 `json:"pick_rate"`
					BanRate  float64 `json:"ban_rate"`
					RoleRate float64 `json:"role_rate"`
					TierData struct {
						Tier     int `json:"tier"`
						Rank     int `json:"rank"`
						RankDiff int `json:"rank_diff"`
					} `json:"tier_data"`
				} `json:"stats"`
			} `json:"positions"`
		} `json:"summary"`
		Opponents [][]struct { // does not exist if ARAM/URF and some normal
			ChampionID int `json:"champion_id"`
			Play       int `json:"play"`
			Win        int `json:"win"`
			Meta       struct {
				ID        int      `json:"id"`
				Key       string   `json:"key"`
				Name      string   `json:"name"`
				ImageURL  string   `json:"image_url"`
				EnemyTips []string `json:"enemy_tips"`
				AllyTips  []string `json:"ally_tips"`
				Skins     []struct {
					Name         string `json:"name"`
					HasChromas   bool   `json:"has_chromas"`
					SplashImage  string `json:"splash_image"`
					LoadingImage string `json:"loading_image"`
					TilesImage   string `json:"tiles_image"`
				} `json:"skins"`
				Passive struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					ImageURL    string `json:"image_url"`
					VideoURL    string `json:"video_url"`
				} `json:"passive"`
				Spells []struct {
					Key          string `json:"key"`
					Name         string `json:"name"`
					Description  string `json:"description"`
					MaxRank      int    `json:"max_rank"`
					RangeBurn    []int  `json:"range_burn"`
					CooldownBurn []int  `json:"cooldown_burn"`
					CostBurn     []int  `json:"cost_burn"`
					Tooltip      string `json:"tooltip"`
					ImageURL     string `json:"image_url"`
					VideoURL     string `json:"video_url"`
				} `json:"spells"`
			} `json:"meta"`
			WinRate float64 `json:"win_rate"`
		} `json:"opponents"`
	} `json:"summary"`
	//Meta struct {
	//	Runes []struct {
	//		ID           int    `json:"id"`
	//		PageID       int    `json:"page_id"`
	//		SlotSequence int    `json:"slot_sequence"`
	//		RuneSequence int    `json:"rune_sequence"`
	//		Key          string `json:"key"`
	//		Name         string `json:"name"`
	//		ShortDesc    string `json:"short_desc"`
	//		LongDesc     string `json:"long_desc"`
	//		ImageURL     string `json:"image_url"`
	//	} `json:"runes"`
	//	RunePages []struct {
	//		ID          int    `json:"id"`
	//		Name        string `json:"name"`
	//		Description string `json:"description"`
	//		Slogan      string `json:"slogan"`
	//		ImageURL    string `json:"image_url"`
	//	} `json:"runePages"`
	//	StatMods []struct {
	//		ID        int    `json:"id"`
	//		Name      string `json:"name"`
	//		ShortDesc string `json:"short_desc"`
	//		ImageURL  string `json:"image_url"`
	//	} `json:"statMods"`
	//	Items []struct {
	//		ID        int         `json:"id"`
	//		Name      string      `json:"name"`
	//		ImageURL  string      `json:"image_url"`
	//		IsMythic  bool        `json:"is_mythic"`
	//		IntoItems []int       `json:"into_items"`
	//		FromItems interface{} `json:"from_items"`
	//		Gold      struct {
	//			Sell        int  `json:"sell"`
	//			Total       int  `json:"total"`
	//			Base        int  `json:"base"`
	//			Purchasable bool `json:"purchasable"`
	//		} `json:"gold"`
	//		Plaintext   string `json:"plaintext"`
	//		Description string `json:"description"`
	//	} `json:"items"`
	//	Spells []struct {
	//		ID          int    `json:"id"`
	//		Key         string `json:"key"`
	//		Name        string `json:"name"`
	//		Description string `json:"description"`
	//		ImageURL    string `json:"image_url"`
	//	} `json:"spells"`
	//} `json:"meta"`
	SummonerSpells []struct {
		Ids      []int   `json:"ids"`
		Win      int     `json:"win"`
		Play     int     `json:"play"`
		PickRate float64 `json:"pick_rate"`
	} `json:"summoner_spells"`
	Trends struct { // does not exist if ARAM/URF and some normal
		TotalRank         int `json:"total_rank"`
		TotalPositionRank int `json:"total_position_rank"`
		Win               []struct {
			Version   string    `json:"version"`
			Rate      float64   `json:"rate"`
			Average   float64   `json:"average"`
			Rank      int       `json:"rank"`
			CreatedAt time.Time `json:"created_at"`
		} `json:"win"`
		Ban []struct {
			Version   string    `json:"version"`
			Rate      float64   `json:"rate"`
			Average   float64   `json:"average"`
			Rank      int       `json:"rank"`
			CreatedAt time.Time `json:"created_at"`
		} `json:"ban"`
		Pick []struct {
			Version   string    `json:"version"`
			Rate      float64   `json:"rate"`
			Average   float64   `json:"average"`
			Rank      int       `json:"rank"`
			CreatedAt time.Time `json:"created_at"`
		} `json:"pick"`
	} `json:"trends"`
	GameLengths []struct {
		GameLength int     `json:"game_length"`
		Rate       float64 `json:"rate"`
		Average    float64 `json:"average"`
		Rank       int     `json:"rank"`
	} `json:"game_lengths"`
	Skills []struct {
		Order    []string `json:"order"`
		Play     int      `json:"play"`
		Win      int      `json:"win"`
		PickRate float64  `json:"pick_rate"`
	} `json:"skills"`
	SkillMasteries []struct {
		Ids      []string `json:"ids"`
		Play     int      `json:"play"`
		Win      int      `json:"win"`
		PickRate float64  `json:"pick_rate"`
		Builds   []struct {
			Order    []string `json:"order"`
			Play     int      `json:"play"`
			Win      int      `json:"win"`
			PickRate float64  `json:"pick_rate"`
		} `json:"builds"`
	} `json:"skill_masteries"`
	Runes []struct {
		ID               int     `json:"id"`
		PrimaryPageID    int     `json:"primary_page_id"`
		PrimaryRuneIds   []int   `json:"primary_rune_ids"`
		SecondaryPageID  int     `json:"secondary_page_id"`
		SecondaryRuneIds []int   `json:"secondary_rune_ids"`
		StatModIds       []int   `json:"stat_mod_ids"`
		Play             int     `json:"play"`
		Win              int     `json:"win"`
		PickRate         float64 `json:"pick_rate"`
	} `json:"runes"`
	RunePages []struct {
		ID              int     `json:"id"`
		PrimaryPageID   int     `json:"primary_page_id"`
		SecondaryPageID int     `json:"secondary_page_id"`
		Play            int     `json:"play"`
		PickRate        float64 `json:"pick_rate"`
		Win             int     `json:"win"`
		Builds          []struct {
			ID               int     `json:"id"`
			PrimaryPageID    int     `json:"primary_page_id"`
			PrimaryRuneIds   []int   `json:"primary_rune_ids"`
			SecondaryPageID  int     `json:"secondary_page_id"`
			SecondaryRuneIds []int   `json:"secondary_rune_ids"`
			StatModIds       []int   `json:"stat_mod_ids"`
			Play             int     `json:"play"`
			Win              int     `json:"win"`
			PickRate         float64 `json:"pick_rate"`
		} `json:"builds"`
	} `json:"rune_pages"`
	CoreItems []struct {
		Ids      []int   `json:"ids"`
		Win      int     `json:"win"`
		Play     int     `json:"play"`
		PickRate float64 `json:"pick_rate"`
	} `json:"core_items"`
	Boots []struct {
		Ids      []int   `json:"ids"`
		Win      int     `json:"win"`
		Play     int     `json:"play"`
		PickRate float64 `json:"pick_rate"`
	} `json:"boots"`
	StarterItems []struct {
		Ids      []int   `json:"ids"`
		Win      int     `json:"win"`
		Play     int     `json:"play"`
		PickRate float64 `json:"pick_rate"`
	} `json:"starter_items"`
	LastItems []struct {
		Ids      []int   `json:"ids"`
		Win      int     `json:"win"`
		Play     int     `json:"play"`
		PickRate float64 `json:"pick_rate"`
	} `json:"last_items"`
}

type OPGGResponse struct {
	Props struct {
		PageProps struct {
			ChampionMetaList []OPGGChampList `json:"championMetaList"`
			Data             OPGGChampData   `json:"data"`
		} `json:"pageProps"`
	} `json:"props"`
}
