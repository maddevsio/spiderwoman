#!/bin/bash
cd ../../api
rsync -v -P -e ssh api "$1":spiderwoman/api/api
rsync -v -P -a -r -e ssh templates/ "$1":spiderwoman/api/templates
rsync -v -P -a -r -e ssh images/ "$1":spiderwoman/api/images
rsync -v -P -a -r -e ssh assets/ "$1":spiderwoman/api/assets