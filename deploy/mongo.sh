echo "Creating MongoDB network..."
docker network create mongoCluster

echo "Starting MongoDB container..."
docker run -d --rm \
  -p 27017:27017 \
  --name mongo1 \
  --network mongoCluster \
  mongo:5 \
  mongod --replSet myReplicaSet --bind_ip 0.0.0.0

echo "Waiting for MongoDB to start (20 seconds)..."
sleep 20

echo "Initializing replica set..."
docker exec -it mongo1 mongosh --eval "rs.initiate({
 _id: \"myReplicaSet\",
 members: [
   { _id: 0, host: \"mongo1:27017\" }
 ]
})"

echo "MongoDB setup complete!"