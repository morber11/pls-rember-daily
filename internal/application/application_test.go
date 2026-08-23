package application

import (
	"errors"
	"strings"
	"testing"
	"time"

	"plsremberdaily/assets"
)

var testVideo = assets.Video{Name: "test.mp4"}

type fakePlatform struct {
	acquired         bool
	acquireError     error
	startupError     error
	playError        error
	registered       bool
	released         bool
	playCalled       bool
	displayed        bool
	played           assets.Video
	markDisplayOrder *[]string
}

func (fake *fakePlatform) Acquire() (func(), bool, error) {
	if fake.acquireError != nil {
		return nil, false, fake.acquireError
	}
	if !fake.acquired {
		return nil, false, nil
	}

	return func() { fake.released = true }, true, nil
}

func (fake *fakePlatform) RegisterStartup() error {
	fake.registered = true
	return fake.startupError
}

func (fake *fakePlatform) Play(video assets.Video, onStarted func() error) error {
	fake.playCalled = true
	fake.played = video
	if err := onStarted(); err != nil {
		return err
	}

	fake.displayed = true
	*fake.markDisplayOrder = append(*fake.markDisplayOrder, "display")
	return fake.playError
}

type fakeReminder struct {
	multipleRuns     bool
	shouldShow       bool
	shouldShowErr    error
	markErr          error
	checked          bool
	marked           bool
	markDisplayOrder *[]string
}

func (fake *fakeReminder) AllowMultipleRuns() bool {
	return fake.multipleRuns
}

func (fake *fakeReminder) ShouldShowToday(time.Time) (bool, error) {
	fake.checked = true
	return fake.shouldShow, fake.shouldShowErr
}

func (fake *fakeReminder) MarkShownToday(time.Time) error {
	fake.marked = true
	*fake.markDisplayOrder = append(*fake.markDisplayOrder, "mark")
	return fake.markErr
}

type fakeSchedule struct {
	video    assets.Video
	selected bool
}

func (fake *fakeSchedule) Video(time.Time) assets.Video {
	fake.selected = true
	return fake.video
}

type fakeClock struct {
	now            time.Time
	sleepDurations []time.Duration
}

func (fake *fakeClock) Now() time.Time {
	return fake.now
}

func (fake *fakeClock) Sleep(duration time.Duration) {
	fake.sleepDurations = append(fake.sleepDurations, duration)
}

type logEntry struct {
	format string
	args   []any
}

type fakeLogger struct {
	entries []logEntry
}

func (fake *fakeLogger) Printf(format string, args ...any) {
	fake.entries = append(fake.entries, logEntry{format: format, args: args})
}

type runnerFakes struct {
	platform *fakePlatform
	reminder *fakeReminder
	schedule *fakeSchedule
	clock    *fakeClock
	logger   *fakeLogger
}

func newTestRunner() (Runner, runnerFakes) {
	markDisplayOrder := []string{}
	fakes := runnerFakes{
		platform: &fakePlatform{
			acquired:         true,
			markDisplayOrder: &markDisplayOrder,
		},
		reminder: &fakeReminder{
			shouldShow:       true,
			markDisplayOrder: &markDisplayOrder,
		},
		schedule: &fakeSchedule{video: testVideo},
		clock: &fakeClock{
			now: time.Date(2026, time.August, 23, 9, 0, 0, 0, time.Local),
		},
		logger: &fakeLogger{},
	}

	runner := NewRunner(
		fakes.platform,
		fakes.reminder,
		fakes.schedule,
		fakes.clock,
		fakes.logger,
	)

	return runner, fakes
}

func TestRunAcquireError(t *testing.T) {
	runner, fakes := newTestRunner()
	wantErr := errors.New("acquire")
	fakes.platform.acquireError = wantErr

	if err := runner.Run(); !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
	if fakes.platform.registered || fakes.platform.playCalled || fakes.platform.released {
		t.Fatal("Run() continued after instance acquisition failed")
	}
}

func TestRunExistingInstance(t *testing.T) {
	runner, fakes := newTestRunner()
	fakes.platform.acquired = false

	if err := runner.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if fakes.platform.registered || fakes.platform.playCalled || fakes.platform.released {
		t.Fatal("Run() continued without acquiring the instance")
	}
}

