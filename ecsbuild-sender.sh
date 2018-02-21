#!/bin/bash

tag=$1
if [[ -z "${tag}" ]]; then
    echo "Provide tag as argument."
    exit 1
fi

docker build --rm -f dockerfile.sender -t sender:${tag} .
docker tag sender:${tag} 963826138034.dkr.ecr.ap-northeast-1.amazonaws.com/sender:${tag}
`aws ecr get-login --no-include-email --region ap-northeast-1`
docker push 963826138034.dkr.ecr.ap-northeast-1.amazonaws.com/sender:${tag}
docker rmi $(docker images --filter "dangling=true" --no-trunc -q)