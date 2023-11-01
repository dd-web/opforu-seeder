package main

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetType string

const (
	AssetTypeImage AssetType = "image"
	AssetTypeVideo AssetType = "video"
)

// unique asset - make references if it already exists
type AssetSource struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`

	Details struct {
		Avatar FileCtx `json:"avatar" bson:"avatar"`
		Source FileCtx `json:"source" bson:"source"`
	}

	AssetType AssetType            `json:"asset_type" bson:"asset_type"`
	Uploaders []primitive.ObjectID `json:"uploaders" bson:"uploaders"`

	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// reference (populated from source when send to client)
type Asset struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`

	SourceID  primitive.ObjectID `json:"source_id" bson:"source_id"`
	AccountID primitive.ObjectID `json:"account_id" bson:"account_id"`

	FileName string   `json:"file_name" bson:"file_name"`
	Tags     []string `json:"tags" bson:"tags"`

	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type FileCtx struct {
	Height    uint16 `json:"height" bson:"height"`
	Width     uint16 `json:"width" bson:"width"`
	Size      uint64 `json:"size" bson:"size"`
	SizeStr   string `json:"size_str" bson:"size_str"`
	URL       string `json:"url" bson:"url"`
	Extension string `json:"extension" bson:"extension"`
}

// generates asset sources to create Assets from (references)
func (s *MongoStore) GenerateAssetSources(min, max int) {
	assetCount := RandomIntBetween(min, max)

	for i := 0; i < assetCount; i++ {
		fmt.Print("\033[G\033[K")
		fmt.Printf(" - Generating Assets: %v/%v", i+1, assetCount)
		assetsrc := GenerateAssetSource(i)
		s.cAssetSrcMap[i] = assetsrc
	}
	fmt.Print("\n")
}

// generates assets for a post or thread and returns their id's
func (s *MongoStore) GenerateAssetCount(count int, creatorId primitive.ObjectID) ([]primitive.ObjectID, error) {
	if count > len(s.cAssetSrcMap) {
		return nil, fmt.Errorf("invalid asset source count %d out of bounds", count)
	}

	ids := []primitive.ObjectID{}
	indexId := RandomIntBetween(9, len(s.cAssetSrcMap)-count)

	for i := 0; i < count; i++ {
		asset, err := s.GenerateAsset(indexId+i, creatorId)
		if err != nil {
			fmt.Printf("Error generating asset for post: %v\n - skipping\n", err)
			continue
		}
		ids = append(ids, asset.ID)
	}

	return ids, nil
}

// creates an asset from the source locaated at the index and returns a pointer to it
func (s *MongoStore) GenerateAsset(index int, creator primitive.ObjectID) (*Asset, error) {
	if index > len(s.cAssetSrcMap) {
		return nil, fmt.Errorf("invalid asset source index %d", index)
	}

	ts := time.Now().UTC()

	assetSource := s.cAssetSrcMap[index]
	assetSource.Uploaders = append(assetSource.Uploaders, creator)

	asset := &Asset{
		ID:        primitive.NewObjectID(),
		SourceID:  assetSource.ID,
		AccountID: creator,
		FileName:  SelectAnyWord(),
		Tags:      GetRandomTags(),
		CreatedAt: &ts,
		UpdatedAt: &ts,
	}

	s.cAssets = append(s.cAssets, asset)

	return asset, nil
}

// saves assets to the db
func (s *MongoStore) PersistAssets() error {
	docs := []interface{}{}
	for _, asset := range s.cAssets {
		docs = append(docs, asset)
	}
	return s.PersistDocuments(docs, "assets")
}

// saves asset sources to the db
func (s *MongoStore) PersistAssetSrc() error {
	docs := []interface{}{}
	for _, assetSrc := range s.cAssetSrcMap {
		docs = append(docs, assetSrc)
	}
	return s.PersistDocuments(docs, "asset_sources")
}
