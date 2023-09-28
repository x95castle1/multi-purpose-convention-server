#!/bin/bash

GLOBIGNORE=metadata.yml
# Get the list of files sorted by version number
files=$(ls -v ./carvel/packagerepository/packages/multi-purpose-convention-server.conventions.tanzu.vmware.com/*.yml)
# Count the number of files
count=$(echo "$files" | wc -l)
# Calculate how many files to delete
del=$(expr $count - 5)
# If there are more than 5 files, delete the oldest ones
if [ $del -gt 0 ]; then
  echo "Deleting $del files"
  # Loop through the first $del files and delete them
  i=1
  while [ $i -le $del ]; do
    file=$(echo "$files" | sed -n ${i}p)
    echo "Deleting $file"
    rm $file
    i=$(expr $i + 1 )
  done
else
  echo "Nothing to delete"
fi
# Unset the GLOBIGNORE variable
unset GLOBIGNORE