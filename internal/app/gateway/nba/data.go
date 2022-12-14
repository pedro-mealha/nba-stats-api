package nba

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	gameDateFormat          = "2006-01-02"
	gameTimeFormat          = "2006-01-02T15:04:05Z"
	gameHourFormat          = "15:04 UTC"
	NBA            LeagueID = "00"
	WNBA           LeagueID = "10"
)

type (
	LeagueID string
	GameDate time.Time
	GameTime time.Time

	GetScoreboardCommand struct {
		Date     string
		LeagueID LeagueID
	}

	GetBoxscoreCommand struct {
		GameID string
	}

	ScoreboardData struct {
		Scoreboard Scoreboard `json:"scoreboard"`
	}

	Scoreboard struct {
		GameDate GameDate `json:"gameDate"`
		Games    []Game   `json:"games"`
	}

	Game struct {
		ID       string   `json:"gameId"`
		StartsAt GameTime `json:"gameTimeUTC"`
		HomeTeam Team     `json:"homeTeam"`
		AwayTeam Team     `json:"awayTeam"`
	}

	Team struct {
		ID      int64    `json:"teamId"`
		Name    string   `json:"teamName"`
		Tricode string   `json:"teamTricode"`
		Stats   Stats    `json:"statistics"`
		Players []Player `json:"players"`
	}

	BoxscoreData struct {
		Boxscore Boxscore `json:"game"`
	}

	Boxscore struct {
		ID       string `json:"gameId"`
		HomeTeam Team   `json:"homeTeam"`
		AwayTeam Team   `json:"awayTeam"`
	}

	Player struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"familyName"`
		Position  string `json:"position"`
		Stats     Stats  `json:"statistics"`
	}

	Stats struct {
		Minutes   string  `json:"minutes"`
		FGM       int64   `json:"fieldGoalsMade"`
		FGA       int64   `json:"fieldGoalsAttempted"`
		FGP       float64 `json:"fieldGoalsPercentage"`
		ThreeFGM  int64   `json:"threePointersMade"`
		ThreeFGA  int64   `json:"threePointersAttempted"`
		ThreeFGP  float64 `json:"threePointersPercentage"`
		FTM       int64   `json:"freeThrowsMade"`
		FTA       int64   `json:"freeThrowsAttempted"`
		FTP       float64 `json:"freeThrowsPercentage"`
		RO        int64   `json:"reboundsOffensive"`
		RD        int64   `json:"reboundsDefensive"`
		RT        int64   `json:"reboundsTotal"`
		RTeam     int64   `json:"reboundsTeam"`
		AST       int64   `json:"assists"`
		STL       int64   `json:"steals"`
		BLK       int64   `json:"blocks"`
		TO        int64   `json:"turnovers"`
		TOT       int64   `json:"turnoversTeam"`
		FP        int64   `json:"foulsPersonal"`
		FD        int64   `json:"foulsDrawn"`
		PT        int64   `json:"points"`
		PlusMinus float64 `json:"plusMinusPoints"`
	}
)

func (gd GameDate) String() string {
	return time.Time(gd).String()
}

func (gd GameDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(gd).Format(gameDateFormat))

	return []byte(stamp), nil
}

func (gd *GameDate) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == "" || str == "null" {
		return nil
	}

	parsed, err := time.Parse(gameDateFormat, str)
	if err != nil {
		return err
	}

	*gd = GameDate(parsed)

	return nil
}

func (gt GameTime) String() string {
	return time.Time(gt).String()
}

func (gt GameTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(gt).Format(gameHourFormat))

	return []byte(stamp), nil
}

func (gt *GameTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == "" || str == "null" {
		return nil
	}

	parsed, err := time.Parse(gameTimeFormat, str)
	if err != nil {
		return err
	}

	*gt = GameTime(parsed)

	return nil
}

func ParseLeague(l string) LeagueID {
	switch l {
	case "wnba":
		return WNBA
	case "nba":
		fallthrough
	default:
		return NBA
	}
}
