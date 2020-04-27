package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	colors "golang.org/x/image/colornames"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

func loadPicture(pic string) (pixel.Picture, error) {
	file, err := os.Open(pic)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel animation",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	handleError(err)

	spritesheet, err := loadPicture("komachi-spritesheet.png")
	handleError(err)
	anim := &animation{
		spritesheet: spritesheet,
		frames: []pixel.Rect{
			pixel.R(0, 0, 90, 90),
			pixel.R(90, 0, 180, 90),
			pixel.R(180, 0, 270, 90),
			pixel.R(270, 0, 360, 90),
		},
		delays: []time.Duration{
			120 * time.Millisecond,
			120 * time.Millisecond,
			120 * time.Millisecond,
			120 * time.Millisecond,
		},
		currentFrame: 0,
		active:       false,
		loop:         true,
	}

	fps := 0
	perSecond := time.Tick(time.Second)

	anim.start()
	/*go func() {
		time.Sleep(200 * time.Millisecond)
		anim.stop()
	}()*/

	for !win.Closed() {
		win.Clear(colors.White)

		pos := pixel.V(screenWidth/2, screenHeight/2)
		transform := pixel.IM.Moved(pos).Scaled(pos, 2)
		sprite := anim.getCurrentSprite()
		sprite.Draw(win, transform)

		fmt.Println(anim.currentFrame)

		win.Update()

		// Show FPS in the window title.
		fps++

		select {
		case <-perSecond:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, fps))
			fps = 0

		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
