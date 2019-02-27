#!/usr/bin/env bash

# pull the official mongo docker container
docker pull mongo:4.0.5

# create network
docker network create mongo-network

# create mongos
docker run -d --net my-mongo-cluster -p 30000:27017 --name mongo1000 mongo mongod --replSet mongo-network --port 27017
docker run -d --net my-mongo-cluster -p 30001:27017 --name mongo2000 mongo mongod --replSet mongo-network --port 27017
docker run -d --net my-mongo-cluster -p 30002:27017 --name mongo3000 mongo mongod --replSet mongo-network --port 27017

# add hosts
# 127.0.0.1       mongo1 mongo2 mongo3

############## setup replica set ################
# docker exec -it mongo2000 mongo1000
# db = (new Mongo('localhost:27017')).getDB('test')
# config={"_id":"mongo-network","members":[{"_id":0,"host":"mongo1:27017"},{"_id":1,"host":"mongo2:27017"},{"_id":2,"host":"mongo3:27017"}]}
# rs.initiate(config)
#################################################

# connection URI
## mongodb://localhost:27017,localhost:27018,localhost:27019/{db}?replicaSet=mongo-network
