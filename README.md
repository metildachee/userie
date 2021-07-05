# Architecture
![image](https://user-images.githubusercontent.com/65015150/124505374-eeed7000-ddfb-11eb-81ab-83f42609f66b.png)

# Set up env
## Install and run Elastic search 
1. Install and start [Docker Desktop](https://www.docker.com/products/docker-desktop)
2. Run 
    ```
    docker network create elastic
    docker pull docker.elastic.co/elasticsearch/elasticsearch:7.13.2
    docker run --name es01-test --net elastic -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.13.2
    ```
3. To verify it is running successfully
    ```
    curl -X GET http://localhost:9200
    ```
## Install and run Kibana
1. In a new terminal session, run
    ```
    docker pull docker.elastic.co/kibana/kibana:7.13.2
    docker run --name kib01-test --net elastic -p 5601:5601 -e "ELASTICSEARCH_HOSTS=http://es01-test:9200" docker.elastic.co/kibana/kibana:7.13.2
    ```
2. To access Kibana, go to http://localhost:5601/app/dev_tools#/console

### Add mapping in Kibana
1. Create index
    ```
    PUT /usersg0
    ```
2. Add mapping
    ```
    PUT /usersg0/_mapping
    {
      "properties": {
        "name": {
          "type": "text"
        },
        "dob": {
          "type": "long"
        },
        "address": {
          "type": "text"
        },
        "description": {
          "type": "text"
        },
        "ctime": {
          "type": "long"
        }
      }
    }
    ```

### Install and start tracer (logging)
1. Run
    ```
    docker pull jaegertracing/all-in-one
    docker run -d -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest
    ```
2. Go to http://localhost:16686/ to see traces

# Start server
There are 2 options to start the server
1. Build and run
    ```
    go build main.go
    nohup ./UserSearch &
    ```
2. Run locally
    ```
   go run main.go
   ```
# Testing
1. Testing api
    ```
    cd api
    go test Test -v
    ```
2. Testing dao
    ```
    cd dao
    go test Test -v
    ```
