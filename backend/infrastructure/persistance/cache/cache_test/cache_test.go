package cache_test

import (
	"context"
	"testing"

	"github.com/XsedoX/RoomPlay/infrastructure/persistance/cache"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_configuration"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	tests_initializer.InitializeDatabaseContainer()
	tests_initializer.RunTestsWithDatabase(m)
}

func setupMocks(t *testing.T) (
	*sqlx.Tx,
	context.Context,
) {
	txx,
		ctx := tests_initializer.GetTxxAndCtx(t, false)

	return txx, ctx
}

func TestCacheSetAndGet(t *testing.T) {
	txx,
		ctx := setupMocks(t)

	config := mock_configuration.MockConfiguration{}
	type testData struct {
		Name  string
		Value int
	}

	testDataSlice := []testData{
		{Name: "Test1", Value: 1},
		{Name: "Test2", Value: 2},
	}

	cache := cache.NewCache[[]testData](config.CacheSimilarityThreshold())

	setErr := cache.Set("test_key", testDataSlice, ctx, txx)
	require.NoError(t, setErr)

	retrievedData, getErr := cache.Get("test_key", ctx, txx)
	require.NoError(t, getErr)

	require.Equal(t, testDataSlice, retrievedData)
}
