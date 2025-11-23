# Simple Inventory API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Most endpoints require authentication via Bearer token in the Authorization header:
```
Authorization: Bearer <token>
```

## API Endpoints

### Authentication

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure123"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "secure123"
}
```

**Response:**
```json
{
  "token": "abc123...",
  "expires_at": "2024-01-01T00:00:00Z",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "is_active": true
  }
}
```

#### Logout
```http
POST /api/v1/auth/logout
Authorization: Bearer <token>
```

#### Get Current User
```http
GET /api/v1/auth/me
Authorization: Bearer <token>
```

---

### Products

#### Create Product
```http
POST /api/v1/products
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Laptop",
  "sku": "LAP-001",
  "barcode": "1234567890",
  "description": "Dell Laptop",
  "quantity": 10,
  "min_quantity": 5,
  "price": 999.99,
  "location_id": 1
}
```

#### Get Product
```http
GET /api/v1/products/:id
Authorization: Bearer <token>
```

#### List Products
```http
GET /api/v1/products?limit=20&offset=0
Authorization: Bearer <token>
```

#### Search Products
```http
GET /api/v1/products/search?q=laptop&limit=20&offset=0
Authorization: Bearer <token>
```

#### Get Low Stock Products
```http
GET /api/v1/products/low-stock?limit=20&offset=0
Authorization: Bearer <token>
```

#### Scan Barcode
```http
GET /api/v1/products/scan?barcode=1234567890
Authorization: Bearer <token>
```

#### Update Product
```http
PUT /api/v1/products/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Laptop",
  "price": 899.99
}
```

#### Delete Product
```http
DELETE /api/v1/products/:id
Authorization: Bearer <token>
```

---

### Locations

#### Create Location
```http
POST /api/v1/locations
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Main Warehouse",
  "code": "WH-001",
  "description": "Main storage facility",
  "building": "Building A",
  "floor": "Floor 1",
  "aisle": "A1",
  "shelf": "S1"
}
```

#### Get Location
```http
GET /api/v1/locations/:id
Authorization: Bearer <token>
```

#### List Locations
```http
GET /api/v1/locations?limit=20&offset=0
Authorization: Bearer <token>
```

#### Search Locations
```http
GET /api/v1/locations/search?q=warehouse&limit=20&offset=0
Authorization: Bearer <token>
```

#### Update Location
```http
PUT /api/v1/locations/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Warehouse",
  "is_active": true
}
```

#### Delete Location
```http
DELETE /api/v1/locations/:id
Authorization: Bearer <token>
```

---

### Inventory Management

#### Receive Inventory
```http
POST /api/v1/inventory/receive
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": 1,
  "quantity": 50,
  "location_id": 1,
  "notes": "New shipment received"
}
```

#### Issue Inventory
```http
POST /api/v1/inventory/issue
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": 1,
  "quantity": 10,
  "notes": "Issued to production"
}
```

#### Adjust Inventory
```http
POST /api/v1/inventory/adjust
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": 1,
  "new_quantity": 45,
  "notes": "Stock count adjustment"
}
```

#### Transfer Inventory
```http
POST /api/v1/inventory/transfer
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": 1,
  "quantity": 20,
  "from_location_id": 1,
  "to_location_id": 2,
  "notes": "Transfer to secondary warehouse"
}
```

#### Get Product Transactions
```http
GET /api/v1/inventory/transactions/product/:product_id?limit=20&offset=0
Authorization: Bearer <token>
```

---

## Error Responses

All errors follow this format:
```json
{
  "error": "error_code",
  "message": "Detailed error message"
}
```

Common HTTP Status Codes:
- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Authentication required
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Barcode Printing

To implement barcode printing, you can use libraries like:
- Go: `github.com/boombuler/barcode`
- Frontend: `JsBarcode` or `react-barcode`

Example barcode generation can be added as a separate endpoint or handled on the frontend.
