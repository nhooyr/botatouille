package music

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/log"
	"runtime"
)

type guildMusic struct {
	q           *queue
	framesChan  chan []byte
	dvc         *discordgo.VoiceConnection
	startChan   chan *video
	// TODO gotta make it a single channel, can't have this shit happen concurrenctly
	pauseChan   chan struct{}
	unPauseChan chan struct{}
	stopChan    chan struct{}
	youtubeDL   *exec.Cmd
	ffmpeg      *exec.Cmd
}

func newGuildMusic(dvc *discordgo.VoiceConnection) *guildMusic {
	gm := &guildMusic{
		// TODO needs to be a linked list
		framesChan:  make(chan []byte, 100000),
		dvc:         dvc,
		startChan:   make(chan *video),
		pauseChan:   make(chan struct{}),
		unPauseChan: make(chan struct{}),
		stopChan:    make(chan struct{}),
	}
	gm.q = newQueue(gm.startChan)
	go gm.play()
	return gm
}

func (gm *guildMusic) play() {
	go gm.opusSender()
	for {
		select {
		case v := <-gm.startChan:
			gm.play_id(v.id)
			for {
				v := gm.q.next()
				if v == nil {
					break
				}
				gm.play_id(v.id)
			}
		case <-gm.stopChan:
			return
		}
	}
}

func (gm *guildMusic) opusSender() {
	for {
		select {
		case f := <-gm.framesChan:
			gm.dvc.OpusSend <- f
		case <-gm.pauseChan:
			select {
			case <-gm.unPauseChan:
				continue
			case <-gm.stopChan:
				return
			}
		case <-gm.stopChan:
			return
		}
	}
}

func (gm *guildMusic) unPause() {
	gm.unPauseChan <- struct{}{}
}

func (gm *guildMusic) pause() {
	gm.pauseChan <- struct{}{}
}

func (gm *guildMusic) stop() error {
	// One for opusSender and the other for play_id if it's running
	// TODO actually need 3
	gm.stopChan <- struct{}{}
	gm.stopChan <- struct{}{}
	select {
	case gm.stopChan <- struct{}{}:
		err := gm.youtubeDL.Process.Kill()
		if err != nil {
			log.Print(err)
		}
		err = gm.ffmpeg.Process.Kill()
		if err != nil {
			log.Print(err)
		}
	default:
	}
	return gm.dvc.Disconnect()
}

// TODO how to handle errors?
func (gm *guildMusic) play_id(id string) {
	// TODO limit number of bytes read
	url := fmt.Sprintf("https://youtu.be/%v", id)
	youtubeDL := exec.Command("youtube-dl", "-q", "-f", "bestaudio", "-o", "-", url)
	youtubeDlOut, err := youtubeDL.StdoutPipe()
	if err != nil {
		log.Print(err)
		return
	}
	// TODO read stderr for errors
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0",
		"-f", "data", "-map", "0:a", "-ar", "48k", "-ac", "2",
		"-acodec", "libopus", "-b:a", "128k", "pipe:1")
	ffmpeg.Stdin = youtubeDlOut
	ffmpegOut, err := ffmpeg.StdoutPipe()
	if err != nil {
		log.Print(err)
		return
	}
	runtime.LockOSThread()
	err = youtubeDL.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = ffmpeg.Start()
	if err != nil {
		log.Fatal(err)
	}
	// TODO store process state in music handler
	for {
		select {
		case <-gm.stopChan:
			return
		default:
		}
		// I read in the RFC that frames will not be bigger than this size
		p := make([]byte, 1275)
		n, err := ffmpegOut.Read(p)
		if err != nil {
			if err == io.EOF {
				return
			}
			go gm.stop()
			<-gm.stopChan
			return
		}
		gm.framesChan <- p[:n]
	}
}
