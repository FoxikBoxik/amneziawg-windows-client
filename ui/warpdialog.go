/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2019-2022 WireGuard LLC. All Rights Reserved.
 */

package ui

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"

	"github.com/amnezia-vpn/amneziawg-windows-client/l18n"
)

func RunWarpDialog(owner walk.Form) error {
	var disposables walk.Disposables
	defer disposables.Treat()

	dlg, err := walk.NewDialogWithFixedSize(owner)
	if err != nil {
		return err
	}
	disposables.Add(dlg)
	dlg.SetTitle("WARP для AmneziaWG")
	if icon, err := loadLogoIcon(32); err == nil {
		dlg.SetIcon(icon)
	}
	// Modern size
	dlg.SetMinMaxSize(walk.Size{480, 340}, walk.Size{480, 340})

	mainLayout := walk.NewVBoxLayout()
	mainLayout.SetMargins(walk.Margins{22, 10, 22, 10})
	mainLayout.SetSpacing(7)
	dlg.SetLayout(mainLayout)

	// Title
	title, err := walk.NewTextLabel(dlg)
	if err != nil {
		return err
	}
	titleFont, _ := walk.NewFont("Segoe UI", 15, walk.FontBold)
	title.SetFont(titleFont)
	title.SetTextAlignment(walk.AlignHCenterVNear)
	title.SetText(l18n.Sprintf("WARP для AmneziaWG"))

	// Subtitle
	sub, err := walk.NewTextLabel(dlg)
	if err != nil {
		return err
	}
	subFont, _ := walk.NewFont("Segoe UI", 9, walk.FontItalic)
	sub.SetFont(subFont)
	sub.SetTextAlignment(walk.AlignHCenterVNear)
	sub.SetText(l18n.Sprintf("Инструкция"))

	// Steps container
	stepsContainer, err := walk.NewComposite(dlg)
	if err != nil {
		return err
	}
	stepsLayout := walk.NewVBoxLayout()
	stepsLayout.SetMargins(walk.Margins{0, 10, 0, 0})
	stepsLayout.SetSpacing(5)
	stepsContainer.SetLayout(stepsLayout)

	stepFont, _ := walk.NewFont("Segoe UI", 10, walk.FontBold)

	step1, _ := walk.NewTextLabel(stepsContainer)
	step1.SetFont(stepFont)
	step1.SetText(l18n.Sprintf("1. Откройте сайт"))

	step2, _ := walk.NewTextLabel(stepsContainer)
	step2.SetFont(stepFont)
	step2.SetText(l18n.Sprintf("2. Сгенерируйте конфигурацию."))

	step21, _ := walk.NewTextLabel(stepsContainer)
	step21Font, _ := walk.NewFont("Segoe UI", 8, walk.FontBold)
	step21.SetFont(step21Font)
	step21.SetText(l18n.Sprintf("2.1. Выбирайте AmneziaWG - AWG 1.5‘’"))

	step3, _ := walk.NewTextLabel(stepsContainer)
	step3.SetFont(stepFont)
	step3.SetText(l18n.Sprintf("3. Импортируйте .conf файл через кнопку 'Импортировать туннели из файла'"))

	step4, _ := walk.NewTextLabel(stepsContainer)
	step4.SetFont(stepFont)
	step4.SetText(l18n.Sprintf("4. Подключите туннель и проверьте подключение"))
	

	open2ipBtn, _ := walk.NewPushButton(stepsContainer)
	open2ipBtn.SetText("Открыть 2ip.io")
	open2ipBtn.Clicked().Attach(func() {
		win.ShellExecute(dlg.Handle(), windows.StringToUTF16Ptr("open"), windows.StringToUTF16Ptr("https://2ip.io/"), nil, nil, win.SW_SHOWNORMAL)
	})

	step5, _ := walk.NewTextLabel(stepsContainer)
	step5.SetFont(stepFont)
	step5.SetText(l18n.Sprintf("5. Если написало что у вас пройвайдер Cloudflare, то все успешно выполненно!"))
	step51, _ := walk.NewTextLabel(stepsContainer)
	step51Font, _ := walk.NewFont("Segoe UI", 8, walk.FontBold)
	step51.SetFont(step51Font)
	step51.SetText(l18n.Sprintf("5.1. Если провайдер не Cloudflare или нет подключения попробуйте обратиться в чат ниже."))
	// Tips
	tips, _ := walk.NewTextLabel(dlg)
	tipsFont, _ := walk.NewFont("Segoe UI", 9, walk.FontItalic)
	tips.SetFont(tipsFont)
	tips.SetTextAlignment(walk.AlignHNearVNear)
	tips.SetText(l18n.Sprintf("Советы: не делитесь ключами\nвозникли проблемы — получите новый ключ или напишите в чат ниже."))

	// Info panel at bottom
	info, err := walk.NewTextLabel(dlg)
	if err != nil {
		return err
	}
	infoFont, _ := walk.NewFont("Segoe UI", 8, 0)
	info.SetFont(infoFont)
	info.SetTextAlignment(walk.AlignHNearVNear)
	info.SetText(l18n.Sprintf("made by foxinabox"))

	// Buttons
	btns, err := walk.NewComposite(dlg)
	if err != nil {
		return err
	}
	hl := walk.NewHBoxLayout()
	hl.SetMargins(walk.Margins{VNear: 8})
	btns.SetLayout(hl)
	walk.NewHSpacer(btns)

	openBtn, err := walk.NewPushButton(btns)
	if err != nil {
		return err
	}
	openBtnFont, _ := walk.NewFont("Segoe UI", 10, walk.FontBold)
	openBtn.SetFont(openBtnFont)
	openBtn.SetText(l18n.Sprintf("Открыть сайт"))
	openBtn.Clicked().Attach(func() {
		win.ShellExecute(dlg.Handle(), nil, windows.StringToUTF16Ptr("https://warp-generator.github.io/warp/"), nil, nil, win.SW_SHOWNORMAL)
	})

	closeBtn, err := walk.NewPushButton(btns)
	if err != nil {
		return err
	}
	closeBtn.SetText(l18n.Sprintf("Закрыть"))
	closeBtn.Clicked().Attach(dlg.Accept)
	walk.NewHSpacer(btns)

	dlg.SetDefaultButton(closeBtn)
	dlg.SetCancelButton(closeBtn)

	disposables.Spare()
	dlg.Run()
	return nil
}
