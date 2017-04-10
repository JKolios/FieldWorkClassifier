FieldWorkClassifier
===================

A microservice-based system that generates farm work timetables based on
(emulated) JSON payloads from mobile devices. Uses ElasticSearch v5 for
 Indexing, Geo-queries and Aggregations. Written in Go.


Structure
---------

FieldWorkClassifier comprises two microservices: Indexer and QueryRunner.

Indexer's role is to receive JSON payloads from either HTTP or Websocket
API endpoints, preprocess and index them in a form allowing for efficient
searches.

QueryRunner runs Elasticsearch queries on the data that was stored by the
Indexer. The results of these queries are returned through its HTTP API.

The project also includes two helper utilities: PolygonImporter, which
imports GeoJSON-formatted field coordinates into the Indexer, and
RandomDataGenerator, which (as the name implies) generates random test
data.

Usage with Docker
-----------------

* Install Docker-compose, if needed.
* To start the service, run `docker-compose up --build`. It may take up
to 60 seconds for the ElasticSearch cluster to initialize. The
microservices will connect to the cluster as soon as it's available.
* Initialize the field geodata by running `go run PolygonImporter` under
`PolygonImporter`. You can use the provided polygons.json or substitute it
with any valid GeoJSON FeatureCollection file.
* Running `go run RandomDataGenerator` under `RandomDataGenerator` will
populate the cluster with randomized data.

Alternatively, Indexer and QueryRunner instances can also be created as
independent Docker containers.