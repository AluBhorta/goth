#!/bin/bash

openssl rand -base64 741 > mongodb.key
chmod 600 mongodb.key

mongod --keyFile mongodb.key --replSet myReplica
