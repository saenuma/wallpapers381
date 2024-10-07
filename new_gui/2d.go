package main

import (
	"image"
	"strconv"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
)

type Ctx struct {
	WindowWidth  int
	WindowHeight int
	ggCtx        *gg.Context
	OldFrame     image.Image
}

func New2dCtx(wWidth, wHeight int) Ctx {
	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := getDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: wWidth, WindowHeight: wHeight, ggCtx: ggCtx}
	return ctx
}

func (ctx *Ctx) drawButtonA(btnId, originX, originY int, text, textColor, bgColor string) g143.Rect {
	// draw bounding rect
	textW, textH := ctx.ggCtx.MeasureString(text)
	width, height := textW+50, textH+25
	ctx.ggCtx.SetHexColor(bgColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(height))
	ctx.ggCtx.Fill()

	// draw text
	ctx.ggCtx.SetHexColor(textColor)
	ctx.ggCtx.DrawString(text, float64(originX)+25, float64(originY)+fontSize+5)

	// save dimensions
	btnARect := g143.NewRect(originX, originY, int(width), int(height))
	objCoords[btnId] = btnARect
	return btnARect
}

func (ctx *Ctx) drawButtonB(btnId int, text, textColor string) g143.Rect {
	textW, textH := ctx.ggCtx.MeasureString(text)
	originX := (ctx.WindowWidth - int(textW)) / 2
	originY := ctx.WindowHeight - int(textH) - 30

	// draw bounding rect
	width, height := textW+50, textH+25
	ctx.ggCtx.SetHexColor("#fff")
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(height))
	ctx.ggCtx.Fill()

	// draw text
	ctx.ggCtx.SetHexColor(textColor)
	ctx.ggCtx.DrawString(text, float64(originX)+25, float64(originY)+fontSize+5)

	// save dimensions
	btnARect := g143.NewRect(originX, originY, int(width), int(height))
	objCoords[btnId] = btnARect
	return btnARect
}

func (ctx *Ctx) drawInput(inputId, originX, originY, writtenNum int) g143.Rect {
	ctx.ggCtx.SetHexColor("#fff")
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), 100, float64(originY)+fontSize+5)
	ctx.ggCtx.Fill()

	ctx.ggCtx.SetHexColor("#909BD0")
	ctx.ggCtx.DrawRectangle(float64(originX), 50, 100, 3)
	ctx.ggCtx.Fill()

	entryRect := g143.Rect{Width: 100, Height: 50, OriginX: originX, OriginY: 10}
	objCoords[inputId] = entryRect

	lineNoStr := strconv.Itoa(writtenNum)
	ctx.ggCtx.SetHexColor("#444")
	ctx.ggCtx.DrawString(lineNoStr, float64(originX+15), float64(originY)+fontSize+5)
	return entryRect
}

func nextHorizontalCoords(aRect g143.Rect, margin int) (int, int) {
	nextOriginX := aRect.OriginX + aRect.Width + margin
	nextOriginY := aRect.OriginY
	return nextOriginX, nextOriginY
}
