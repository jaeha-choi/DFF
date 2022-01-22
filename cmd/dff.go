package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jaeha-choi/DFF/internal/core"
	"github.com/jaeha-choi/DFF/internal/updater"
	"github.com/jaeha-choi/DFF/pkg/log"
	"os"
)

func main() {
	var logOut *os.File

	logOut, err := os.OpenFile("dff.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logOut = os.Stdout
	}

	client := core.Initialize(logOut)

	// Does not seem to work because of w.ShowAndRun()?
	// Always update configuration file upon exit
	//defer func() {
	//	logOut.Close()
	//	_ = client.WriteConfig()
	//}()

	a := app.New()
	w := a.NewWindow(core.ProjectName + " " + core.Version)

	w.SetOnClosed(func() {
		client.Log.Debug("Saving configuration")
		if err = client.WriteConfig(); err != nil {
			client.Log.Debug(err)
			client.Log.Errorf("Error while writing configuration")
		}
		client.Log.Debug("Exiting...")
		logOut.Close()
		os.Exit(0)
	})

	sl := widget.NewSlider(1, 5)
	sl.Value = client.Interval
	sl.Step = 0.5
	sl.OnChanged = func(f float64) {
		client.Interval = f
	}

	infoTextStyle := fyne.TextStyle{
		Bold:      true,
		Italic:    false,
		Monospace: true,
	}
	status := widget.NewLabelWithStyle("Not running", fyne.TextAlignCenter, infoTextStyle)
	selectedChamp := widget.NewLabelWithStyle("Not selected", fyne.TextAlignCenter, infoTextStyle)

	enableRunesCheck := widget.NewCheck("", func(b bool) {
		client.EnableRune = b
	})
	enableRunesCheck.SetChecked(client.EnableRune)

	enableItemsCheck := widget.NewCheck("", func(b bool) {
		client.EnableItem = b
	})
	enableItemsCheck.SetChecked(client.EnableItem)

	roleSelect := widget.NewSelect(nil, nil)
	roleSelect.PlaceHolder = "No champion selected"

	runeSelect := widget.NewSelect(nil, nil)
	runeSelect.PlaceHolder = "No rune selected"

	enableSpellCheck := widget.NewCheck("", func(b bool) {
		client.EnableSpell = b
	})
	enableSpellCheck.SetChecked(client.EnableSpell)

	enableDFlash := widget.NewCheck("", func(b bool) {
		client.DFlash = b
	})
	enableDFlash.SetChecked(client.DFlash)

	enableDebugging := widget.NewCheck("", func(b bool) {
		client.Debug = b
		client.Log.Debug("Debug mode updated: ", b)
		if b {
			client.Log.Mode = log.DEBUG
		} else {
			client.Log.Mode = log.INFO
		}
	})
	enableDebugging.SetChecked(client.Debug)

	checkUpdateButton := widget.NewButton("Check Update", func() {
		updater.Update(client.Log, w)
	})

	go func() {
		for {
			client.Run(w, status, roleSelect, selectedChamp, runeSelect)
		}
	}()

	///lol-lobby/v1/lobby/availability
	///lol-lobby/v1/lobby/countdown
	///riotclient/get_region_locale

	w.SetContent(
		container.NewVBox(
			container.NewHBox(
				container.NewVBox(
					widget.NewLabel(core.ProjectName+" "+core.Version),
					widget.NewLabel("Program Status:"),
					status,
					container.NewHBox(widget.NewLabel("Current Champion:")),
					selectedChamp),
				container.NewVBox(
					checkUpdateButton,
					container.NewHBox(widget.NewLabel("Debug"), enableDebugging),
					container.NewHBox(widget.NewLabel("Auto runes"), enableRunesCheck),
					container.NewHBox(widget.NewLabel("Auto items"), enableItemsCheck),
					container.NewHBox(widget.NewLabel("Auto spells"), enableSpellCheck),
					container.NewHBox(widget.NewLabel("Left Flash"), enableDFlash),
					widget.NewLabel("Polling interval"),
					sl,
				)),
			roleSelect,
			runeSelect))

	w.SetFixedSize(true)
	w.ShowAndRun()
}
