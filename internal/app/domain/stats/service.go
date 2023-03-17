package stats

import (
	"context"
	"fmt"

	"github.com/pedro-mealha/nba-stats-api/internal/app/gateway/nba"
)

type (
	Provider interface {
		GetScoreboard(context.Context, nba.GetScoreboardCommand) (Scoreboard, error)
		GetBoxscore(context.Context, nba.GetBoxscoreCommand) (Boxscore, error)
	}

	Service struct {
		a nba.API
	}
)

func NewService(a nba.API) *Service { return &Service{a} }

func (s *Service) GetScoreboard(ctx context.Context, cmd nba.GetScoreboardCommand) (Scoreboard, error) {
	sb, err := s.a.GetScoreboard(ctx, cmd)
	if err != nil {
		return Scoreboard{}, fmt.Errorf("failed to get scoreboard: %w", err)
	}

	return NewScoreboard(sb), nil
}

func (s *Service) GetBoxscore(ctx context.Context, cmd nba.GetBoxscoreCommand) (Boxscore, error) {
	bs, err := s.a.GetBoxscore(ctx, cmd)
	if err != nil {
		return Boxscore{}, fmt.Errorf("failed to get boxscore: %w", err)
	}

	return NewBoxscore(bs), nil
}
