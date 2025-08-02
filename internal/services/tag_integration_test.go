//go:build integration
// +build integration

package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/services"
	"github.com/garnizeh/englog/internal/testutils"
)

// TestTagServiceIntegration tests the full tag management flow with database
// "Tags are the labels that bring order to chaos." 🏷️
func TestTagServiceIntegration(t *testing.T) {
	db := testutils.DB(t)

	// Create test logger
	testLogger := logging.NewTestLogger()

	tagService := services.NewTagService(db, testLogger)

	t.Run("FullTagLifecycle", func(t *testing.T) {
		ctx := context.Background()

		// Create a new tag
		createReq := &models.TagRequest{
			Name:        fmt.Sprintf("integration-tag-%d", time.Now().UnixNano()),
			Color:       "#FF5733",
			Description: stringPtr("Integration test tag"),
		}

		// Test tag creation
		createdTag, err := tagService.CreateTag(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, createdTag)
		assert.Equal(t, createReq.Name, createdTag.Name)
		assert.Equal(t, createReq.Color, createdTag.Color)
		assert.Equal(t, *createReq.Description, *createdTag.Description)
		assert.NotEqual(t, uuid.Nil, createdTag.ID)

		// Test tag retrieval by ID
		retrievedTag, err := tagService.GetTag(ctx, createdTag.ID.String())
		require.NoError(t, err)
		assert.Equal(t, createdTag.ID, retrievedTag.ID)
		assert.Equal(t, createdTag.Name, retrievedTag.Name)
		assert.Equal(t, createdTag.Color, retrievedTag.Color)

		// Test tag retrieval by name
		tagByName, err := tagService.GetTagByName(ctx, createdTag.Name)
		require.NoError(t, err)
		assert.Equal(t, createdTag.ID, tagByName.ID)

		// Test tag update
		updateReq := &models.TagRequest{
			Name:        createdTag.Name + "-updated",
			Color:       "#00FF00",
			Description: stringPtr("Updated description"),
		}

		updatedTag, err := tagService.UpdateTag(ctx, createdTag.ID.String(), updateReq)
		require.NoError(t, err)
		assert.Equal(t, updateReq.Name, updatedTag.Name)
		assert.Equal(t, updateReq.Color, updatedTag.Color)
		assert.Equal(t, *updateReq.Description, *updatedTag.Description)

		// Test tag deletion
		err = tagService.DeleteTag(ctx, createdTag.ID.String())
		require.NoError(t, err)

		// Verify tag is deleted
		_, err = tagService.GetTag(ctx, createdTag.ID.String())
		assert.Error(t, err)
	})

	t.Run("TagSearchAndFiltering", func(t *testing.T) {
		ctx := context.Background()

		// Create multiple tags for testing
		testTags := []*models.TagRequest{
			{
				Name:        fmt.Sprintf("search-dev-%d", time.Now().UnixNano()),
				Color:       "#FF0000",
				Description: stringPtr("Development tag"),
			},
			{
				Name:        fmt.Sprintf("search-test-%d", time.Now().UnixNano()),
				Color:       "#00FF00",
				Description: stringPtr("Testing tag"),
			},
			{
				Name:        fmt.Sprintf("search-bug-%d", time.Now().UnixNano()),
				Color:       "#0000FF",
				Description: stringPtr("Bug tracking tag"),
			},
		}

		var createdTagIDs []string
		for _, tagReq := range testTags {
			tag, err := tagService.CreateTag(ctx, tagReq)
			require.NoError(t, err)
			createdTagIDs = append(createdTagIDs, tag.ID.String())
		}

		// Clean up created tags
		defer func() {
			for _, tagID := range createdTagIDs {
				_ = tagService.DeleteTag(ctx, tagID)
			}
		}()

		// Test GetAllTags
		allTags, err := tagService.GetAllTags(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allTags), 3)

		// Test SearchTags
		searchResults, err := tagService.SearchTags(ctx, "search-", 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(searchResults), 3)

		// Verify search results contain our created tags
		foundNames := make(map[string]bool)
		for _, tag := range searchResults {
			foundNames[tag.Name] = true
		}

		for _, originalTag := range testTags {
			assert.True(t, foundNames[originalTag.Name], "Search should find tag: %s", originalTag.Name)
		}
	})

	t.Run("ConcurrentTagOperations", func(t *testing.T) {
		ctx := context.Background()

		// Test concurrent tag creation
		const numGoroutines = 10
		results := make(chan *models.Tag, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				tagReq := &models.TagRequest{
					Name:        fmt.Sprintf("concurrent-tag-%d-%d", index, time.Now().UnixNano()),
					Color:       "#FF0000",
					Description: stringPtr(fmt.Sprintf("Concurrent tag %d", index)),
				}

				tag, err := tagService.CreateTag(ctx, tagReq)
				if err != nil {
					errors <- err
					return
				}
				results <- tag
			}(i)
		}

		// Collect results
		var createdTags []*models.Tag
		for i := 0; i < numGoroutines; i++ {
			select {
			case tag := <-results:
				createdTags = append(createdTags, tag)
			case err := <-errors:
				t.Errorf("Concurrent operation failed: %v", err)
			case <-time.After(10 * time.Second):
				t.Error("Timeout waiting for concurrent operations")
			}
		}

		// Clean up created tags
		defer func() {
			for _, tag := range createdTags {
				_ = tagService.DeleteTag(ctx, tag.ID.String())
			}
		}()

		assert.Equal(t, numGoroutines, len(createdTags))

		// Verify all tags are unique
		tagNames := make(map[string]bool)
		for _, tag := range createdTags {
			assert.False(t, tagNames[tag.Name], "Tag name should be unique: %s", tag.Name)
			tagNames[tag.Name] = true
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		ctx := context.Background()

		// Test creating duplicate tags
		tagReq := &models.TagRequest{
			Name:        fmt.Sprintf("duplicate-tag-%d", time.Now().UnixNano()),
			Color:       "#FF0000",
			Description: stringPtr("Original tag"),
		}

		tag1, err := tagService.CreateTag(ctx, tagReq)
		require.NoError(t, err)

		defer func() {
			_ = tagService.DeleteTag(ctx, tag1.ID.String())
		}()

		// Try to create duplicate
		_, err = tagService.CreateTag(ctx, tagReq)
		assert.Error(t, err)

		// Test getting non-existent tag
		_, err = tagService.GetTag(ctx, uuid.New().String())
		assert.Error(t, err)

		// Test updating non-existent tag
		updateReq := &models.TagRequest{
			Name:        "updated-name",
			Color:       "#00FF00",
			Description: stringPtr("Updated"),
		}
		_, err = tagService.UpdateTag(ctx, uuid.New().String(), updateReq)
		assert.Error(t, err)

		// Test deleting non-existent tag
		err = tagService.DeleteTag(ctx, uuid.New().String())
		assert.Error(t, err)

		// Test invalid UUID formats
		invalidUUIDs := []string{"", "invalid", "not-a-uuid", "12345"}
		for _, invalidUUID := range invalidUUIDs {
			_, err = tagService.GetTag(ctx, invalidUUID)
			assert.Error(t, err)
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		ctx := context.Background()

		// Test with minimal tag data
		minimalTag := &models.TagRequest{
			Name:        fmt.Sprintf("minimal-%d", time.Now().UnixNano()),
			Color:       "#000000",
			Description: nil,
		}

		tag, err := tagService.CreateTag(ctx, minimalTag)
		require.NoError(t, err)
		assert.Nil(t, tag.Description)

		defer func() {
			_ = tagService.DeleteTag(ctx, tag.ID.String())
		}()

		// Test with maximum length fields
		maxTag := &models.TagRequest{
			Name:        fmt.Sprintf("max-length-tag-name-with-many-chars-%d", time.Now().UnixNano()),
			Color:       "#FFFFFF",
			Description: stringPtr("This is a very long description that tests the maximum length handling for tag descriptions in the database and service layer"),
		}

		tag2, err := tagService.CreateTag(ctx, maxTag)
		require.NoError(t, err)

		defer func() {
			_ = tagService.DeleteTag(ctx, tag2.ID.String())
		}()

		// Test search with empty query
		emptyResults, err := tagService.SearchTags(ctx, "", 10)
		require.NoError(t, err)
		assert.NotNil(t, emptyResults)

		// Test search with very large limit
		largeResults, err := tagService.SearchTags(ctx, "tag", 1000)
		require.NoError(t, err)
		assert.NotNil(t, largeResults)
	})
}

