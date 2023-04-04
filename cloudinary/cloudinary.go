// Package cloudinary contains the cloudinary SDK
// to upload media to the cloudinary CDN
package cloudinary

import (
	"context"
	"log"

	"github.com/VinukaThejana/go-utils/errors"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// InitCloudinary is a function that is used to initalize the
// cloudinary config
func InitCloudinary(projectID, apiKey, apiSecret string) (*cloudinary.Cloudinary, errors.Status) {
	client, err := cloudinary.NewFromParams(projectID, apiKey, apiSecret)
	if err != nil {
		log.Println("Failed to initialize cloudinary")
		log.Println(err)
		return nil, errors.InternalServerError
	}

	return client, errors.Okay
}

// Cloudinary is a struct that is used to manage SDKs related
// to cloudinary
type Cloudinary struct {
	Client *cloudinary.Cloudinary
}

// Save the given image to the cloudinary CDN
func (c Cloudinary) Save(path, hash, url string, eager, transformation *string) string {
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
		log.Println("Failed to upload the image to clodinary !")
		log.Println("The url is : ", url, "\nThe hash is : ", hash)
		log.Println(err)
		return url
	}

	return uploadResult.SecureURL
}
