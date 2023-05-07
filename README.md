## Rust Showcase QR Backend

This repository contains the backend implementation for the Rust Showcase QR application. The backend exposes RESTful endpoints to perform CRUD operations on `opdata` resource. The following endpoints are available:

- `GET` http://localhost:8080/opdata
- `GET` http://localhost:8080/opdata/{id}
- `POST` http://localhost:8080/opdata
- `PUT` http://localhost:8080/opdata/{id}
- `DELETE` http://localhost:8080/opdata/{id}

These endpoints can be used to fetch, create, update or delete `opdata` records in the system. 

Please note that this implementation assumes a running instance of MongoDB on the system, and expects the necessary configuration parameters to be available as environment variables or in a `.env` file.
