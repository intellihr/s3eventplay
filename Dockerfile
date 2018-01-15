# build stage
FROM intellihr/s3eventplay:onbuild AS build
ADD . ./
RUN make build-alpine

# final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=build /go/src/github.com/intellihr/s3eventplay/s3eventplay .
ENTRYPOINT ["./s3eventplay"]
CMD ["help"]
