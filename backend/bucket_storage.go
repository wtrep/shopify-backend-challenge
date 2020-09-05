package backend

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"time"
)

func UploadToBucket(data io.Reader, bucket, object string) error {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, data); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}
