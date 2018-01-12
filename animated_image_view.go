package grip

import (
	"image"
	"image/gif"
	"os"
	"time"

	termbox "github.com/nsf/termbox-go"
	asciiart "github.com/sparkymat/goasciiart"
	"github.com/sparkymat/grip/event"
)

type AnimationFrame struct {
	Image                image.Image
	scaledAscii          []byte
	DurationMilliseconds int
}

type AnimatedImageView struct {
	setCellFn      SetCellFn
	Frames         []*AnimationFrame
	LoopCount      int
	imageView      ImageView
	rect           Rect
	visibleRect    Rect
	loopsRemaining int
}

func NewAnimatedImageViewForGifFile(filePath string) (AnimatedImageView, error) {
	gifFile, err := os.Open(filePath)
	if err != nil {
		return AnimatedImageView{}, err
	}

	gifImage, err := gif.DecodeAll(gifFile)
	if err != nil {
		return AnimatedImageView{}, err
	}

	frames := []*AnimationFrame{}

	for index, img := range gifImage.Image {
		frames = append(frames, &AnimationFrame{Image: img, DurationMilliseconds: gifImage.Delay[index] * 10})
	}

	return AnimatedImageView{
		Frames:    frames,
		LoopCount: gifImage.LoopCount,
		imageView: ImageView{},
	}, nil
}

func (ai *AnimatedImageView) Initialize(setCellFn SetCellFn) {
	ai.setCellFn = setCellFn
	ai.imageView.Initialize(setCellFn)

	ai.loopsRemaining = ai.LoopCount
	ai.SetFrame(0)
}

func (ai *AnimatedImageView) FrameCount() int {
	return len(ai.Frames)
}

func (ai *AnimatedImageView) SetFrame(idx int) {
	ai.imageView.Image = ai.Frames[idx].Image
	ai.imageView.scaleAscii = ai.Frames[idx].scaledAscii
	ai.Draw()

	if idx == len(ai.Frames)-1 && ai.loopsRemaining == 0 && ai.LoopCount != 0 {
		return
	}

	nextIdx := idx + 1
	if nextIdx >= len(ai.Frames) {
		nextIdx = 0
		ai.loopsRemaining -= 1
		if ai.loopsRemaining == 0 {
			ai.loopsRemaining = 0
		}
	}

	timer := time.NewTimer(time.Millisecond * time.Duration(ai.Frames[idx].DurationMilliseconds))
	go func() {
		<-timer.C
		ai.SetFrame(nextIdx)
	}()
}

func (ai *AnimatedImageView) Resize(rect, visibleRect Rect) {
	ai.rect = rect
	ai.visibleRect = visibleRect
	for _, frame := range ai.Frames {
		var scaleAscii []byte
		if visibleRect.Size.Width > 0 {
			scaleAscii = asciiart.Convert2AsciiOfWidth(frame.Image, visibleRect.Size.Width-1)
		}
		frame.scaledAscii = scaleAscii
	}
	ai.imageView.Resize(rect, visibleRect)
}

func (ai *AnimatedImageView) Draw() {
	ai.imageView.Draw()
}

func (ai *AnimatedImageView) OnEvent(e event.Event) {
}

func (ai *AnimatedImageView) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	ai.setCellFn(Point{x, y}, ColoredRune{ch, fg, bg})
}
