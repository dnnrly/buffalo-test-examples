#!/usr/bin/env bash
set -exo pipefail

sleep 10 # Wait for the DB to become available

buffalo db -e test migrate up
buffalo test