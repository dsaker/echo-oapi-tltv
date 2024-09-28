package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	newUser := InsertUserParams{
		Name:           util.RandomName(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
		TitleID:        util.ValidTitleId,
		Flipped:        false,
		OgLanguageID:   util.ValidOgLanguageId,
		NewLanguageID:  util.ValidNewLanguageId,
	}

	user, err := testQueries.InsertUser(context.Background(), newUser)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, newUser.Name, user.Name)
	require.Equal(t, newUser.Email, user.Email)
	require.Equal(t, newUser.Flipped, user.Flipped)
	require.Equal(t, newUser.TitleID, user.TitleID)
	require.Equal(t, newUser.OgLanguageID, user.OgLanguageID)
	require.Equal(t, newUser.NewLanguageID, user.NewLanguageID)
	require.NotZero(t, user.Created)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
