# Inventory Manager (backend)
This repo contains the backend for the inventory manager. It is a gofiber app that uses a postgres database to store information about the inventories.

## Getting Started

### Prerequisites

- Go
- Postgres

### Installing

1. Clone the repo
2. Run `go build` in the root directory
3. Run the binary

## Usage

### Creating an inventory

To create an inventory, you can use the following command:

```
curl -X POST -H "Content-Type: application/json" -d '{"name": "My Inventory"}' http://localhost:3000/api/inventories
```

This will create a new inventory with the name "My Inventory".

### Getting all inventories

To get all inventories, you can use the following command:

```
curl -X GET http://localhost:3000/api/inventories
```

This will return a JSON array containing all the inventories.

### Getting a specific inventory

To get a specific inventory, you can use the following command:

```
curl -X GET http://localhost:3000/api/inventories/1
```

This will return a JSON object containing the details of the inventory with the ID 1.

### Updating an inventory

To update an inventory, you can use the following command:

```
curl -X PUT -H "Content-Type: application/json" -d '{"name": "New Name"}' http://localhost:3000/api/inventories/1
```

This will update the name of the inventory with the ID 1 to "New Name".

### Deleting an inventory

To delete an inventory, you can use the following command:

```
curl -X DELETE http://localhost:3000/api/inventories/1
```

This will delete the inventory with the ID 1.
