#!/usr/bin/env bash

while getopts t:v: flag
do
    case "${flag}" in
        t) tar=${OPTARG};;
        v) version=${OPTARG};;
    esac
done

fpm -s tar -t deb -C $tar --name node --version $version --iteration 1  --description "A verim node" .
