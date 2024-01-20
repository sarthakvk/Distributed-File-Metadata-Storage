# Distributed File Metadata Storage

**Project Overview:**

This project focuses on building a Go application adhering to the principles of Hexagonal Architecture. The application is designed to accomplish the following key tasks:

1. **File Processing:** Calculate file sizes accurately.
2. **Database Update:** Dynamically update a database with fileName-fileSize mappings using a decentralized consensus mechanism.
3. **Logging:** Implement a comprehensive Logger to manage Log, Trace, Debug, and Metrics for the API flow.


**Application Documentation:**

The source code is meticulously self-documented with insightful comments, ensuring not only clarity but also facilitating seamless maintenance.

**Application Usage:**

* **Application Prerequisites:**
  - Begin by cloning the repository and navigating to the project's root directory. All subsequent commands assume execution from this root directory.
  - Launch Redis containers as specified in the `docker-compose.yml` file. Each Redis container will be exclusively utilized by a designated node.
    ```bash
    docker compose up
    ```

* **Application Startup:**
  - Execute the following commands to initiate the nodes:
    ```bash
    go run cmd/httpd/httpd.go -node-id A -address localhost:8000 -http-port 9000 -redis-addr localhost:6379 -bootstrap
    ```
    ```bash
    go run cmd/httpd/httpd.go -node-id B -address localhost:8001 -http-port 9001 -redis-addr localhost:7379
    ```
    ```bash
    go run cmd/httpd/httpd.go -node-id C -address localhost:8002 -http-port 9002 -redis-addr localhost:8379
    ```

  After successfully launching the nodes, ensure that the leader node (i.e., `node A`, as it was designated with the bootstrap option) adds the other two nodes to the cluster.

* **Adding Nodes API:**
  - Add node B to the cluster:
    ```bash
    curl --location 'localhost:9000/add-replica' \
    --header 'Content-Type: application/json' \
    --data '{
        "node_id": "B",
        "address": "localhost:8001"
    }'
    ```
  - Add node C to the cluster:
    ```bash
    curl --location 'localhost:9000/add-replica' \
    --header 'Content-Type: application/json' \
    --data '{
        "node_id": "C",
        "address": "localhost:8002"
    }'
    ```

  All nodes are now seamlessly integrated into the cluster, each equipped with its dedicated Redis instance for storage.

* **File Uploading API:**
  - To upload a file, utilize the following command. It is imperative to note that file uploads are exclusively processed by the leader node, aligning with the fundamental RAFT principle where all operations transpire on the leader.
    ```bash
    curl --location 'localhost:9000/upload-file' \
    --form 'file=@"<FULL PATH OF THE FILE>"'
    ```

* **Verifying File Upload:**
  - To authenticate the file's upload and retrieve its size, leverage the following API. This validation process can also be conducted on replica nodes to ensure unwavering data consistency.
    ```bash
    curl --location 'localhost:9000/key-store' \
    --header 'Content-Type: application/json' \
    --data '{
        "command": "GET",
        "key": "<FileName>"
    }'
    ```

- Additionally, the API can serve as a distributed key-value store. It is crucial to execute all commands on the leader node. While running `GET` on non-leader nodes provides the value stored in that specific node, operations like `SET`, `DELETE`, and `GET_OR_CREATE` are exclusively processed on the leader node.

- **Examples:**
  1. **SET:**
     - **Request:**
          ```bash
          curl --location 'localhost:9000/key-store' \
          --header 'Content-Type: application/json' \
          --data '{
              "command": "SET",
              "key": "FOO",
              "value": "BAR"
          }'
          ```
      - **Response:**
          ```json
          {"created":true}
          ```
  2. **DELETE:**
      ```bash
      curl --location 'localhost:9000/key-store' \
      --header 'Content-Type: application/json' \
      --data '{
          "command": "DELETE",
          "key": "FOO"
      }'
      ```
 1. **GET_OR_CREATE:**
    - **Request:**
      ```bash
      curl --location 'localhost:9000/key-store' \
      --header 'Content-Type: application/json' \
      --data '{
          "command": "GET_OR_CREATE",
          "key": "FOO",
          "value": "FOO"
      }'
      ```
     - **Response:**
      ```json
      {"value":"BAR"}
      ```
