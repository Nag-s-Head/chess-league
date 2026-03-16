package githubapi_test

import (
	"testing"

	githubapi "github.com/Nag-s-Head/chess-league/handlers/utils/github_api"
	"github.com/stretchr/testify/require"
)

func TestGetOrganisation(t *testing.T) {
	t.Run("Lots of members", func(t *testing.T) {
		members, err := githubapi.GerOrganisationMembers("EpicGames")
		require.NoError(t, err)
		require.NotEmpty(t, members)
	})

	t.Run("Some members", func(t *testing.T) {
		members, err := githubapi.GerOrganisationMembers("SquireTournamentServices")
		require.NoError(t, err)
		require.NotEmpty(t, members)
	})

	t.Run("Our org", func(t *testing.T) {
		members, err := githubapi.GerOrganisationMembers("Nag-s-Head")
		require.NoError(t, err)
		require.NotEmpty(t, members)
	})
}

func TestGetUser(t *testing.T) {
	user, err := githubapi.GetUser("djpiper28")
	require.NoError(t, err)
	require.NotEmpty(t, user)
}
