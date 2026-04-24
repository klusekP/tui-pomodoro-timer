package alarm

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type Alarm interface {
	Play() error
}

type Noop struct{}

func (Noop) Play() error { return nil }

type Bell struct{}

func (Bell) Play() error {
	_, err := fmt.Fprint(os.Stderr, "\a\a\a")
	return err
}

type Mac struct {
	SoundPath string
}

func (m Mac) Play() error {
	path := m.SoundPath
	if path == "" {
		path = "/System/Library/Sounds/Glass.aiff"
	}
	return exec.Command("afplay", path).Run()
}

type Linux struct{}

func (Linux) Play() error {
	candidates := [][]string{
		{"paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga"},
		{"paplay", "/usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga"},
		{"aplay", "/usr/share/sounds/alsa/Front_Center.wav"},
	}
	for _, c := range candidates {
		if _, err := exec.LookPath(c[0]); err != nil {
			continue
		}
		if err := exec.Command(c[0], c[1:]...).Run(); err == nil {
			return nil
		}
	}
	return nil
}

type Windows struct{}

func (Windows) Play() error {
	return exec.Command(
		"powershell", "-c",
		"[console]::beep(880,400); [console]::beep(660,400); [console]::beep(880,600)",
	).Run()
}

type Composite struct {
	alarms []Alarm
}

func NewComposite(alarms ...Alarm) Composite {
	return Composite{alarms: alarms}
}

func (c Composite) Play() error {
	for _, a := range c.alarms {
		if a == nil {
			continue
		}
		_ = a.Play()
	}
	return nil
}

func NewSystem() Alarm {
	alarms := []Alarm{Bell{}}
	switch runtime.GOOS {
	case "darwin":
		alarms = append(alarms, Mac{})
	case "linux":
		alarms = append(alarms, Linux{})
	case "windows":
		alarms = append(alarms, Windows{})
	}
	return NewComposite(alarms...)
}
