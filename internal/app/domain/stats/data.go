package stats

import (
	"fmt"
	"regexp"

	"github.com/pedro-mealha/nba-stats-api/internal/app/gateway/nba"
)

const (
	// Matches PT25M12.02S or PT04M05.00S
	minutesRegex = `(?:PT)([1-9]{1}|(?:0))(\d{1})(?:M)(\d{2})(?:\.\d{2}S)`

	zeroMins = "0:00"
)

type (
	Scoreboard struct {
		Date  nba.GameDate `json:"date"`
		Games []Game       `json:"games"`
	}

	Game struct {
		ID       string       `json:"id"`
		StartsAt nba.GameTime `json:"starts_at"`
		HomeTeam Team         `json:"home_team"`
		AwayTeam Team         `json:"away_team"`
	}

	Boxscore struct {
		GameID   string `json:"game_id"`
		HomeTeam Team   `json:"home_team"`
		AwayTeam Team   `json:"away_team"`
	}

	Team struct {
		ID      int64    `json:"id"`
		Name    string   `json:"name"`
		Tricode string   `json:"tricode"`
		Stats   Stats    `json:"stats"`
		Players []Player `json:"players,omitempty"`
	}

	Player struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Position  string `json:"position"`
		Stats     Stats  `json:"stats"`
	}

	Stats struct {
		Minutes   string  `json:"min"`
		FGM       int64   `json:"fgm"`
		FGA       int64   `json:"fga"`
		FGP       float64 `json:"fgp"`
		ThreeFGM  int64   `json:"3fgm"`
		ThreeFGA  int64   `json:"3fga"`
		ThreeFGP  float64 `json:"3fgp"`
		FTM       int64   `json:"ftm"`
		FTA       int64   `json:"fta"`
		FTP       float64 `json:"ftp"`
		RO        int64   `json:"oreb"`
		RD        int64   `json:"dreb"`
		RT        int64   `json:"reb"`
		RTeam     int64   `json:"rebt"`
		AST       int64   `json:"ast"`
		STL       int64   `json:"stl"`
		BLK       int64   `json:"blk"`
		TO        int64   `json:"to"`
		TOT       int64   `json:"tot"`
		FP        int64   `json:"pf"`
		FD        int64   `json:"fd"`
		PT        int64   `json:"pts"`
		PlusMinus float64 `json:"plus_minus"`
	}
)

func NewScoreboard(sb nba.ScoreboardData) Scoreboard {
	return Scoreboard{
		Date:  sb.Scoreboard.GameDate,
		Games: addGames(sb.Scoreboard.Games),
	}
}

func NewBoxscore(bs nba.BoxscoreData) Boxscore {
	b := Boxscore{
		GameID: bs.Boxscore.ID,
		HomeTeam: Team{
			ID:      bs.Boxscore.HomeTeam.ID,
			Name:    bs.Boxscore.HomeTeam.Name,
			Tricode: bs.Boxscore.HomeTeam.Tricode,
			Stats:   statsDecorator(bs.Boxscore.HomeTeam.Stats),
			Players: addPlayers(bs.Boxscore.HomeTeam.Players),
		},
		AwayTeam: Team{
			ID:      bs.Boxscore.AwayTeam.ID,
			Name:    bs.Boxscore.AwayTeam.Name,
			Tricode: bs.Boxscore.AwayTeam.Tricode,
			Stats:   statsDecorator(bs.Boxscore.AwayTeam.Stats),
			Players: addPlayers(bs.Boxscore.AwayTeam.Players),
		},
	}

	return b
}

func addGames(gs []nba.Game) []Game {
	gg := make([]Game, len(gs))

	for i, g := range gs {
		gg[i] = Game{
			ID:       g.ID,
			StartsAt: g.StartsAt,
			HomeTeam: Team{
				ID:      g.HomeTeam.ID,
				Name:    g.HomeTeam.Name,
				Tricode: g.HomeTeam.Tricode,
			},
			AwayTeam: Team{
				ID:      g.AwayTeam.ID,
				Name:    g.AwayTeam.Name,
				Tricode: g.AwayTeam.Tricode,
			},
		}
	}

	return gg
}

func addPlayers(ps []nba.Player) []Player {
	pp := make([]Player, len(ps))

	for i, p := range ps {
		pp[i] = Player{
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Position:  p.Position,
			Stats:     statsDecorator(p.Stats),
		}
	}

	return pp
}

func parseMinutes(min string) string {
	re := regexp.MustCompile(minutesRegex)
	if match := re.FindStringSubmatch(min); len(match) > 0 {
		if match[1] == "0" {
			match[1] = ""
		}

		return fmt.Sprintf("%s%s:%s", match[1], match[2], match[3])
	}

	return zeroMins
}

func parsePercentages(p float64) float64 {
	return p * 100
}

func statsDecorator(s nba.Stats) Stats {
	sts := Stats(s)

	sts.FGP = parsePercentages(sts.FGP)
	sts.FTP = parsePercentages(sts.FTP)
	sts.ThreeFGP = parsePercentages(sts.ThreeFGP)

	sts.Minutes = parseMinutes(sts.Minutes)

	return sts
}
