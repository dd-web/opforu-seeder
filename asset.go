package main

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetChecksum struct {
	MD5    string `json:"md5" bson:"md5"`
	SHA256 string `json:"sha256" bson:"sha256"`
}

// DO NOT inline unless you know what you're doing
type PairedAsset[T comparable] struct {
	Source T `json:"source" bson:"source"`
	Thumb  T `json:"thumb" bson:"thumb"`
}

type MediaSource struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`

	FileSize PairedAsset[int]    `json:"file_size" bson:"file_size"`
	URL      PairedAsset[string] `json:"url" bson:"url"`

	MimeType  string    `json:"mime_type" bson:"mime_type"`
	MediaType string    `json:"media_type" bson:"media_type"` // image or video
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	Checksum AssetChecksum `json:"checksum" bson:"checksum"`

	FileExtension string `json:"file_extension" bson:"file_extension"`
}

type Media struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Source       primitive.ObjectID `json:"source" bson:"source"`
	UserFilename string             `json:"user_filename" bson:"user_filename"`

	URL PairedAsset[string] `json:"url" bson:"url"`

	// below reference the source (not thumbnail)
	FileSize  int       `json:"file_size" bson:"file_size"`
	MediaType string    `json:"media_type" bson:"media_type"` // image or video
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// new empty media source - a media struct musn't be generated without a source reference
func NewEmptyMediaSource() *MediaSource {
	return &MediaSource{}
}

// new empty media
func NewEmptyMedia() *Media {
	return &Media{}
}

// randomize media source values - never use with an existing media source
func (m *MediaSource) Randomize(id int) {
	mType := GetRandomMediaType()
	src, thumb := FormatImageUrls(id)
	ext := GetRandomExtByType(mType)

	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now().UTC()
	m.UpdatedAt = time.Now().UTC()
	m.URL.Source = src
	m.URL.Thumb = thumb
	m.FileExtension = ext
	m.MediaType = mType
	m.MimeType = mType + "/" + ext
}

// Generate Media Sources - Posts will pick from these to create references
func (s *MongoStore) GenerateMediaSources(min, max int) {
	mediaSourceCount := RandomIntBetween(min, max)

	for i := 0; i < mediaSourceCount; i++ {
		fmt.Print("\033[G\033[K")
		fmt.Printf(" - Generating Media Sources: %v/%v", i+1, mediaSourceCount)
		mediaSource := NewEmptyMediaSource()
		mediaSource.Randomize(i)
		s.cMediaSourceMap[i] = mediaSource
	}
	fmt.Print("\n")
}

// generates a list of media ids for a post
func (s *MongoStore) GenerateMediaCount(num int) ([]primitive.ObjectID, error) {
	if num > len(s.cMediaSourceMap) {
		return nil, fmt.Errorf("invalid media source count %d", num)
	}

	mediaIds := []primitive.ObjectID{}
	msIx := RandomIntBetween(9, len(s.cMediaSourceMap)-num)

	for i := 0; i < num; i++ {
		media, err := s.GenerateMediaForPost(msIx + i)
		if err != nil {
			fmt.Printf("Error generating media for post: %v\n - skipping\n", err)
			continue
		}
		mediaIds = append(mediaIds, media.ID)
	}

	return mediaIds, nil
}

// Generate Media from id for a post
func (s *MongoStore) GenerateMediaForPost(id int) (*Media, error) {
	if id > len(s.cMediaSourceMap) {
		return nil, fmt.Errorf("invalid media source id %d", id)
	}

	msrc := s.cMediaSourceMap[id]
	m := NewEmptyMedia()

	m.ID = primitive.NewObjectID()
	m.Source = msrc.ID
	m.URL.Source = msrc.URL.Source
	m.URL.Thumb = msrc.URL.Thumb
	m.MediaType = msrc.MediaType
	m.CreatedAt = time.Now().UTC()
	m.UpdatedAt = time.Now().UTC()

	s.cMedia = append(s.cMedia, m)

	return m, nil
}

// persist media
func (s *MongoStore) PersistMedia() error {
	docs := []interface{}{}
	for _, media := range s.cMedia {
		docs = append(docs, media)
	}
	return s.PersistDocuments(docs, "media")
}

// persist media sources
func (s *MongoStore) PersistMediaSources() error {
	docs := []interface{}{}
	for _, msrc := range s.cMediaSourceMap {
		docs = append(docs, msrc)
	}
	return s.PersistDocuments(docs, "media_sources")
}
