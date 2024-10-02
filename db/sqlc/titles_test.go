package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

func createRandomTitle(t *testing.T) Title {

	args := InsertTitleParams{
		Title:        util.RandomString(8),
		NumSubs:      util.RandomInt32(),
		OgLanguageID: util.ValidOgLanguageId,
	}

	title, err := testQueries.InsertTitle(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, title)

	require.Equal(t, args.Title, title.Title)
	require.Equal(t, args.NumSubs, title.NumSubs)

	require.NotZero(t, title.ID)

	return title
}

func TestAddTitle(t *testing.T) {
	createRandomTitle(t)
}
