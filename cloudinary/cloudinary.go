// Package cloudinary contains the connection to the cloudinary SDK
// and provides various storage functionalites on Cloudinary
package cloudinary

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Cloudinary holds the connection client to cloudinary
type Cloudinary struct {
	Client *cloudinary.Cloudinary
}

// Init is used to initialzie the connection to cloudinary
func (c Cloudinary) Init(projectID, apiKey, apiSecret string) (client *cloudinary.Cloudinary, err error) {
	client, err = cloudinary.NewFromParams(projectID, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Save the given image to the cloudinary CDN
func (c Cloudinary) Save(path, hash, url string, eager, transformation *string) (secureURL string, err error) {
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
