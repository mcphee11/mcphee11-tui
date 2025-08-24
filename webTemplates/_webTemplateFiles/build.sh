#!/bin/bash

########################
# Created as an example by https://github.com/mcphee11 Version 1
# used to build single html file
# requires: terser
# requires: uglifycss
# requires: html-minifier
########################

########## Minify the JavaScript file ##########
terser index.js -c -m -o index.min.js

# Read the minified JavaScript file content into a variable
JS_CONTENT=$(<index.min.js)

# Create a temporary file for our new content
TMP_FILE_JS=$(mktemp)

# Create a copy of index.html with the new JavaScript content
while IFS= read -r linejs; do
  if [[ "$linejs" =~ \/\/MINIFIED_JS ]]; then
    printf '%s' "$JS_CONTENT" >> "$TMP_FILE_JS"
  else
    printf '%s\n' "$linejs" >> "$TMP_FILE_JS"
  fi
done < index.html

# Now, use grep to create the index.min.html file.
grep -v '<script src="/index.js"></script>' "$TMP_FILE_JS" > index.min.html

# Remove the temporary file
rm "$TMP_FILE_JS"
rm index.min.js

echo "Minified & embedded index.js..."

########## Minify the CSS file ##########
uglifycss index.css > index.min.css

# Read the minified CSS file content into a variable
CSS_CONTENT=$(<index.min.css)

# Create a temporary file for our new content
TMP_FILE_CSS=$(mktemp)

# Create a copy of index.min.html with the new CSS content
while IFS= read -r linecss; do
  if [[ "$linecss" =~ \#MINIFIED_CSS ]]; then
    printf '%s' "$CSS_CONTENT" >> "$TMP_FILE_CSS"
  else
    printf '%s\n' "$linecss" >> "$TMP_FILE_CSS"
  fi
done < index.min.html

# Now, use grep to create the index.min.html file.
grep -v '<link href="/index.css" rel="stylesheet" />' "$TMP_FILE_CSS" > index.min.html

# Remove the temporary file
rm "$TMP_FILE_CSS"
rm index.min.css

echo "Minified & embedded index.css..."

########## Minify entire HTML file now ##########
html-minifier --collapse-whitespace --remove-comments index.min.html

echo "\nFinished"
