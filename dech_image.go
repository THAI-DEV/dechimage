package dechimage

import (
	"bufio"
	"unicode/utf8"

	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type ImageConfigImgType struct {
	Dpi          float64
	FontFile     string
	Hinting      string
	Size         float64
	Spacing      float64
	WhiteOnBlack bool
}

type service struct {
	imageConfig ImageConfigImgType
}

func New(imageConfig ImageConfigImgType) *service {
	return &service{
		imageConfig: imageConfig,
	}
}

func (rcv *service) GenImage(title string, dataTxt []string, footer string, imageW int, imageH int, outputImage string) bool {
	// Read the font data.
	isOk := false

	colLen := utf8.RuneCountInString(dataTxt[0])

	titleLen := utf8.RuneCountInString(title)
	if titleLen > colLen {
		colLen = titleLen
	}

	footerLen := utf8.RuneCountInString(footer)
	if footerLen > colLen {
		colLen = footerLen
	}

	numRow := len(dataTxt)

	imageW = calAutoWidth(colLen, imageW)
	imageH = calAutoHight(numRow, imageH)

	fontBytes, err := os.ReadFile(rcv.imageConfig.FontFile)
	if err != nil {
		log.Println(err)
		return isOk
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return isOk
	}

	// Draw the background and the guidelines.
	fg, bg := image.Black, image.White
	// ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	if rcv.imageConfig.WhiteOnBlack {
		fg, bg = image.White, image.Black
		// ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	}

	rgba := image.NewRGBA(image.Rect(0, 0, imageW, imageH))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	// for i := 0; i < 200; i++ {
	// 	rgba.Set(10, 10+i, ruler)
	// 	rgba.Set(10+i, 10, ruler)
	// }

	// Draw the text.
	h := font.HintingNone
	switch rcv.imageConfig.Hinting {
	case "full":
		h = font.HintingFull
	}
	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    rcv.imageConfig.Size,
			DPI:     rcv.imageConfig.Dpi,
			Hinting: h,
		}),
	}
	y := 10 + int(math.Ceil(rcv.imageConfig.Size*rcv.imageConfig.Dpi/72))
	dy := int(math.Ceil(rcv.imageConfig.Size * rcv.imageConfig.Spacing * rcv.imageConfig.Dpi / 72))
	d.Dot = fixed.Point26_6{
		X: (fixed.I(imageW) - d.MeasureString(title)) / 2,
		Y: fixed.I(y),
	}

	if len(title) > 0 {
		d.DrawString(title)
		y += dy
	}

	for _, s := range dataTxt {
		d.Dot = fixed.P(10, y)
		d.DrawString(s)
		y += dy
	}

	//* By DECH
	if len(footer) > 0 {
		d.Dot = fixed.P(20, y)
		y += dy
		d.DrawString(footer)
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create(outputImage)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	isOk = true

	return isOk
	// applog.WriteInfo("Write Image OK.")
	// log.Println("Write Image OK.")
}

func calAutoHight(numRow int, height int) int {
	autoHight := height

	higthPerRow := 18
	higthExtend := 0

	if autoHight == 0 {
		if numRow == 0 {
			higthExtend = 100
		} else {
			if numRow >= 0 && numRow <= 3 {
				higthExtend = 45
			}

			if numRow > 3 {
				higthExtend = higthExtend - ((numRow - 25) * 2)
			}
		}

		autoHight = (higthPerRow * numRow) + higthExtend
	}

	// fmt.Printf("Row : %d  . Hight : %d\n", numRow, autoHight)

	return autoHight
}

func calAutoWidth(colLen int, width int) int {
	autoWidh := width
	widthPerLen := 10
	// widthExtend := 0

	if autoWidh == 0 {
		// if colLen > 45 {
		// 	widthExtend = ((colLen - 45) * 2) - widthExtend
		// }
		// autoWidh = (widthPerLen * colLen) - widthExtend
		autoWidh = (widthPerLen * colLen)

		if colLen > 45 {
			autoWidh += 10
		} else {
			autoWidh += 15
		}
	}

	// fmt.Printf("Col Len : %d  , width : %d\n", colLen, autoWidh)
	return autoWidh
}
