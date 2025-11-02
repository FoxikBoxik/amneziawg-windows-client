/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2019-2022 WireGuard LLC. All Rights Reserved.
 */

package ui

import (
	"runtime"
	"strings"

	"github.com/lxn/walk"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"

	"github.com/amnezia-vpn/amneziawg-windows-client/l18n"
	"github.com/amnezia-vpn/amneziawg-windows-client/version"
)

var (
	easterEggIndex     = -1
	showingAboutDialog *walk.Dialog
)

func onAbout(owner walk.Form) {
	showError(runAboutDialog(owner), owner)
}

func runAboutDialog(owner walk.Form) error {
	if showingAboutDialog != nil {
		showingAboutDialog.Show()
		raise(showingAboutDialog.Handle())
		return nil
	}

	vbl := walk.NewVBoxLayout()
	vbl.SetMargins(walk.Margins{80, 20, 80, 20})
	vbl.SetSpacing(10)

	var disposables walk.Disposables
	defer disposables.Treat()

	var err error
	showingAboutDialog, err = walk.NewDialogWithFixedSize(owner)
	if err != nil {
		return err
	}
	defer func() {
		showingAboutDialog = nil
	}()
	disposables.Add(showingAboutDialog)
	showingAboutDialog.SetTitle(l18n.Sprintf("О DementiaWG"))
	showingAboutDialog.SetLayout(vbl)
	if icon, err := loadLogoIcon(32); err == nil {
		showingAboutDialog.SetIcon(icon)
	}

	font, _ := walk.NewFont("Segoe UI", 9, 0)
	showingAboutDialog.SetFont(font)

	iv, err := walk.NewImageView(showingAboutDialog)
	if err != nil {
		return err
	}
	iv.SetCursor(walk.CursorHand())
	iv.MouseUp().Attach(func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton {
			win.ShellExecute(showingAboutDialog.Handle(), nil, windows.StringToUTF16Ptr("https://www.youtube.com/watch?v=dQw4w9WgXcQ"), nil, nil, win.SW_SHOWNORMAL)
		} else if easterEggIndex >= 0 && button == walk.RightButton {
			if icon, err := loadSystemIcon("moricons", int32(easterEggIndex), 128); err == nil {
				iv.SetImage(icon)
				easterEggIndex++
			} else {
				easterEggIndex = -1
				if logo, err := loadLogoIcon(128); err == nil {
					iv.SetImage(logo)
				}
			}
		}
	})
	if logo, err := loadLogoIcon(128); err == nil {
		iv.SetImage(logo)
	}
	iv.Accessibility().SetName(l18n.Sprintf("AmneziaWG logo image"))

	wgLbl, err := walk.NewTextLabel(showingAboutDialog)
	if err != nil {
		return err
	}
	wgFont, _ := walk.NewFont("Segoe UI", 16, walk.FontBold)
	wgLbl.SetFont(wgFont)
	wgLbl.SetTextAlignment(walk.AlignHCenterVNear)
	wgLbl.SetText("DemetiaWG")

	detailsLbl, err := walk.NewTextLabel(showingAboutDialog)
	if err != nil {
		return err
	}
	detailsLbl.SetTextAlignment(walk.AlignHCenterVNear)
	detailsLbl.SetText(l18n.Sprintf("App version: %s\nWintun version: %s\nGo version: %s\nOperating system: %s\nArchitecture: %s", version.Number, version.WintunVersion(), strings.TrimPrefix(runtime.Version(), "go"), version.OsName(), version.Arch()))

	copyrightLbl, err := walk.NewTextLabel(showingAboutDialog)
	if err != nil {
		return err
	}
	copyrightFont, _ := walk.NewFont("Segoe UI", 7, 0)
	copyrightLbl.SetFont(copyrightFont)
	copyrightLbl.SetTextAlignment(walk.AlignHCenterVNear)
	copyrightLbl.SetText("Copyright © 2022-2025 AmneziaVPN.\nModification for AmenziaWG with new Protocol\nchanges and modifications made by FoxinaBox. \nAll Rights Not Reserved.")

	// Links grid (2x2)
	linksCP, err := walk.NewComposite(showingAboutDialog)
	if err != nil {
		return err
	}
	gl := walk.NewGridLayout()
	gl.SetSpacing(8)
	gl.SetMargins(walk.Margins{VNear: 10})
	gl.SetColumnStretchFactor(0, 1)
	gl.SetColumnStretchFactor(1, 1)
	linksCP.SetLayout(gl)

	// Action buttons opening URLs in a 2x2 grid
	githubPB, err := walk.NewPushButton(linksCP)
	if err != nil {
		return err
	}
	githubPB.SetAlignment(walk.AlignHCenterVNear)
	githubPB.SetText(l18n.Sprintf("GitHub"))
	githubPB.Clicked().Attach(func() {
		win.ShellExecute(showingAboutDialog.Handle(), nil, windows.StringToUTF16Ptr("https://github.com/FoxikBoxik/amneziawg-windows-client"), nil, nil, win.SW_SHOWNORMAL)
	})
	gl.SetRange(githubPB, walk.Rectangle{0, 0, 1, 1})

	telegramPB, err := walk.NewPushButton(linksCP)
	if err != nil {
		return err
	}
	telegramPB.SetAlignment(walk.AlignHCenterVNear)
	telegramPB.SetText(l18n.Sprintf("Telegram"))
	telegramPB.Clicked().Attach(func() {
		win.ShellExecute(showingAboutDialog.Handle(), nil, windows.StringToUTF16Ptr("https://t.me/findllimonix"), nil, nil, win.SW_SHOWNORMAL)
	})
	gl.SetRange(telegramPB, walk.Rectangle{1, 0, 1, 1})

	supportPB, err := walk.NewPushButton(linksCP)
	if err != nil {
		return err
	}
	supportPB.SetAlignment(walk.AlignHCenterVNear)
	supportPB.SetText(l18n.Sprintf("Поддержать автора"))
	supportPB.Clicked().Attach(func() {
		win.ShellExecute(showingAboutDialog.Handle(), nil, windows.StringToUTF16Ptr("https://www.donationalerts.com/r/foxinabox"), nil, nil, win.SW_SHOWNORMAL)
	})
	gl.SetRange(supportPB, walk.Rectangle{0, 1, 1, 1})

	questionsPB, err := walk.NewPushButton(linksCP)
	if err != nil {
		return err
	}
	questionsPB.SetAlignment(walk.AlignHCenterVNear)
	questionsPB.SetText(l18n.Sprintf("Получить WARP"))
	questionsPB.Clicked().Attach(func() {
		showError(RunWarpDialog(showingAboutDialog), showingAboutDialog)
	})
	gl.SetRange(questionsPB, walk.Rectangle{1, 1, 1, 1})

	// Bottom row with only Close button
	buttonCP, err := walk.NewComposite(showingAboutDialog)
	if err != nil {
		return err
	}
	hbl := walk.NewHBoxLayout()
	hbl.SetMargins(walk.Margins{VNear: 10})
	buttonCP.SetLayout(hbl)
	walk.NewHSpacer(buttonCP)
	closePB, err := walk.NewPushButton(buttonCP)
	if err != nil {
		return err
	}
	closePB.SetAlignment(walk.AlignHCenterVNear)
	closePB.SetText(l18n.Sprintf("Close"))
	closePB.Clicked().Attach(showingAboutDialog.Accept)
	walk.NewHSpacer(buttonCP)

	showingAboutDialog.SetDefaultButton(closePB)
	showingAboutDialog.SetCancelButton(closePB)

	disposables.Spare()

	showingAboutDialog.Run()

	return nil
}

// WARP dialog moved to warpdialog.go (RunWarpDialog)
