#!/bin/bash

# get all the app icons
echo "creating icons from $icon ..."
response=$(curl -o --request POST  -H "Content-Type: multipart/form-data" -F fileName=@$icon -F padding=0.3 -F color=transparent -F platform=windows11 -F platform=android -F platform=ios  https://appimagegenerator-prod-dev.azurewebsites.net/api/image)
echo "$response"
eval "$(jq -M -r '@sh "response_url=\(.Uri)"' <<< "$response")"
echo "Downloading from: https://appimagegenerator-prod-dev.azurewebsites.net$response_url ..."
wget https://appimagegenerator-prod-dev.azurewebsites.net$response_url
file=${response_url:5}

# unzip and app images into dir
unzip $file -d AppImages
rm $file