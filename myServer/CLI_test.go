package poker_test

import (
	"bytes"
	"fmt"
	poker "go-practice/myServer"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

type SpyBlindAlerter struct {
	alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, ScheduledAlert{at, amount})
}

type GameSpy struct {
	StartedWith    int
	FinishedWith   string
	StartCalled    bool
	FinishedCalled bool
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
	g.StartCalled = true
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
	g.FinishedCalled = true
}

var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

func TestCLI(t *testing.T) {

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("3", "Chris"+poker.WinsMsg)
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})
	t.Run("shouldn't register winner if the user entry is not correct", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("3", "Chris is a maniac")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		asserGameNotFInished(t, game)
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &GameSpy{}

		in := userSends("8", "Cleo"+poker.WinsMsg)
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		game := &GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("pies")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func assertScheduledAlert(t testing.TB, got ScheduledAlert, want ScheduledAlert) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Alert %s is not equal to %s", got, want)
	}
}

func assertMessagesSentToUser(t testing.TB, out *bytes.Buffer, messages ...string) {
	t.Helper()
	for _, message := range messages {
		if !bytes.ContainsAny(out.Bytes(), message) {
			t.Errorf("Message %s is missing in %s", message, out.String())
		}
	}
}
func assertGameNotStarted(t testing.TB, game *GameSpy) {
	if game.StartCalled {
		t.Errorf("Game has started when it shouldn't had")

	}
}

func assertGameStartedWith(t testing.TB, game *GameSpy, num int) {
	if game.StartedWith != num {
		t.Errorf("Game started with wrong number %d it should be %d", game.StartedWith, num)
	}
}
func assertFinishCalledWith(t testing.TB, game *GameSpy, player string) {
	if !bytes.ContainsAny([]byte(game.FinishedWith), player) {
		t.Errorf("Game ended with wrong player %s, it should be %s", game.FinishedWith, player)

	}
}

func asserGameNotFInished(t testing.TB, game *GameSpy) {
	if game.FinishedCalled {
		t.Error("Game shouln't have finished")
	}
}
