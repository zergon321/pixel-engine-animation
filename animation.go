package main

import (
	"time"

	"github.com/faiface/pixel"
)

type animation struct {
	spritesheet  pixel.Picture
	frames       []pixel.Rect
	delays       []time.Duration
	timeout      <-chan time.Time
	cancel       chan interface{}
	currentFrame int
	active       bool
	loop         bool
}

func (anim *animation) start() {
	if anim.active {
		return
	}

	anim.currentFrame = 0
	anim.timeout = time.After(anim.delays[anim.currentFrame])
	anim.cancel = make(chan interface{})
	anim.active = true

	go anim.process()
}

func (anim *animation) process() {
	for {
		select {
		case <-anim.timeout:
			anim.currentFrame++

			if anim.currentFrame >= len(anim.frames) {
				anim.currentFrame = 0

				if !anim.loop {
					anim.active = false
					close(anim.cancel)

					anim.cancel = nil
					anim.timeout = nil

					return
				}
			}

			anim.timeout = time.After(anim.delays[anim.currentFrame])

		case <-anim.cancel:
			anim.currentFrame = 0
			anim.active = false

			close(anim.cancel)
			anim.cancel = nil
			anim.timeout = nil

			return
		}
	}
}

func (anim *animation) stop() {
	if anim.active {
		go func() {
			anim.cancel <- true
		}()
	}
}

func (anim *animation) getSprite(frameInd int) *pixel.Sprite {
	return pixel.NewSprite(anim.spritesheet, anim.frames[frameInd])
}

func (anim *animation) getCurrentSprite() *pixel.Sprite {
	return anim.getSprite(anim.currentFrame)
}
