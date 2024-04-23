package test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/docker/docker/api/types/filters"
	"github.com/lib/pq"
	"github.com/skraio/banner-service/internal/data"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

type BannerTestSuite struct {
	suite.Suite
	ctx context.Context
	data.BannerModel
	pgContainer        *postgres.PostgresContainer
	pgConnectionString string
}

func (suite *BannerTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := postgres.RunContainer(
		suite.ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithDatabase("bannerservice"),
		postgres.WithUsername("bannerservice"),
		postgres.WithPassword("pass777"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	suite.NoError(err)

	connStr, err := pgContainer.ConnectionString(suite.ctx)
	suite.NoError(err)

	db, err := sql.Open("postgres", "postgres://bannerservice:pass777@localhost/bannerservice?sslmode=disable")
	suite.NoError(err)

	suite.BannerModel.DB = db
	suite.pgContainer = pgContainer
	suite.pgConnectionString = connStr

	err = suite.BannerModel.DB.Ping()
	suite.NoError(err)
}

func (suite *BannerTestSuite) TearDownSuite() {
	err := suite.pgContainer.Terminate(suite.ctx)
	suite.NoError(err)
}

func (suite *BannerTestSuite) SetupTest() {
	suite.BannerModel.DB.Exec(`
    create table if not exists banners ( 
        banner_id bigserial primary key,
        tag_ids integer[] not null,
        feature_id integer not null,
        content json not null,
        is_active boolean not null,
        created_at timestamp(3) with time zone not null default now(),
        updated_at timestamp(3) with time zone not null default now()
    );`)
}

func (suite *BannerTestSuite) TearDownTest() {
	suite.BannerModel.DB.Exec("DROP TABLE IF EXISTS banners CASCADE;")
}

func (suite *BannerTestSuite) BeforeTest(_ string, testName string) {
	if testName == "TestShowBanner" || testName == "TestUpdateBanner" || testName == "TestDeleteBanner" {
		newBanner := &data.Banner{FeatureID: 123, TagIDs: []int64{1, 2, 3}, IsActive: true, Content: data.Content{Title: "ABCD", Text: "WQER", URL: "https://exm.com"}}
		err := suite.BannerModel.Insert(newBanner)
		suite.NoError(err)
	}
}

func (suite *BannerTestSuite) TestCreateBanner() {
	var count int
	err := suite.BannerModel.DB.QueryRow("SELECT COUNT(*) FROM banners").Scan(&count)
	suite.NoError(err)
	suite.Equal(0, count, "There should not be any records in the database")

	newBanner := &data.Banner{
		FeatureID: 123,
		TagIDs:    []int64{1, 2, 3},
		IsActive:  true,
		Content: data.Content{
			Title: "ABCD",
			Text:  "WQER",
			URL:   "https://exm.com",
		},
	}
	err = suite.BannerModel.Insert(newBanner)
	suite.NoError(err)

	err = suite.BannerModel.DB.QueryRow("SELECT COUNT(*) FROM banners").Scan(&count)
	suite.NoError(err)
	suite.Equal(1, count, "There should be exactly one record in the database")

	var banner data.Banner
	err = suite.BannerModel.DB.QueryRow("SELECT feature_id, tag_ids, is_active, content, created_at, updated_at FROM banners").Scan(
		&banner.FeatureID,
		pq.Array(&banner.TagIDs),
		&banner.IsActive,
		&banner.Content,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)
	suite.NoError(err)

	suite.Equal(newBanner.FeatureID, banner.FeatureID, "FeatureID should match the inserted value")
	suite.ElementsMatch(newBanner.TagIDs, banner.TagIDs, "TagIDs should match the inserted values")
	suite.Equal(newBanner.IsActive, banner.IsActive, "IsActive status should match")
	suite.Equal(newBanner.Content.Title, banner.Content.Title, "Content Title should match")
	suite.Equal(newBanner.Content.Text, banner.Content.Text, "Content Text should match")
	suite.Equal(newBanner.Content.URL, banner.Content.URL, "Content URL should match")

	suite.WithinDuration(newBanner.CreatedAt, banner.CreatedAt, time.Second, "CreatedAt should be within reasonable range")
	suite.WithinDuration(newBanner.UpdatedAt, banner.UpdatedAt, time.Second, "UpdatedAt should be within reasonable range")
}

func (suite *BannerTestSuite) TestShowBanner() {
	var count int
	err := suite.BannerModel.DB.QueryRow("SELECT COUNT(*) FROM banners").Scan(&count)
	suite.NoError(err)
	suite.Equal(1, count, "There should be exactly one record in the database")

	var filters data.UserFilters
	filters.TagID = 2
	filters.FeatureID = 123
	filters.UseLastRevision = false
	role := data.RoleUser

	originalBanner, err := suite.BannerModel.Get(filters, role)
	suite.NoError(err)

	var banner data.Banner
	err = suite.BannerModel.DB.QueryRow("SELECT banner_id, tag_ids, is_active, content, created_at, updated_at FROM banners").Scan(
		&banner.BannerID,
		pq.Array(&banner.TagIDs),
		&banner.IsActive,
		&banner.Content,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)
	suite.NoError(err)

	suite.Equal(originalBanner.BannerID, banner.BannerID, "BannerID should match the inserted value")
	suite.Equal(originalBanner.IsActive, banner.IsActive, "IsActive status should match")
	suite.Equal(originalBanner.Content.Title, banner.Content.Title, "Content Title should match")
	suite.Equal(originalBanner.Content.Text, banner.Content.Text, "Content Text should match")
	suite.Equal(originalBanner.Content.URL, banner.Content.URL, "Content URL should match")

	suite.WithinDuration(originalBanner.CreatedAt, banner.CreatedAt, time.Second, "CreatedAt should be within reasonable range")
	suite.WithinDuration(originalBanner.UpdatedAt, banner.UpdatedAt, time.Second, "UpdatedAt should be within reasonable range")
}

func (suite *BannerTestSuite) TestUpdateBanner() {
	var count int
	err := suite.BannerModel.DB.QueryRow("SELECT COUNT(*) FROM banners").Scan(&count)
	suite.NoError(err)
	suite.Equal(1, count, "There should be exactly one record in the database")

	originalBanner, err := suite.BannerModel.GetByID(1)
	suite.NoError(err)

	newBannerFields := &data.Banner{
		BannerID:  originalBanner.BannerID,
		TagIDs:    originalBanner.TagIDs,
		FeatureID: 9876,
		Content: data.Content{
			Title: "New Title",
			Text:  "New Text",
			URL:   originalBanner.Content.URL,
		},
		IsActive:  originalBanner.IsActive,
		UpdatedAt: originalBanner.UpdatedAt,
	}

	err = suite.BannerModel.Update(newBannerFields)
	suite.NoError(err)

	err = suite.BannerModel.DB.QueryRow("SELECT COUNT(*) FROM banners").Scan(&count)
	suite.NoError(err)
	suite.Equal(1, count, "There should be exactly one record in the database")

	var updatedBanner data.Banner
	err = suite.BannerModel.DB.QueryRow("SELECT feature_id, tag_ids, is_active, content, created_at, updated_at FROM banners").Scan(
		&updatedBanner.FeatureID,
		pq.Array(&updatedBanner.TagIDs),
		&updatedBanner.IsActive,
		&updatedBanner.Content,
		&updatedBanner.CreatedAt,
		&updatedBanner.UpdatedAt,
	)
	suite.NoError(err)

	suite.Equal(newBannerFields.FeatureID, updatedBanner.FeatureID, "FeatureID should match the inserted value")
	suite.Equal(originalBanner.IsActive, updatedBanner.IsActive, "IsActive status should match")
	suite.Equal(newBannerFields.Content.Title, updatedBanner.Content.Title, "Content Title should match")
	suite.Equal(newBannerFields.Content.Text, updatedBanner.Content.Text, "Content Text should match")
	suite.Equal(originalBanner.Content.URL, updatedBanner.Content.URL, "Content URL should match")

	suite.WithinDuration(originalBanner.CreatedAt, updatedBanner.CreatedAt, time.Second, "CreatedAt should be within reasonable range")
}

func (suite *BannerTestSuite) TestDeleteBanner() {
	var count int
	err := suite.BannerModel.DB.QueryRow("SELECT COUNT(*) FROM banners").Scan(&count)
	suite.NoError(err)
	suite.Equal(1, count, "There should be exactly one record in the database")

	err = suite.BannerModel.Delete(1)
	suite.NoError(err)

	err = suite.BannerModel.DB.QueryRow("SELECT COUNT(*) FROM banners").Scan(&count)
	suite.NoError(err)
	suite.Equal(0, count, "There should not be any records in the database")
}

func TestBannerDatabaseSuite(t *testing.T) {
	suite.Run(t, new(BannerTestSuite))
}
