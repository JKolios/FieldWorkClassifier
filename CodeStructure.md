Code Structure
==============

FieldWorkClassifier contains the following 5 Go packages:

1. **Common** Contains common utility packages which are reused throughout the system.

    Package **config** contains a JSON config file structure and parsing code.

    Package **esclient** contains a portable method for connecting to ES clusters.

    Package **geojson** contains types for parsing GeoJSON-formatted files.

    Package **utils** contains a portable method for logging fatal errors.


2. **Indexer**:  Contains the application code for the Indexer microservice. This is further split into the **api** package, which implements the Indexer's HTTP and websocket API, and the **es** package, which contains all the code that deals with elasticsearch index setup, classification and the indexing of incoming documents. The Indexer service also receives and processes the GeoJSON data that describes the locations of the fields used when classifying data.

3. **QueryRunner**:  Contains the application code for the QueryRunner microservice. This is further split into the **api** package, which implements the QueryRunner's HTTP API, and the **es** package, which handles the creation of elasticsearch queries and the parsing of their results.

4. **PolygonImporter**:  A utility application that parses GeoJSON files containing field locations and imports them into elasticsearch via the Indexer.

5. **RandomDataGenerator**: A utility application which generates random "device" data payloads. Used for testing the system.


Both microservices have been intentionally divided into an API and an Elasticsearch (query) layer. This decouples their functionality and allows the query layer to interact with other potential data sources. Potential sources relevant to this project could be a Task Queue or a cache implementation (in the case of the QueryRunner).