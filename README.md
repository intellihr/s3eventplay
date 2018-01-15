# S3 Event Play
Play Events From S3 Bucket **Parallelly** and **Sequentially** into `stdout`

## Running locally?

### Build
```
docker-compose build app
```

### Run With Docker

```
docker-compose build dev
docker-compose build app

# assume AWS credentials are stored in .env file
docker run --rm --env-file .env intellihr/s3eventplay:latest <options>
```
