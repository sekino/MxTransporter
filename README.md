[TODO]
insert mage

MxTransporter is a middleware that accurately carries change streams of MongoDB in real time. For infrastructure, you can easily use this middleware by creating a container image with Dockerfile on any platform and deploying it.

<br>

# Guide

---
## Build with samples
We have prepared a samples to build on AWS and GCP container orchestration services.
This can be easily constructed by setting environment variables as described and executing commands.

See ```docs/``` for details.

<br>

## Deploy to your container environment
With the Dockerfile provided, you can easily run MxTransporter by simply building the container image and deploying it to your favorite container environment.

### Requirement

- Build a Dockerfile, create an image, and create a container based on that image.

- Mount the persistent volume on the container to store the resume token. See the Change streams section of this README for more information.

- Allow access from the container to MongoDB

- Add permissions so that the container can access the export destination.

- Have the container read the required environment variables. All you need is a "Run locally" section in this README to add to your ```.env```.

<br>

## Run locally
### Requirement
- Set the following environment variables in ```.env```.

```
## mongodb+srv://<user name>:<passward>@xxx.yyy.mongodb.net/test?retryWrites=true&w=majority
MONGODB_HOST=
MONGODB_DATABASE=
MONGODB_COLLECTION=

## You have to specify this environment variable if you want to export BigQuery, Pub/Sub.
PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS=

## You have to specify this environment variable if you want to export BigQuery.
BIGQUERY_DATASET=
BIGQUERY_TABLE=

## You have to specify this environment variable if you want to export Kinesis Stream.
KINESIS_STREAM_NAME=
KINESIS_STREAM_REGION=

## One resume token is saved in the location specified here.
PERSISTENT_VOLUME_DIR=

## Specify the location you want to export. (ex. EXPORT_DESTINATION=bigquery )
EXPORT_DESTINATION=
```

- Allow access from the IP of the local machine on the mongoDB.


- Run ```go run ./cmd/main.go``` in the root directory.

<br>

# Architects

--- 

![image](https://user-images.githubusercontent.com/37132477/140257547-fd5417fe-abe3-4bdc-8aad-c08d96e19d0f.png)


<br>

# Specification

---

## MongoDB

### Connection to MongoDB
Allow the public IP of the MxTransporter container on the mongoDB side. This allows you to watch the changed streams that occur.

### Change streams
change streams output the change events that occurred in the database and are the same as the logs stored in oplog. And it has a unique token called resume token, which can be used to get events after a specific event.

In this system, resume token is saved in Persistent Volume associated with the container, and when a new container is started, the resume token is referenced and change streams acquisition starts from that point.

The resume token of the change streams just before the container stopped is stored in the persistent volume, so you can refer to it and get again the change streams that you missed while the container stopped and the new container started again.

The resume token is stored in the directory where the PVC is mounted.

```PERSISTENT_VOLUME_DIR``` is an environment variable given to the container.

```
{$PERSISTENT_VOLUME_DIR}/{year}/{month}/{day}
```

The resume token is saved in ```{year}-{month}-{day}.dat```.

```
$ pwd
{$PERSISTENT_VOLUME_DIR}/{year}/{month}/{day}

$ ls
{year}-{month}-{day}.dat

$ cat {year}-{month}-{day}.dat
T7466SLQD7J49BT7FQ4DYERM6BYGEMVD9ZFTGUFLTPFTVWS35FU4BHUUH57J3BR33UQSJJ8TMTK365V5JMG2WYXF93TYSA6BBW9ZERYX6HRHQWYS
```

When getting change-streams by referring to resumu token, it is designed to specify resume token in ```startAfrter``` of ```Collection.Watch()```.

<br>

## Export change streams
MxTransporter export change streams to the following description.

- Google Cloud BigQuery
- Google Cloud Pub/Sub
- Amazon Kinesis Data Streams

### BigQuery
Create a BigQuery Table with a schema like the one below.

Table schema
```
[
    {
      "mode": "NULLABLE",
      "name": "id",
      "type": "STRING"
    },
    {
      "mode": "NULLABLE",
      "name": "operationType",
      "type": "STRING"
    },
    {
      "mode": "NULLABLE",
      "name": "clusterTime",
      "type": "TIMESTAMP"
    },
    {
      "mode": "NULLABLE",
      "name": "fullDocument",
      "type": "STRING"
    },
    {
      "mode": "NULLABLE",
      "name": "ns",
      "type": "STRING"
    },
    {
      "mode": "NULLABLE",
      "name": "documentKey",
      "type": "STRING"
    },
    {
      "mode": "NULLABLE",
      "name": "updateDescription",
      "type": "STRING"
    }
]
```

### Pub/Sub
No special preparation is required. Create a Topic with the MongoDB Database name, and a Subscription with the MongoDB Collection name from which the change streams originated.

Change streams are sent to that subscription in a pipe (|) separated CSV.

### Kinesis Data Streams
No special preparation is required. If you want to separate the data warehouse table for each MongoDB collection for which you want to get change streams, use Kinesis Firehose and devise the output destination.

Change streams are sent to that in a pipe (|) separated CSV.

<br>

# Copyright

---

CAM, Inc. All rights reserved.