#!/bin/bash
echo "
	8888888      888888                 8   8  8                   
	8  8  8 eeee 8    8 e   e eeee eeee 8   8  8 eeeee eeeee  eeee 
	8e 8  8 8  8 8eeee8 8   8 8    8    8e  8  8 8   8 8   8  8    
	88 8  8 8e   88     8eee8 8eee 8eee 88  8  8 8eee8 8eee8e 8eee 
	88 8  8 88   88     88  8 88   88   88  8  8 88  8 88   8 88   
	88 8  8 88e8 88     88  8 88ee 88ee 88ee8ee8 88  8 88   8 88ee 


Ensure you already have the Google Cloud cli installed and setup https://cloud.google.com/sdk/docs/install
Ensure you have terser cli installed and setup https://github.com/terser/terser

This is a automation building script to deploy the banking PWA from: https://github.com/mcphee11/mcphee11-tui

enter ctrl+c to exit.

"
read -p "Continue? (Y/N): " confirm && [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]] || exit 1

#clean public folder and recreate
if [ -d public ]
then
echo "removing public folder"
rm -r public
mkdir public
mkdir public/AppImages
else
echo "making public folder"
mkdir public
fi

# copy app images into dir
cp -r AppImages public/

#minify & uglify root files
terser script.js -c -m  -o public/script.min.js
terser service-worker.js -c -m  -o public/service-worker.min.js
terser genesys.js -c -m  -o public/genesys.min.js
# move other files
cp -r svgs public/svgs
cp index.css public/index.css
cp manifest.json public/manifest.json
#update script.min.js
sed -i "/service-worker.js/s//service-worker.min.js/" public/script.min.js

# update index.html to use script.min.js & logos
sed '/src=\"script.js\"/s//src=\"script.min.js\"/' index.html > public/index.html

# update home.html to use script.min.js & logos
sed '/src=\"script.js\"/s//src=\"script.min.js\"/' home.html > public/home.html
sed -i '/src=\"genesys.js\"/s//src=\"genesys.min.js\"/' public/home.html

echo "finished building to public folder.
Do you want to upload to your Google Cloud bucket $bucketName..."
read -p "Continue? (Y/N): " confirm && [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]] || exit 1

#remove existing files from GCP
gcloud storage rm --recursive gs://$buckName/$shortName/

#copy folder to GCP
gcloud storage cp --recursive public/. gs://$buckName/$shortName/

#public access
gcloud storage objects update --recursive gs://$buckName/$shortName/ --add-acl-grant=entity=AllUsers,role=READER

echo "Completed build"