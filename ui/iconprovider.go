/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2019-2022 WireGuard LLC. All Rights Reserved.
 */

package ui

import (
	"github.com/lxn/walk"

	"github.com/amnezia-vpn/amneziawg-windows-client/l18n"
	"github.com/amnezia-vpn/amneziawg-windows-client/manager"
)

type widthAndState struct {
	width int
	state manager.TunnelState
}

type widthAndDllIdx struct {
	width int
	idx   int32
	dll   string
}

type widthAndPath struct {
	width int
	path  string
}

var cachedOverlayIconsForWidthAndState = make(map[widthAndState]walk.Image)
var cachedImagesForWidthAndPath = make(map[widthAndPath]walk.Image)

func iconWithOverlayForState(state manager.TunnelState, size int) (icon walk.Image, err error) {
	icon = cachedOverlayIconsForWidthAndState[widthAndState{size, state}]
	if icon != nil {
		return
	}

	wireguardIcon, err := loadLogoImage(size)
	if err != nil {
		return
	}

	if state == manager.TunnelStopped {
		return wireguardIcon, err
	}

	iconSize := wireguardIcon.Size()
	w := int(float64(iconSize.Width) * 0.65)
	h := int(float64(iconSize.Height) * 0.65)
	overlayBounds := walk.Rectangle{iconSize.Width - w, iconSize.Height - h, w, h}
	overlayIcon, err := imageForState(state, overlayBounds.Width)
	if err != nil {
		return
	}

	icon = walk.NewPaintFuncImage(walk.Size{size, size}, func(canvas *walk.Canvas, bounds walk.Rectangle) error {
		if err := canvas.DrawImageStretched(wireguardIcon, bounds); err != nil {
			return err
		}
		if err := canvas.DrawImageStretched(overlayIcon, overlayBounds); err != nil {
			return err
		}
		return nil
	})

	cachedOverlayIconsForWidthAndState[widthAndState{size, state}] = icon

	return
}

var cachedIconsForWidthAndState = make(map[widthAndState]*walk.Icon)

func iconForState(state manager.TunnelState, size int) (icon *walk.Icon, err error) {
	icon = cachedIconsForWidthAndState[widthAndState{size, state}]
	if icon != nil {
		return
	}
	switch state {
	case manager.TunnelStarted:
		icon, err = loadSystemIcon("imageres", -106, size)
	case manager.TunnelStopped:
		icon, err = walk.NewIconFromResourceIdWithSize(8, walk.Size{size, size})
	default:
		icon, err = loadSystemIcon("shell32", -16739, size)
	}
	if err == nil {
		cachedIconsForWidthAndState[widthAndState{size, state}] = icon
	}
	return
}

// imageForState возвращает изображение для состояния туннеля
func imageForState(state manager.TunnelState, size int) (img walk.Image, err error) {
	// Для состояний используем стандартные иконки
	icon, err := iconForState(state, size)
	if err != nil {
		return nil, err
	}

	// Конвертируем иконку в изображение
	return iconToImage(icon, size)
}

// loadPNGImage загружает PNG файл и масштабирует его до нужного размера
func loadPNGImage(path string, size int) (walk.Image, error) {
	cacheKey := widthAndPath{size, path}
	if cached, exists := cachedImagesForWidthAndPath[cacheKey]; exists {
		return cached, nil
	}

	// Загружаем изображение из файла
	img, err := walk.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}

	// Масштабируем изображение если нужно
	imgSize := img.Size()
	if imgSize.Width != size || imgSize.Height != size {
		scaledImg, err := scaleImage(img, size, size)
		if err != nil {
			return nil, err
		}
		cachedImagesForWidthAndPath[cacheKey] = scaledImg
		return scaledImg, nil
	}

	cachedImagesForWidthAndPath[cacheKey] = img
	return img, nil
}

// scaleImage масштабирует изображение до указанных размеров
func scaleImage(src walk.Image, width, height int) (walk.Image, error) {
	return walk.NewPaintFuncImage(walk.Size{width, height}, func(canvas *walk.Canvas, bounds walk.Rectangle) error {
		return canvas.DrawImageStretched(src, bounds)
	}), nil
}

// iconToImage конвертирует иконку в изображение
func iconToImage(icon *walk.Icon, size int) (walk.Image, error) {
	return walk.NewPaintFuncImage(walk.Size{size, size}, func(canvas *walk.Canvas, bounds walk.Rectangle) error {
		return canvas.DrawIconStretched(icon, bounds)
	}), nil
}

// loadLogoImage загружает логотип из PNG файла
func loadLogoImage(size int) (walk.Image, error) {
	// Загружаем PNG логотип из конкретного пути
	img, err := loadPNGImage("icon/wireguard.png", size)
	if err != nil {
		// Fallback к старому методу загрузки иконки если PNG не найден
		icon, err := loadLogoIcon(size)
		if err != nil {
			return nil, err
		}
		return iconToImage(icon, size)
	}
	return img, nil
}

func textForState(state manager.TunnelState, withEllipsis bool) (text string) {
	switch state {
	case manager.TunnelStarted:
		text = l18n.Sprintf("Active")
	case manager.TunnelStarting:
		text = l18n.Sprintf("Activating")
	case manager.TunnelStopped:
		text = l18n.Sprintf("Inactive")
	case manager.TunnelStopping:
		text = l18n.Sprintf("Deactivating")
	case manager.TunnelUnknown:
		text = l18n.Sprintf("Unknown state")
	}
	if withEllipsis {
		switch state {
		case manager.TunnelStarting, manager.TunnelStopping:
			text += "…"
		}
	}
	return
}

var cachedSystemIconsForWidthAndDllIdx = make(map[widthAndDllIdx]*walk.Icon)

func loadSystemIcon(dll string, index int32, size int) (icon *walk.Icon, err error) {
	icon = cachedSystemIconsForWidthAndDllIdx[widthAndDllIdx{size, index, dll}]
	if icon != nil {
		return
	}
	icon, err = walk.NewIconFromSysDLLWithSize(dll, int(index), size)
	if err == nil {
		cachedSystemIconsForWidthAndDllIdx[widthAndDllIdx{size, index, dll}] = icon
	}
	return
}

func loadShieldIcon(size int) (icon *walk.Icon, err error) {
	icon, err = loadSystemIcon("imageres", -1028, size)
	if err != nil {
		icon, err = loadSystemIcon("imageres", 1, size)
	}
	return
}

var cachedLogoIconsForWidth = make(map[int]*walk.Icon)

func loadLogoIcon(size int) (icon *walk.Icon, err error) {
	icon = cachedLogoIconsForWidth[size]
	if icon != nil {
		return
	}
	icon, err = walk.NewIconFromResourceIdWithSize(7, walk.Size{size, size})
	if err == nil {
		cachedLogoIconsForWidth[size] = icon
	}
	return
}
