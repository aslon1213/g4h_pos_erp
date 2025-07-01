echo "Creating MongoDB network..."
docker network create mongoCluster

echo "Starting MongoDB container..."
docker run -d  --rm \
  -p 27017:27017 \
  --name mongo1 \
  --network mongoCluster \
  -e MONGO_INITDB_ROOT_USERNAME=<username> \
  -e MONGO_INITDB_ROOT_PASSWORD=<password> \
  -v ./mongo-keyfile:/data/keyfile \
  mongo:5 \
  mongod --bind_ip_all --replSet myReplicaSet --auth --keyFile /data/keyfile

echo "Waiting for MongoDB to start (10 seconds)..."
sleep 10

echo "Initializing replica set..."
docker exec -it mongo1 mongosh -u <username> -p <password> --eval "rs.initiate({
 _id: \"myReplicaSet\",
 members: [
   { _id: 0, host: \"localhost:27017\" }
 ]
})"
# run this to change host resolving
docker exec -it mongo1 mongosh -u <username> -p <password> --eval "cfg = rs.conf() 
cfg.members[0].host = \"<host>:27017\"
rs.reconfig(cfg, { force: true })
"
echo "MongoDB setup complete!"


# run redis container as well
echo "Starting Redis container..."  
# docker run  -d --name redis --net mongoCluster -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
echo "Redis setup complete!"