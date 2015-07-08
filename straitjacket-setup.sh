#!/bin/bash
set -e

cp -R files/etc/apparmor.d/* /etc/apparmor.d/
service apparmor reload

for language in languages/*; do
  image_name=straitjacket-$(basename $language)
  echo "Building image $image_name"
  docker build -t $image_name $language
done
