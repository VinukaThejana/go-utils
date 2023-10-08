// Package vision is used to identify harmful content in images with Google Cloud Platform
package vision

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	cloudVision "cloud.google.com/go/vision/apiv1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"
)

// Vision struct contains the connection to the vision client
type Vision struct {
	Client *cloudVision.ImageAnnotatorClient
	Redis  *redis.Client
}

// Init is used to initialize the connection with the vision SDK
func (v *Vision) Init(googleKey string) error {
	auth, err := base64.StdEncoding.DecodeString(googleKey)
	if err != nil {
		return err
	}

	client, err := cloudVision.NewImageAnnotatorClient(context.Background(), option.WithCredentialsJSON(auth))
	if err != nil {
		return err
	}

	v.Client = client
	return nil
}

// Detect is used to detect wether the image contains explicit content
func (v Vision) Detect(fileName string, hash string) (visionTypes VisonTypes, err error) {
	ctx := context.Background()

	// Check the Redis database to reduce API calls
	state := v.Redis.Get(ctx, hash).Val()
	if state != "" {
		if state != ProperContent.String() {
			return ParseVsionTypes(state), ErrVisionFailed
		}

		return ProperContent, nil
	}

	file, err := os.Open(fileName)
	if err != nil {
		return UnknwonContent, ErrInternalServerError
	}
	defer file.Close()

	image, err := cloudVision.NewImageFromReader(file)
	if err != nil {
		return UnknwonContent, ErrInternalServerError
	}

	props, err := v.Client.DetectSafeSearch(context.Background(), image, nil)
	if err != nil {
		return UnknwonContent, ErrInternalServerError
	}

	adult := props.Adult.Enum().String()
	violence := props.Violence.Enum().String()
	spoof := props.Spoof.Enum().String()
	medical := props.Medical.Enum().String()
	racy := props.Racy.Enum().String()

	if adult == VeryLikely.String() || adult == Likely.String() || adult == Possible.String() {
		return AdultContent, ErrVisionFailed
	}
	if violence == VeryLikely.String() {
		return ViolenceContent, ErrVisionFailed
	}

	if spoof == VeryLikely.String() {
		return SpoofContent, ErrVisionFailed
	}

	if medical == VeryLikely.String() {
		return MedicalContent, ErrVisionFailed
	}

	if racy == VeryLikely.String() {
		return RacyContent, ErrVisionFailed
	}

	v.Redis.Set(ctx, hash, ProperContent.String(), 0)
	return ProperContent, nil
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

// String convert the given vision type enum to a string
func (v VisonTypes) String() string {
	return string(v)
}

// ParseVsionTypes parse the string format of the enum to the relevant enum
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

// Sevearity is an Enum to represent the vision sevearity
type Sevearity string

// String convert the VisonSevearity enum to its string representative
func (vs Sevearity) String() string {
	return string(vs)
}

const (
	// VeryLikely is to represent that the image has a high risk in the given category
	VeryLikely Sevearity = "VERY_LIKELY"
	// Likely is to represent that the image has a medium risk in the given category
	Likely Sevearity = "LIKELY"
	// Possible is to represent that the image has a low risk in the given category
	Possible Sevearity = "POSSIBLE"
)

var (
	// ErrVisionFailed is to represent that the vision API failed
	ErrVisionFailed = fmt.Errorf("cloud vision failed")
	// ErrInternalServerError is to represent that the server failed
	ErrInternalServerError = fmt.Errorf("something went wrong")
)
