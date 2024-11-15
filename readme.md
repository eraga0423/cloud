# Triple-S (Simple Storage Service)

**Triple-S** is a simplified cloud storage service, modeled after Amazon S3. It provides RESTful API endpoints for managing storage buckets and objects (files) through HTTP. The project allows users to create, list, and delete buckets, as well as upload, retrieve, and delete objects within those buckets. All operations are stored in local directories and metadata is saved in CSV files.

## Features

- **Bucket Management**: 
  - Create buckets with unique, valid names.
  - List all existing buckets.
  - Delete empty buckets.

- **Object Operations**:
  - Upload objects to a bucket.
  - Retrieve object content.
  - Delete objects from a bucket.

- **RESTful API**:
  - Endpoints for bucket and object management.
  - HTTP status codes for error handling and success messages.
  
- **CSV Metadata Storage**:
  - Bucket metadata stored in `buckets.csv`.
  - Object metadata stored in `objects.csv` within each bucket's directory.

## Installation

1. Clone the repository.
2. Build the project:
   ```bash
   $ go build -o triple-s .



