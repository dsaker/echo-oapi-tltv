package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"math"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

func createRandomTitle(t *testing.T) Title {

	args := InsertTitleParams{
		Title:      util.RandomString(8),
		NumSubs:    util.RandomInt32(100, 1000),
		LanguageID: util.RandomInt64(math.MinInt64, math.MaxInt64),
	}

	title, err := testQueries.InsertTitle(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, title)

	require.Equal(t, args.Title, title.Title)
	require.Equal(t, args.NumSubs, title.NumSubs)
	require.Equal(t, args.LanguageID, title.LanguageID)

	require.NotZero(t, title.ID)

	return title
}

func TestAddTitle(t *testing.T) {
	createRandomTitle(t)
}
