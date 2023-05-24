// Package cloudinary contains the cloudinary SDK
// to upload media to the cloudinary CDN
package cloudinary

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Cloudinary is a struct that contains the cloudinary client
type Cloudinary struct {
	Client *cloudinary.Cloudinary
}

// Init is a function that is used to initialize cloudinary
func (c Cloudinary) Init(projectID, apiKey, apiSecret string) (*cloudinary.Cloudinary, error) {
	client, err := cloudinary.NewFromParams(projectID, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Save the given image to the cloudinary CDN
func (c Cloudinary) Save(path, hash, url string, eager, transformation *string) (string, error) {
	if eager == nil {
		defaultEager := "f_avif|f_jp2|f_webp/fl_awebp"
		eager = &defaultEager
	}

	if transformation == nil {
		defaultTransformation := "f_auto/q_auto"
		transformation = &defaultTransformation
	}

	uploadResult, err := c.Client.Upload.Upload(context.Background(), url, uploader.UploadParams{
		PublicID:       hash,
		Folder:         path,
		Overwrite:      api.Bool(true),
		Eager:          *eager,
		EagerAsync:     api.Bool(true),
		Transformation: *transformation,
	})
	if err != nil {
		return url, err
	}

	return uploadResult.SecureURL, nil
}
