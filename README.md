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

### Run With Binary

```
./s3eventplay --help
```

### Use Inside Docker image

```
# Include "s3eventplay"
RUN curl -sL -o /usr/local/bin/s3eventplay \
    https://github.com/intellihr/s3eventplay/releases/download/v0.0.2/s3eventplay_0.0.2_linux_amd64 \
    && chmod +x /usr/local/bin/s3eventplay
```


## Usage

```
USAGE:
   s3eventplay [global options] command [command options] [arguments...]

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --bucket value, -b value  target AWS S3 bucket, required [$S3_BUCKET]
   --stream value, -s value  target kinesis stream name or folder name that stores the event files, required [$EVENT_STREAM]
   --batch value, -n value   target batch size (i.e. number of concurrent file downloads). (default: 5) [$BATCH_SIZE]
   --dates value, -d value   target dates (date range) to play the events (e.g. YYYY-MM-DD or YYYY-MM-DD~YYYY-MM-DD) [$EVENT_DATES]
   --help, -h                show help
   --version, -v             print the version
```
