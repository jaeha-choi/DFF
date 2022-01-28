package datatype

type Spells struct {
	//SelectedSkinID int32 `json:"selectedSkinId"`
	Spell1ID int64 `json:"spell1Id"`
	Spell2ID int64 `json:"spell2Id"`
	//WardSkinID     int64 `json:"wardSkinId"`
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
