package db_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestFrom(t *testing.T) {
	d := testutils.GetDb(t)
	d.Close()
}

func TestPlayerCapitilisationFix1(t *testing.T) {
	name := "danny piper"
	require.Equal(t, "Danny Piper", db.InternalFixPlayerNameCapitals(name))
}

func TestPlayerCapitilisationFix2(t *testing.T) {
	name := "rhys"
	require.Equal(t, "Rhys", db.InternalFixPlayerNameCapitals(name))
}

func TestMigrationsOnPrototypeDatabase(t *testing.T) {
	t.Parallel()

	const (
		port          = 5444
		containerName = "chessleagueoldmigrationtest"
		user          = "test"
		password      = "test"
		dbName        = "test"
	)

	dockerKill := exec.CommandContext(t.Context(), "docker", "kill", containerName)
	dockerKill.Stdout = os.Stdout
	dockerKill.Stderr = os.Stderr
	t.Logf("Executing docker kill: %s", dockerKill.Run())

	dockerRm := exec.CommandContext(t.Context(), "docker", "container", "rm", containerName)
	dockerRm.Stdout = os.Stdout
	dockerRm.Stderr = os.Stderr
	t.Logf("Exexcuting docker container rm: %s", dockerRm.Run())

	t.Logf("Starting test database")
	cmd := exec.CommandContext(t.Context(), "docker", "run",
		"--name", containerName,
		"-e", fmt.Sprintf("POSTGRES_USER=%s", user),
		"-e", fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
		"-e", fmt.Sprintf("POSTGRES_DB=%s", dbName),
		"-p", fmt.Sprintf("%d:5432", port),
		"postgres")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	require.NoError(t, err, "Cannot start test database")

	t.Cleanup(func() {
		exec.Command("docker", "stop", containerName).Run()
	})

	defer cmd.Process.Kill()

	var sqlDb *sqlx.DB
	const (
		maxTries = 30
		wait     = time.Second / 2
	)

	for tries := range maxTries {
		tries++
		pgDb, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=localhost port=%d sslmode=disable",
			user,
			password,
			dbName,
			port))
		if err != nil {
			t.Log("Could not connect to database trying again...", "waiting for", wait, "err", err)
			time.Sleep(wait)
		} else {
			sqlDb = pgDb
			break
		}
	}

	require.NotNil(t, sqlDb)

	bytes, err := os.ReadFile("./nagshead.sql")
	require.NoError(t, err)

	_, err = sqlDb.Exec(string(bytes))
	require.NoError(t, err, "Must execute datbase export")

	database, err := db.From(sqlDb)
	require.NoError(t, err)
	defer database.Close()
}
