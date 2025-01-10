package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bluicezhen/emqx-s3-backup/lib/emqx"
	s3 "github.com/bluicezhen/emqx-s3-backup/lib/s3upload"
)

func main() {
	tempFile := getEmqxData()
	defer tempFile.Close()

	upload2s3(tempFile)
}

/**
 * @description: Upload a data backup file to S3
 * @param {*os.File} tempFile
 */
func upload2s3(tempFile *os.File) {
	logger := log.New(os.Stdout, "emqx-s3-backup: ", log.LstdFlags|log.Lshortfile)

	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		panic("S3_BUCKET is not set")
	}

	region := os.Getenv("S3_REGION")
	if region == "" {
		panic("S3_REGION is not set")
	}

	path := os.Getenv("S3_PATH")

	err := s3.UploadFile(tempFile, bucket, filepath.Base(tempFile.Name()), region, path)
	if err != nil {
		logger.Fatalf("Failed to upload data to S3: %v", err)
	}

	logger.Printf("Data uploaded to S3: %s", tempFile.Name())
}

/**
 * @description: Get a data backup file from EMQX
 * @return {*os.File} tempFile
 */
func getEmqxData() *os.File {
	logger := log.New(os.Stdout, "emqx-s3-backup: ", log.LstdFlags|log.Lshortfile)

	filename, err := emqx.NewEMQX().DataExport()
	if err != nil {
		logger.Fatalf("Failed to export data: %v", err)
	}
	logger.Printf("Data exported to file: %s", filename)

	tempFile, err := emqx.NewEMQX().DownloadData(filename.Filename, filename.Node)
	if err != nil {
		logger.Fatalf("Failed to download data: %v", err)
	}
	logger.Printf("Data downloaded to file: %s", tempFile.Name())

	return tempFile
}
