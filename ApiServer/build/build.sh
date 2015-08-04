#!/bin/bash
#cp -r ../conf ./
cp ../ApiServer ./
cp -r ../certs ./
docker build -t selapiserver .
docker save  selapiserver > SelApiServer.tar
