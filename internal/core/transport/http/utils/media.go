package core_http_utils

import (
	"fmt"
	"strings"
)

const (
	TagImages  = "images"
	TagAvatars = "avatars"
	TagDocs    = "docs"
)

func ValidateTagParam(tag string) error {
	allowedTags := map[string]struct{}{}
	tags := fmt.Sprintf("%s,%s,%s", TagAvatars, TagDocs, TagImages)

	for allowedTag := range strings.SplitSeq(tags, ",") {
		allowedTags[strings.TrimSpace(allowedTag)] = struct{}{}
	}

	if _, ok := allowedTags[tag]; !ok {
		return fmt.Errorf("invalid media s3 tag")
	}

	return nil
}

func GetBucketFolder(tag string) string {
	var folder string

	switch tag {
	case TagImages:
		folder = "images/"
	case TagAvatars:
		folder = "avatars/"
	case TagDocs:
		folder = "documents/"
	default:
		folder = "misc/"
	}

	return folder
}
