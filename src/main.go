package main

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

//go:embed assets/logo.png
var logo []byte

const version = "0.0.0"

func main() {
	a := app.New()
	icon := fyne.NewStaticResource("logo.png", logo)
	a.SetIcon(icon)
	window := a.NewWindow("Adobe Fonts Dump")
	window.SetIcon(icon)
	window.Resize(fyne.NewSize(560, 460))

	image := canvas.NewImageFromResource(icon)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(560, 160))
	title := widget.NewLabel("adobe-fonts-dump")
	title.Alignment = fyne.TextAlignCenter
	versionLabel := widget.NewLabel(fmt.Sprintf("v%s", version))
	versionLabel.Alignment = fyne.TextAlignCenter
	footer := widget.NewLabel("made by ch4og")
	footer.Alignment = fyne.TextAlignCenter

	sourcePath := widget.NewEntry()
	source, warning := rootPath()
	sourcePath.SetText(source)
	sourceBrowse := widget.NewButton("Choose...", nil)
	outputPath := widget.NewEntry()
	outputPath.SetText("dump")
	outputBrowse := widget.NewButton("Choose...", nil)
	skip := widget.NewCheck("Skip existing files", nil)
	statusText := "Ready"
	if warning != "" {
		statusText = warning
	}
	status := widget.NewLabel(statusText)
	dumpButton := widget.NewButton("Dump fonts", nil)

	chooseFolder := func(entry *widget.Entry) {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if uri != nil {
				entry.SetText(uriPath(uri.Path()))
			}
		}, window).Show()
	}
	sourceBrowse.OnTapped = func() { chooseFolder(sourcePath) }
	outputBrowse.OnTapped = func() { chooseFolder(outputPath) }

	dumpButton.OnTapped = func() {
		root := sourcePath.Text
		if root == "" {
			status.SetText("Choose a source folder")
			return
		}
		output := outputPath.Text
		if output == "" {
			output = "dump"
		}
		output, err := filepath.Abs(output)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		sourceBrowse.Disable()
		outputBrowse.Disable()
		dumpButton.Disable()
		sourcePath.Disable()
		outputPath.Disable()
		skip.Disable()
		status.SetText("Dumping fonts...")

		go func() {
			fonts, err := restore(root, output, skip.Checked, func(message string) {
				fyne.Do(func() { status.SetText(message) })
			})
			if err != nil {
				fyne.Do(func() {
					status.SetText("Dump failed")
					sourceBrowse.Enable()
					outputBrowse.Enable()
					dumpButton.Enable()
					sourcePath.Enable()
					outputPath.Enable()
					skip.Enable()
					dialog.ShowError(err, window)
				})
				return
			}
			fyne.Do(func() {
				status.SetText(fmt.Sprintf("Restored %d font(s); skipped %d.", fonts.Restored, fonts.Skipped))
				dialog.ShowCustomConfirm(
					"Success",
					"Open output dir",
					"OK",
					widget.NewLabel(fmt.Sprintf("Success.\nRestored %d fonts, skipped %d.", fonts.Restored, fonts.Skipped)),
					func(open bool) {
						if open {
							if err := openDirectory(output); err != nil {
								dialog.ShowError(err, window)
							}
						}
					},
					window,
				)
				sourceBrowse.Enable()
				outputBrowse.Enable()
				dumpButton.Enable()
				sourcePath.Enable()
				outputPath.Enable()
				skip.Enable()
			})
		}()
	}

	window.SetContent(container.NewVBox(
		image,
		title,
		widget.NewLabel("Source folder"),
		container.NewBorder(nil, nil, nil, sourceBrowse, sourcePath),
		widget.NewLabel("Output folder"),
		container.NewBorder(nil, nil, nil, outputBrowse, outputPath),
		skip,
		status,
		dumpButton,
		versionLabel,
		footer,
	))
	window.Show()
	if warning != "" {
		dialog.ShowInformation("Adobe Fonts path", warning, window)
	}
	a.Run()
}
