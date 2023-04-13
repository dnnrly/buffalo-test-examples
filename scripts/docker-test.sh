#!/usr/bin/env bash
set -e

sleep 10 # Wait for the DB to become available

buffalo db -e test migrate up
buffalo test