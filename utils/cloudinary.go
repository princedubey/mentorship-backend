package utils

import (
	"context"
	"mime/multipart"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

// InitCloudinary initializes the Cloudinary client
func InitCloudinary() error {
	var err error
	cld, err = cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return err
	}
	return nil
}

// UploadImage uploads an image to Cloudinary and returns the URL
func UploadImage(file *multipart.FileHeader, folder string) (string, error) {
	ctx := context.Background()

	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Upload to Cloudinary
	uploadResult, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

// DeleteImage deletes an image from Cloudinary
func DeleteImage(publicID string) error {
	ctx := context.Background()
	
	// Extract the public ID from the URL
	if strings.Contains(publicID, "/") {
		publicID = strings.Split(publicID, "/")[1]
	}

	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

// DeleteImagesFromPost deletes all images associated with a post
func DeleteImagesFromPost(mediaURLs []string) error {
	var err error
	for _, url := range mediaURLs {
		if url != "" {
			err = DeleteImage(url)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
