package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type CloudinaryClient struct {
	cld *cloudinary.Cloudinary
}

func Init() *CloudinaryClient {
	cld, err := cloudinary.New()
	if err != nil {
		return nil
	}

	return &CloudinaryClient{cld: cld}
}

func (c *CloudinaryClient) Upload(
	ctx context.Context,
	file multipart.File,
	listingId int,
) (*uploader.UploadResult, error) {
	return c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:     uuid.NewString(),
		Folder:       fmt.Sprintf("live-listings-server/listings/%d", listingId),
		ResourceType: "image",
	})
}
