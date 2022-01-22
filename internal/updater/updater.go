package updater

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jaeha-choi/DFF/internal/core"
	"github.com/jaeha-choi/DFF/pkg/log"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

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

func Update(log *log.Logger, w fyne.Window) {
	req, err := http.Get("https://api.github.com/repos/jaeha-choi/DFF/releases/latest")
	if err != nil || req.StatusCode != http.StatusOK {
		log.Debug(err)
		log.Errorf("Error while getting latest version info. Status code: %d", req.StatusCode)
		widget.ShowPopUpAtPosition(widget.NewLabel("Could not connect to github repository."),
			w.Canvas(), fyne.NewPos(50, 50))
	} else {
		var update GitHubUpdate
		if err = json.NewDecoder(req.Body).Decode(&update); err != nil {
			log.Debug(err)
			log.Error("Error while decoding version")
			return
		}

		if update.TagName != core.Version {
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

			log.Info(name[len(name)-1] + " downloaded.")
			popup := container.NewVBox(widget.NewLabel("Version "+update.TagName+" downloaded."),
				widget.NewLabel("Program will now exit."))
			widget.ShowPopUpAtPosition(popup,
				w.Canvas(), fyne.NewPos(50, 50))
			time.Sleep(3 * time.Second)
			log.Info(core.ProjectName + " closed for update.")
			os.Exit(0)
		} else {
			widget.ShowPopUpAtPosition(widget.NewLabel("No update found"),
				w.Canvas(), fyne.NewPos(50, 50))
		}
	}
}
