// Package vision analyzes the image for potential harmful content
package vision

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	cloudVision "cloud.google.com/go/vision/apiv1"
	"github.com/VinukaThejana/go-utils/errors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"
)

// InitVision is used to initialize the Google Cloud Vision
// instance to detect explicit image content
func InitVision(googleKey string) (*cloudVision.ImageAnnotatorClient, errors.Status) {
	auth, err := base64.StdEncoding.DecodeString(googleKey)
	if err != nil {
		log.Println("Failed to extract the Google private key")
		log.Println(err)
		return nil, errors.InternalServerError
	}

	client, err := cloudVision.NewImageAnnotatorClient(context.Background(), option.WithCredentialsJSON(auth))
	if err != nil {
		log.Println("Failed to initialize the image anotator client")
		log.Println(err)
		return nil, errors.InternalServerError
	}

	return client, errors.Okay
}

// Vision is the struct that contains
// the image anotator instance
type Vision struct {
	Client      *cloudVision.ImageAnnotatorClient
	redisVision *redis.Client
}

// Detect is used to detect wether the
// image contains exmplicit content
func (v Vision) Detect(fileName string, hash string) (errors.Status, VisonTypes) {
	ctx := context.Background()

	// Check the Redis database to reduce API
	// calls
	state := v.redisVision.Get(ctx, hash).Val()
	if state != "" {
		if state != ProperContent.String() {
			return errors.CloudVisionFailed, ParseVsionTypes(state)
		}

		return errors.Okay, ProperContent
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Failed to open the file in cloud vision detector")
		log.Println(err)
		return errors.InternalServerError, UnknwonContent
	}
	defer file.Close()

	image, err := cloudVision.NewImageFromReader(file)
	if err != nil {
		log.Println("Cloud vision failed to read the image")
		log.Println(err)
		return errors.InternalServerError, UnknwonContent
	}

	props, err := v.Client.DetectSafeSearch(context.Background(), image, nil)
	if err != nil {
		log.Println("Failed to analyze by  cloud vision")
		log.Println(err)
		return errors.InternalServerError, UnknwonContent
	}

	adult := props.Adult.Enum().String()
	violence := props.Violence.Enum().String()
	spoof := props.Spoof.Enum().String()
	medical := props.Medical.Enum().String()
	racy := props.Racy.Enum().String()

	if adult == VeryLikely.String() || adult == Likely.String() || adult == Possible.String() {
		log.Println("Cloud vision failed on Adult content")
		return errors.Okay, AdultContent
	}
	if violence == VeryLikely.String() {
		log.Println("Cloud vision failed on violence content")
		return errors.Okay, ViolenceContent
	}

	if spoof == VeryLikely.String() {
		log.Println("Cloud vision failed on spoof content")
		return errors.CloudVisionFailed, SpoofContent
	}

	if medical == VeryLikely.String() {
		log.Println("Cloud vision failed on mediacal content")
		return errors.CloudVisionFailed, MedicalContent
	}

	if racy == VeryLikely.String() {
		log.Println("Cloud vision failed on racy content")
		return errors.CloudVisionFailed, RacyContent
	}

	v.redisVision.Set(ctx, hash, ProperContent.String(), 0)
	return errors.Okay, ProperContent
}

// VisonTypes - The type returned by the Google Cloud vison API
type VisonTypes string

const (
	// ProperContent is to represent proper content
	ProperContent VisonTypes = "PROPER_CONTENT"
	// AdultContent is to represent adult content
	AdultContent VisonTypes = "ADULT_CONTENT"
	// SpoofContent is to represent spoof content
	SpoofContent VisonTypes = "SPOOF_CONTENT"
	// MedicalContent is to represent medical content
	MedicalContent VisonTypes = "MEDICAL_CONTENT"
	// RacyContent is to represent racy content
	RacyContent VisonTypes = "RACY_CONTENT"
	// ViolenceContent is to represent violence content
	ViolenceContent VisonTypes = "VIOLENCE_CONTENT"
	// UnknwonContent is to represent unknown content
	UnknwonContent VisonTypes = "UNKNOWN_CONTENT"
)

// String convert the given vision type enum
// to a string
func (v VisonTypes) String() string {
	return string(v)
}

// ParseVsionTypes parse the string format of the enum
// to the relevant enum
func ParseVsionTypes(v string) VisonTypes {
	vision := map[VisonTypes]struct{}{
		ProperContent:   {},
		AdultContent:    {},
		MedicalContent:  {},
		RacyContent:     {},
		ViolenceContent: {},
		UnknwonContent:  {},
	}

	vis := VisonTypes(v)
	_, ok := vision[vis]
	if !ok {
		return UnknwonContent
	}

	return vis
}

// Sevearity is an Enum to represent the vision
// sevearity
type Sevearity string

// String convert the VisonSevearity enum to its string
// representative
func (vs Sevearity) String() string {
	return string(vs)
}

const (
	// VeryLikely is to represent that the image has a high risk
	// in the given category
	VeryLikely Sevearity = "VERY_LIKELY"
	// Likely is to represent that the image has a medium risk
	// in the given category
	Likely Sevearity = "LIKELY"
	// Possible is to represent that the image has a low risk
	// in the given category
	Possible Sevearity = "POSSIBLE"
)
