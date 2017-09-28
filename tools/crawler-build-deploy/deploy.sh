#!/bin/bash
cd ../../
rsync -v -P -e ssh crawler "$1":spiderwoman/crawler