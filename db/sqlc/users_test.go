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
		OgLanguageID:   util.ValidOgLanguageId,
		NewLanguageID:  util.ValidNewLanguageId,
	}

	user, err := testQueries.InsertUser(context.Background(), newUser)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, newUser.Name, user.Name)
	require.Equal(t, newUser.Email, user.Email)
	require.Equal(t, newUser.TitleID, user.TitleID)
	require.Equal(t, newUser.OgLanguageID, user.OgLanguageID)
	require.Equal(t, newUser.NewLanguageID, user.NewLanguageID)
	require.NotZero(t, user.Created)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestSelectUserById(t *testing.T) {
	user := createRandomUser(t)

	newUser, err := testQueries.SelectUserById(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, newUser.Name, user.Name)
	require.Equal(t, newUser.Email, user.Email)
	require.Equal(t, newUser.TitleID, user.TitleID)
	require.Equal(t, newUser.OgLanguageID, user.OgLanguageID)
	require.Equal(t, newUser.NewLanguageID, user.NewLanguageID)
	require.NotZero(t, user.Created)

	_, err = testQueries.SelectUserById(context.Background(), util.InvalidUserId)

	require.Error(t, err)
	require.Contains(t, err.Error(), "sql: no rows in result set")
}

func TestSelectUserByName(t *testing.T) {
	user := createRandomUser(t)

	newUser, err := testQueries.SelectUserByName(context.Background(), user.Name)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, newUser.Name, user.Name)
	require.Equal(t, newUser.Email, user.Email)
	require.Equal(t, newUser.TitleID, user.TitleID)
	require.Equal(t, newUser.OgLanguageID, user.OgLanguageID)
	require.Equal(t, newUser.NewLanguageID, user.NewLanguageID)
	require.NotZero(t, user.Created)

	_, err = testQueries.SelectUserByName(context.Background(), string(rune(util.InvalidUserId)))

	require.Error(t, err)
	require.Contains(t, err.Error(), "sql: no rows in result set")
}

func TestUpdateUserById(t *testing.T) {
	user := createRandomUser(t)

	updateUserByIdParams := UpdateUserByIdParams{
		ID:             user.ID,
		TitleID:        user.TitleID,
		Email:          user.Email,
		OgLanguageID:   user.OgLanguageID,
		NewLanguageID:  user.NewLanguageID,
		HashedPassword: user.HashedPassword,
	}
	newUser, err := testQueries.UpdateUserById(context.Background(), updateUserByIdParams)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, newUser.Name, user.Name)
	require.Equal(t, newUser.Email, user.Email)
	require.Equal(t, newUser.TitleID, user.TitleID)
	require.Equal(t, newUser.OgLanguageID, user.OgLanguageID)
	require.Equal(t, newUser.NewLanguageID, user.NewLanguageID)
	require.NotZero(t, user.Created)

	updateUserByIdParams.ID = util.InvalidUserId
	_, err = testQueries.UpdateUserById(context.Background(), updateUserByIdParams)

	require.Error(t, err)
	require.Contains(t, err.Error(), "sql: no rows in result set")
}
