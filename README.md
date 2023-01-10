# aws-vmstate

## Build docker image
Pre-requisite: Need to login to any image registry and replace registry in the command below

```
docker build -t  quay.io/talat_shaheen0/aws-vmstate:latest .
```

## Run docker image
Pre-requisite: Need to export the AWS params

```
docker run -e AWS_DEFAULT_REGION=us-east-1 -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID -it quay.io/talat_shaheen0/aws-vmstate:latest

docker run -e AWS_DEFAULT_REGION=us-east-1 -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID -it quay.io/talat_shaheen0/aws-vmstate:latest
```