func TestRunSkipsWhenAlreadyShown(t *testing.T) {
	runner, fakes := newTestRunner()
	fakes.reminder.shouldShow = false

	if err := runner.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !fakes.reminder.checked {
		t.Fatal("Run() did not check whether the reminder was already shown")
	}
	if len(fakes.clock.sleepDurations) != 0 || fakes.schedule.selected || fakes.platform.playCalled {
		t.Fatal("Run() continued after the reminder was already shown")
	}
	if !fakes.platform.released {
		t.Fatal("Run() did not release the acquired instance")
	}
}

func TestRunSuccessfulFlow(t *testing.T) {
	runner, fakes := newTestRunner()

	if err := runner.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if startupDelay == 0 {
		if len(fakes.clock.sleepDurations) != 0 {
			t.Fatalf("sleep durations = %v, want none", fakes.clock.sleepDurations)
		}
	} else if len(fakes.clock.sleepDurations) != 1 || fakes.clock.sleepDurations[0] != startupDelay {
		t.Fatalf("sleep durations = %v, want [%v]", fakes.clock.sleepDurations, startupDelay)
	}

	if !fakes.schedule.selected || fakes.platform.played.Name != testVideo.Name {
		t.Fatalf("played video = %q, want %q", fakes.platform.played.Name, testVideo.Name)
	}

	if got := *fakes.platform.markDisplayOrder; len(got) != 2 || got[0] != "mark" || got[1] != "display" {
		t.Fatalf("mark/display order = %v, want [mark display]", got)
	}

	if !fakes.platform.released {
		t.Fatal("Run() did not release the acquired instance")
	}
}

func TestRunStartupErrorContinues(t *testing.T) {
	runner, fakes := newTestRunner()
	startupErr := errors.New("startup")
	fakes.platform.startupError = startupErr
	fakes.reminder.multipleRuns = true

	if err := runner.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !fakes.platform.displayed {
		t.Fatal("Run() stopped after startup registration failed")
	}
	if len(fakes.logger.entries) != 1 {
		t.Fatalf("log entries = %d, want 1", len(fakes.logger.entries))
	}

	entry := fakes.logger.entries[0]
	if !strings.Contains(entry.format, "startup") {
		t.Fatalf("log category = %q, want startup", entry.format)
	}
	if len(entry.args) != 1 {
		t.Fatalf("logged arguments = %v, want startup error", entry.args)
	}
	loggedErr, ok := entry.args[0].(error)
	if !ok || !errors.Is(loggedErr, startupErr) {
		t.Fatalf("logged error = %v, want %v", entry.args[0], startupErr)
	}
}

func TestRunMarkErrorSkipsPlayback(t *testing.T) {
	runner, fakes := newTestRunner()
	markErr := errors.New("mark")
	fakes.reminder.markErr = markErr

	if err := runner.Run(); !errors.Is(err, markErr) {
		t.Fatalf("Run() error = %v, want %v", err, markErr)
	}
	if !fakes.platform.playCalled {
		t.Fatal("Run() did not prepare playback before marking the reminder")
	}
	if fakes.platform.displayed {
		t.Fatal("Run() displayed the video after marking failed")
	}
	if !fakes.platform.released {
		t.Fatal("Run() did not release the acquired instance")
	}
}

func TestRunPlaybackError(t *testing.T) {
	runner, fakes := newTestRunner()
	playErr := errors.New("play")
	fakes.platform.playError = playErr

	if err := runner.Run(); !errors.Is(err, playErr) {
		t.Fatalf("Run() error = %v, want %v", err, playErr)
	}
	if !fakes.reminder.marked {
		t.Fatal("Run() did not mark the reminder before playback failed")
	}
	if !fakes.platform.released {
		t.Fatal("Run() did not release the acquired instance")
	}
}

func TestRunShouldShowError(t *testing.T) {
	runner, fakes := newTestRunner()
	wantErr := errors.New("should show")
	fakes.reminder.shouldShowErr = wantErr

	if err := runner.Run(); !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
	if len(fakes.clock.sleepDurations) != 0 || fakes.schedule.selected || fakes.platform.playCalled {
		t.Fatal("Run() continued after checking the reminder failed")
	}
	if !fakes.platform.released {
		t.Fatal("Run() did not release the acquired instance")
	}
}
