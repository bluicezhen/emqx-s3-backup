# emqx-s3-backup

EMQX-S3-Backup is a tool for seamlessly backing up EMQX data to AWS S3, ensuring secure and reliable data storage in the cloud.

## Usage

```bash
go run main.go
```

## Configuration

Environment variables:

- `EMQX_URL`: The URL of EMQX.
- `EMQX_API_NAME`: The API name of EMQX.
- `EMQX_API_PASS`: The API password of EMQX.
- `S3_BUCKET`: The bucket name of S3.
- `S3_REGION`: The region of S3.
- `S3_PATH`: The path of S3.
## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
