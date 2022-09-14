package graphic

import (
	"log"

	"github.com/tinne26/etxt"
)

func getTxtRenderer() *etxt.Renderer {
	// load font library
	fontLib := etxt.NewFontLibrary()
	_, _, err := fontLib.ParseDirFonts("data/fonts")
	if err != nil {
		log.Fatalf("Error while loading fonts: %s", err.Error())
	}

	// Display all fonts in the font directory
	/*
		fontLib.EachFont(func(name string, _ *etxt.Font) error {
			fmt.Println(name)
			return nil
		})
	*/

	// check that we have the fonts we want
	// (shown for completeness, you don't need this in most cases)
	if !fontLib.HasFont("Raleway ExtraBold") {
		log.Fatal("missing font 1")
	}

	// check that the fonts have the characters we want
	// (shown for completeness, you don't need this in most cases)
	err = fontLib.EachFont(checkMissingRunes)
	if err != nil {
		log.Fatal(err)
	}

	// create a new text renderer and configure it
	txtRenderer := etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	txtRenderer.SetCacheHandler(glyphsCache.NewHandler())
	txtRenderer.SetFont(fontLib.GetFont("Raleway ExtraBold"))
	// txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter)
	txtRenderer.SetAlign(etxt.Top, etxt.Left)
	txtRenderer.SetSizePx(72)

	return txtRenderer
}

// helper used after loading fonts
func checkMissingRunes(name string, font *etxt.Font) error {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789 .,;:!?-()[]"
	missing, err := etxt.GetMissingRunes(font, alphabet)
	if err != nil {
		return err
	}
	if len(missing) > 0 {
		log.Fatalf("Font '%s' missing runes: %s", name, string(missing))
	}
	return nil
}
