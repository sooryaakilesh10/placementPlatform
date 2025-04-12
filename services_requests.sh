#!/bin/bash

# test userd service test cases
cd services/userd && ./user_test.sh

# test questiond service endpoints
cd ../questiond && hurl --test question.hurl

# test testd service endpoints
cd ../testd && hurl --test test.hurl
