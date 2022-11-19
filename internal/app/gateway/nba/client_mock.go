package nba

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// APIMock mock
type APIMock struct{ mock.Mock }

// GetScoreboard mock
func (m *APIMock) GetScoreboard(ctx context.Context, cmd GetScoreboardCommand) (ScoreboardData, error) {
	args := m.Called(ctx, cmd)

	return args.Get(0).(ScoreboardData), args.Error(1)
}

// GetBoxscore mock
func (m *APIMock) GetBoxscore(ctx context.Context, cmd GetBoxscoreCommand) (BoxscoreData, error) {
	args := m.Called(ctx, cmd)

	return args.Get(0).(BoxscoreData), args.Error(1)
}
