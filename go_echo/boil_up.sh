#!/usr/bin/env bash

cd $(dirname $0)

sqlboiler psql --output models --pkgname models --wipe