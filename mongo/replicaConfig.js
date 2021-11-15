db = new Mongo("localhost:27017").getDB("test");

config = {
  _id: "myReplica",
  members: [{ _id: 0, host: "localhost:27017" }],
};
rs.initiate(config);

uri =
  "mongodb://localhost:27017,localhost:27018,localhost:27019/test?replicaSet=my-mongo-set";
