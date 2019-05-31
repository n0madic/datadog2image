// Package datadog2image take screenshot from DataDog public dashboard
package datadog2image

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/url"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"

	"github.com/chromedp/chromedp"
	"github.com/golang/freetype/truetype"
)

// Dashboard object
type Dashboard struct {
	URL        string
	screenshot []byte
	Error      error
}

// NewDashboard return a new Datadog dashboard object
func NewDashboard(source string) *Dashboard {
	sourceURL, err := url.Parse(source)
	if err != nil {
		return nil
	}
	return &Dashboard{URL: sourceURL.String()}
}

// GetScreenshot from Datadog dashboard via headless Chrome
func (d *Dashboard) GetScreenshot(waitLoading int64) *Dashboard {
	sel := `#sub_board`

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Duration(waitLoading+10)*time.Second)
	defer cancel()

	d.Error = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(d.URL),
		chromedp.WaitVisible(sel, chromedp.ByID),
		chromedp.Sleep(time.Second * time.Duration(waitLoading)),
		chromedp.Screenshot(sel, &d.screenshot, chromedp.NodeVisible, chromedp.ByID),
	})
	return d
}

// AddTimestamp with the current time in the screenshot
func (d *Dashboard) AddTimestamp(timestamp *time.Time) *Dashboard {
	if len(d.screenshot) > 0 {
		img, err := png.Decode(bytes.NewReader(d.screenshot))
		if err != nil {
			d.Error = err
			return d
		}

		b := img.Bounds()
		rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()+40))
		draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)
		addLabel(rgba, 20, b.Dy()+28, timestamp.Format("Mon, 02 Jan 2006 15:04"))

		var buf bytes.Buffer
		err = png.Encode(&buf, rgba)
		if err != nil {
			d.Error = err
		} else {
			d.screenshot = buf.Bytes()
		}
	}
	return d
}

// PNG return bytes array with screenshot
func (d *Dashboard) PNG() []byte {
	return d.screenshot
}

// HTML return index page with embedded screenshot
func (d *Dashboard) HTML(refresh int) []byte {
	return []byte(fmt.Sprintf(`<html>
<head>
	<meta http-equiv="refresh" content="%d">
	<style>
		* {
			margin: 0;
			padding: 0;
		}
		.imgbox {
			display: grid;
			height: 100%%;
		}
		.center-fit {
			max-width: 100%%;
			max-height: 100vh;
			margin: auto;
		}
	</style>
</head>
<body>
	<div class="imgbox">
		<img class="center-fit" src="data:image/png;base64,%s" />
	</div>
</body>
</html>`, refresh, base64.StdEncoding.EncodeToString(d.screenshot)))
}

// Add label on image
func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{0, 0, 0, 128}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	regular, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	faceRegular := truetype.NewFace(regular, &truetype.Options{Size: 28})

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: faceRegular,
		Dot:  point,
	}
	d.DrawString(label)
}
