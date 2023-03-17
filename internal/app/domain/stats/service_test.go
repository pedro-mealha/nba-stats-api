package stats_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/pedro-mealha/nba-stats-api/internal/app/domain/stats"
	"github.com/pedro-mealha/nba-stats-api/internal/app/gateway/nba"
	"github.com/stretchr/testify/suite"
)

var errFailed = errors.New("failed")

type ServiceTestSuite struct {
	suite.Suite

	nm *nba.APIMock
	s  stats.Provider
}

func (s *ServiceTestSuite) SetupTest() {
	s.nm = new(nba.APIMock)

	s.s = stats.NewService(s.nm)
}

func TestStatsService(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestGetScoreboard() {
	tests := []struct {
		scenario string

		sb     nba.ScoreboardData
		nbaErr error

		expErr error
		expRes stats.Scoreboard
	}{
		{
			scenario: "failed to fetch scoreboard from nba api",
			nbaErr:   errFailed,
			expErr:   fmt.Errorf("failed to get scoreboard: %w", errFailed),
		},
		{
			scenario: "fetch scoreboard from nba api",
			sb:       nba.ScoreboardData{},
			expRes: stats.Scoreboard{
				Games: []stats.Game{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.scenario, func() {
			s.SetupTest()

			var (
				ctx = context.Background()

				cmd = nba.GetScoreboardCommand{
					Date:     "2022-10-01",
					LeagueID: nba.NBA,
				}
			)

			s.nm.On("GetScoreboard", ctx, cmd).Return(tt.sb, tt.nbaErr)

			res, err := s.s.GetScoreboard(ctx, cmd)

			s.Equal(tt.expErr, err)
			s.Equal(tt.expRes, res)
		})
	}
}

func (s *ServiceTestSuite) TestGetBoxscore() {
	tests := []struct {
		scenario string

		nbaData nba.BoxscoreData
		nbaErr  error

		expErr error
		expRes stats.Boxscore
	}{
		{
			scenario: "failed to fetch boxscore from nba api",
			nbaErr:   errFailed,
			expErr:   fmt.Errorf("failed to get boxscore: %w", errFailed),
		},
		{
			scenario: "fetch boxscore from nba api",
			nbaData:  nba.BoxscoreData{},
			expRes: stats.Boxscore{
				HomeTeam: stats.Team{
					Stats: stats.Stats{
						Minutes: "0:00",
					},
					Players: []stats.Player{},
				},
				AwayTeam: stats.Team{
					Stats: stats.Stats{
						Minutes: "0:00",
					},
					Players: []stats.Player{},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.scenario, func() {
			s.SetupTest()

			var (
				ctx = context.Background()

				cmd = nba.GetBoxscoreCommand{
					GameID: "005123512",
				}
			)

			s.nm.On("GetBoxscore", ctx, cmd).Return(tt.nbaData, tt.nbaErr)

			res, err := s.s.GetBoxscore(ctx, cmd)

			s.Equal(tt.expErr, err)
			s.Equal(tt.expRes, res)
		})
	}
}
