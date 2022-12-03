package stats

import (
	"fmt"
	"regexp"

	"github.com/WeNeedThePoh/nba-stats-api/internal/app/gateway/nba"
)

const minutesRegex = `(?:PT)(\d{2})(?:M)(\d{2})(?:\.\d{2}S)`

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
		Slug      string `json:"slug"`
		Position  string `json:"position"`
		Stats     Stats  `json:"stats"`
	}

	Stats struct {
		Minutes  string  `json:"min"`
		FGM      int64   `json:"fgm"`
		FGA      int64   `json:"fga"`
		FGP      float64 `json:"fgp"`
		ThreeFGM int64   `json:"3fgm"`
		ThreeFGA int64   `json:"3fga"`
		ThreeFGP float64 `json:"3fgp"`
		FTM      int64   `json:"ftm"`
		FTA      int64   `json:"fta"`
		FTP      float64 `json:"ftp"`
		RO       int64   `json:"oreb"`
		RD       int64   `json:"dreb"`
		RT       int64   `json:"reb"`
		AST      int64   `json:"ast"`
		STL      int64   `json:"stl"`
		BLK      int64   `json:"blk"`
		TO       int64   `json:"to"`
		FP       int64   `json:"pf"`
		PT       int64   `json:"pts"`
	}
)

func NewScoreboard(sb nba.ScoreboardData) Scoreboard {
	return Scoreboard{
		Date:  sb.Scoreboard.GameDate,
		Games: addGames(sb.Scoreboard.Games),
	}
}

func NewBoxscore(bs nba.BoxscoreData) Boxscore {
	return Boxscore{
		GameID: bs.Boxscore.ID,
		HomeTeam: Team{
			ID:      bs.Boxscore.HomeTeam.ID,
			Name:    bs.Boxscore.HomeTeam.Name,
			Tricode: bs.Boxscore.HomeTeam.Tricode,
			Stats:   Stats(bs.Boxscore.HomeTeam.Stats),
			Players: addPlayers(bs.Boxscore.HomeTeam.Players),
		},
		AwayTeam: Team{
			ID:      bs.Boxscore.AwayTeam.ID,
			Name:    bs.Boxscore.AwayTeam.Name,
			Tricode: bs.Boxscore.AwayTeam.Tricode,
			Stats:   Stats(bs.Boxscore.AwayTeam.Stats),
			Players: addPlayers(bs.Boxscore.AwayTeam.Players),
		},
	}
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
			Slug:      p.Slug,
			Position:  p.Position,
			Stats:     Stats(p.Stats),
		}

		pp[i].Stats.Minutes = parseMinutes(pp[i].Stats.Minutes)
	}

	return pp
}

func parseMinutes(min string) string {
	re := regexp.MustCompile(minutesRegex)
	if match := re.FindStringSubmatch(min); len(match) > 0 {
		return fmt.Sprintf("%s:%s", match[1], match[2])
	}

	return "00:00"
}
