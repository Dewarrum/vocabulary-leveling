FROM docker.elastic.co/elasticsearch/elasticsearch:8.14.1
RUN elasticsearch-plugin install analysis-nori