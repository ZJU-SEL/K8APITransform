#!/bin/bash
cp -r ../conf ./
cp ../ApiServer ./
docker build -t selapiserver .
docker save  selapiserver > SelApiServer.tar