// BenchmarkTagService benchmarks tag service operations
func BenchmarkTagService(b *testing.B) {
	db := testutils.DB(b)

	// Create test logger
	testLogger := logging.NewTestLogger()

	tagService := services.NewTagService(db, testLogger)
	ctx := context.Background()

	b.Run("CreateTag", func(b *testing.B) {
		var createdIDs []string

		b.ResetTimer()
		for i := range b.N {
			req := &models.TagRequest{
				Name:        fmt.Sprintf("bench-tag-%d", i),
				Color:       "#FF0000",
				Description: stringPtr("Benchmark tag"),
			}

			tag, err := tagService.CreateTag(ctx, req)
			if err != nil {
				b.Fatalf("Failed to create tag: %v", err)
			}
			createdIDs = append(createdIDs, tag.ID.String())
		}
		b.StopTimer()

		// Cleanup
		for _, id := range createdIDs {
			_ = tagService.DeleteTag(ctx, id)
		}
	})

	b.Run("GetTag", func(b *testing.B) {
		// Setup: create a tag to fetch
		req := &models.TagRequest{
			Name:        "bench-get-tag",
			Color:       "#FF0000",
			Description: stringPtr("Benchmark get tag"),
		}

		tag, err := tagService.CreateTag(ctx, req)
		if err != nil {
			b.Fatalf("Failed to create test tag: %v", err)
		}

		defer func() {
			_ = tagService.DeleteTag(ctx, tag.ID.String())
		}()

		b.ResetTimer()
		for range b.N {
			_, err := tagService.GetTag(ctx, tag.ID.String())
			if err != nil {
				b.Fatalf("Failed to get tag: %v", err)
			}
		}
	})

	b.Run("SearchTags", func(b *testing.B) {
		// Setup: create some tags to search
		for i := 0; i < 10; i++ {
			req := &models.TagRequest{
				Name:        fmt.Sprintf("search-bench-tag-%d", i),
				Color:       "#FF0000",
				Description: stringPtr("Search benchmark tag"),
			}
			_, _ = tagService.CreateTag(ctx, req)
		}

		b.ResetTimer()
		for range b.N {
			_, err := tagService.SearchTags(ctx, "search-bench", 10)
			if err != nil {
				b.Fatalf("Failed to search tags: %v", err)
			}
		}
	})
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
