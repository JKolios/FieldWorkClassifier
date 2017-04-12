# Scaling Proposal

### Overview

The underlying architecture of the FieldWorkClassifier service was
designed with easy horizontal scaling as a primary goal. Both the
overall microservices-based architecture and the choice of Elasticsearch
as its datastore and analytics engine contribute to that goal.

On the most basic level, the system can scale up horizontally by
adding additional instances of the Indexer and QueryRunner services
and distributing all incoming requests between them by using a Load
Balancer (for example, an Amazon ELB). Since the responsibilities of
handling incoming data and handling of queries has been split into
different services, it's possible to scale the one independent from
the other according to usage patterns.

Elasticsearch is an industry-proven and mature search and analytics
engine, with a proven track record for handling very large scale datasets
and generating complex analytics while staying highly performant. One
notable use of Elasticsearch on very large datasets has been at
[CERN](https://www.elastic.co/blog/grid-monitoring-at-cern-with-elastic).


## Scaling for a higher Data Ingestion rate

FieldWorkClassifier is designed so that the bulk of the processing needed
to classify an incoming payload is done at indexing-time (i.e. when
Elasticsearch first stores an incoming payload). This is accomplished by
the use of predefined [percolation queries](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-percolate-query.html).
This was a conscious design choice, as it allows us to have a faster and more responsive query
system by moving the "heavy lifting" to the indexing side. This means that
the Indexer is the most resource-intensive of the two services.

The intended scaling strategy for higher Data Ingestion rate is:

1. Using the [Elasticsearch Bulk Indexing](https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-bulk.html) to
API. This will increase indexing speed while introducing a slight latency
between the time a payload arrives and the time the data is available to
be queried. Real-time indexing does not seem mission critical for this
service.

2. Horizontally scaling up both the Indexer services and the Elasticsearch
cluster.

3. Introducing a Task Queue (for example, RabbitMQ), to store incoming
data before it is passed on to the Indexers. This would spread the processing
load evenly over time, while introducing some further latency.


## Scaling for a higher Query Volume

The intended scaling strategy for a higher Query Volume is:

1. Horizontally scaling up both the QueryRunner services and the Elasticsearch
cluster. Scaling up the Elasticsearch cluster should provide a significant
performance increase, as Elasticsearch uses
[Sharding](https://www.elastic.co/guide/en/elasticsearch/reference/current/_basic_concepts.html#getting-started-shards-and-replicas)
to parallelize query execution.

3. Introducing caching for Query Results. This could be done using a
 Redis cluster, for example. Since having up-to-the-second results is
 not a hard requirement, results can be cached to reduce the load on
 the QueryRunner instances.