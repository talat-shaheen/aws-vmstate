FROM registry.access.redhat.com/ubi8/go-toolset:1.18.4-8.1669838000

# Labels
LABEL name="aws-vmstate" \
    maintainer="xyzcompany.com" \
    vendor="xyzcompany" \
    version="1.0.0" \
    release="1" \
    summary="This service enables state management of AWS cloud vms." \
    description="This service enables state management AWS cloud vms."

# copy code to the build path
COPY go.* ./
COPY aws-vmstate.go .

RUN go mod download

RUN go build -o aws-vmstate

CMD ["bash","-c","./aws-vmstate "]
