# Go Ecommerce API

A complete ecommerce API solution built with Go, featuring Stripe payment integration, JWT-based authentication, and comprehensive admin functionality for managing products and collections.
This repo is not for commercial use and is a submission for roadmap.sh's Go Ecommerce API project (https://roadmap.sh/projects/ecommerce-api)[https://roadmap.sh/projects/ecommerce-api].

## Features

- **JWT Authentication**: Secure admin routes and cart sessions
- **Product Management**: CRUD operations for products
- **Collection Management**: Group products into collections
- **Cart Functionality**: Add, update, remove items with session persistence
- **Stripe Integration**: Process payments with Stripe checkout sessions
- **Admin Portal**: Protected routes for store management

## Prerequisites

- Go 1.16+
- PostgreSQL
- [Air](https://github.com/cosmtrek/air) for hot reloading (optional)
- Stripe account for payment processing

## Environment Variables

Create a `.env` file with the following variables:

```
DB_URL=postgresql://username:password@localhost:5432/dbname
JWT_SECRET=your_jwt_secret_key
STRIPE_SECRET_KEY=your_stripe_secret_key
COOKIE_NAME=your_cookie_name
```

## Getting Started

### Installation

1. Clone the repository
2. Install dependencies:

```bash
go mod download
```

3. Install Air for hot reloading (optional):

```bash
make install-air
```

### Database Setup

1. Create the database schema:

```bash
psql -d your_database_name -f schema.sql
```

2. Seed the database with initial data:

```bash
make run-db
```

### Running the Application

Development mode with hot reloading:

```bash
make dev
```

Build and run the application:

```bash
make build
make run
```

Clean build artifacts:

```bash
make clean
```

## API Routes

### Public Routes

- `GET /api/` - Health check
- `GET /api/products/` - List all products
- `GET /api/products/{id}` - Get product details
- `GET /api/collections/` - List all collections
- `GET /api/collections/{id}` - Get collection details
- `POST /api/new-cart` - Create a new cart session with JWT

### Cart Routes (JWT Protected)

- `GET /api/cart/` - View cart contents
- `POST /api/cart/add` - Add item to cart
- `PUT /api/cart/{productID}/` - Update item quantity
- `DELETE /api/cart/{productID}/` - Remove item from cart
- `DELETE /api/cart/` - Clear cart
- `POST /api/cart/checkout` - Create Stripe checkout session

### Authentication Routes

- `POST /api/auth/login` - Admin login
- `POST /api/auth/logout` - Admin logout

### Admin Routes (JWT Protected)

- `GET /api/admin/` - Admin portal access
- `POST /api/admin/users/register` - Register admin user
- `GET /api/admin/users/{id}/` - Get user details
- `DELETE /api/admin/users/{id}/` - Delete user

#### Product Management (Admin)

- `GET /api/admin/products/` - List all products
- `POST /api/admin/products/` - Create new product
- `GET /api/admin/products/{id}/` - Get product details
- `PUT /api/admin/products/{id}/` - Update product
- `DELETE /api/admin/products/{id}/` - Delete product

#### Collection Management (Admin)

- `GET /api/admin/collections/` - List all collections
- `POST /api/admin/collections/` - Create new collection
- `GET /api/admin/collections/{id}/` - Get collection details
- `PUT /api/admin/collections/{id}/` - Update collection
- `DELETE /api/admin/collections/{id}/` - Delete collection
- `POST /api/admin/collections/{id}/product/{id}/` - Add product to collection
- `DELETE /api/admin/collections/{id}/product/{id}/` - Remove product from collection

## Project Structure

```
go-ecommerce-api/
├── cmd/
│   ├── api/        # API entry point
│   └── db/         # Database seeding
├── internal/
│   ├── auth/       # Authentication middleware
│   ├── db/         # Database models and queries
│   ├── handlers/   # HTTP handlers
│   └── methods/    # Business logic
├── schema.sql      # Database schema
├── query.sql       # SQLC queries
└── sqlc.yaml       # SQLC config
```

## License

[MIT](LICENSE)

