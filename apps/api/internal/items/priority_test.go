package items_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
)

func TestScores(t *testing.T) {
	gatitos, claude, google := uuid.New(), uuid.New(), uuid.New()
	a, b, c := uuid.New(), uuid.New(), uuid.New()
	tests := []struct {
		name   string
		rivals map[uuid.UUID][]uuid.UUID
		want   map[uuid.UUID]items.Score
	}{
		{
			"all sources agree",
			map[uuid.UUID][]uuid.UUID{a: {gatitos, claude, google}},
			map[uuid.UUID]items.Score{a: {3, items.PriorityQuick}},
		},
		{
			"two of three agree",
			map[uuid.UUID][]uuid.UUID{a: {claude, google}, b: {gatitos}},
			map[uuid.UUID]items.Score{a: {2, items.PriorityStandard}, b: {1, items.PriorityFull}},
		},
		{
			"all disagree",
			map[uuid.UUID][]uuid.UUID{a: {gatitos}, b: {claude}, c: {google}},
			map[uuid.UUID]items.Score{a: {1, items.PriorityFull}, b: {1, items.PriorityFull}, c: {1, items.PriorityFull}},
		},
		{
			"one source alone is not agreement",
			map[uuid.UUID][]uuid.UUID{a: {gatitos}},
			map[uuid.UUID]items.Score{a: {1, items.PriorityFull}},
		},
		{
			"a source backing two answers",
			map[uuid.UUID][]uuid.UUID{a: {gatitos, claude}, b: {claude}},
			map[uuid.UUID]items.Score{a: {2, items.PriorityQuick}, b: {1, items.PriorityFull}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := items.Scores(tt.rivals)
			for id, want := range tt.want {
				if got[id] != want {
					t.Errorf("score = %+v, want %+v", got[id], want)
				}
			}
		})
	}
}
