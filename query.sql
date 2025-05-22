-- Users
-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users 
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (
  email, password_hash, is_admin
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1 RETURNING *;


-- Products
-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY name;

-- name: CreateProduct :one
INSERT INTO products (
  name, description, price
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateProduct :exec
UPDATE products
  SET name = $2,
  description = $3,
  price = $4,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- Product Images
-- name: AddProductImage :exec
INSERT INTO product_images (
  product_id, image_url
) VALUES (
  $1, $2
);

-- name: GetProductImages :many
SELECT image_url FROM product_images
WHERE product_id = $1;

-- name: DeleteProductImage :exec
DELETE FROM product_images
WHERE product_id = $1 AND image_url = $2;

-- Product Sizes
-- name: AddProductSize :exec
INSERT INTO product_sizes (
  product_id, size_name, stock
) VALUES (
  $1, $2, $3
);

-- name: GetProductSizes :many
SELECT size_name, stock FROM product_sizes
WHERE product_id = $1;

-- name: UpdateProductStock :exec
UPDATE product_sizes
  SET stock = $3,
  updated_at = NOW()
WHERE product_id = $1 AND size_name = $2;

-- Fit Guide
-- name: GetProductFitGuide :one
SELECT * FROM fit_guides
WHERE product_id = $1 LIMIT 1;

-- name: CreateProductFitGuide :exec
INSERT INTO fit_guides (
  product_id, body_length, sleeve_length, chest_width, shoulder_width,
  arm_hole, front_rise, inseam, hem, back_rise, waist, thigh, knee
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
);

-- name: UpdateProductFitGuide :exec
UPDATE fit_guides
  SET body_length = $2,
  sleeve_length = $3,
  chest_width = $4,
  shoulder_width = $5,
  arm_hole = $6,
  front_rise = $7,
  inseam = $8,
  hem = $9,
  back_rise = $10,
  waist = $11,
  thigh = $12,
  knee = $13,
  updated_at = NOW()
WHERE product_id = $1;

-- Carts
-- name: GetCart :one
SELECT * FROM carts
WHERE id = $1 LIMIT 1;

-- name: CreateCart :one
INSERT INTO carts (
  id
) VALUES (
  $1
)
RETURNING *;

-- name: UpdateCartTimestamp :exec
UPDATE carts
  SET updated_at = NOW()
WHERE id = $1;

-- Cart Items
-- name: GetCartItems :many
SELECT ci.product_id, ci.quantity, p.name, p.description, p.price
FROM cart_items ci
JOIN products p ON ci.product_id = p.id
WHERE ci.cart_id = $1;

-- name: AddCartItem :exec
INSERT INTO cart_items (
  cart_id, product_id, quantity
) VALUES (
  $1, $2, $3
)
ON CONFLICT (cart_id, product_id) 
DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity,
              updated_at = NOW();

-- name: UpdateCartItemQuantity :exec
UPDATE cart_items
  SET quantity = $3,
  updated_at = NOW()
WHERE cart_id = $1 AND product_id = $2;

-- name: RemoveCartItem :exec
DELETE FROM cart_items
WHERE cart_id = $1 AND product_id = $2;

-- name: ClearCart :exec
DELETE FROM cart_items
WHERE cart_id = $1;

-- name: GetCartItemCount :one
SELECT COUNT(*) FROM cart_items
WHERE cart_id = $1;

-- Collections
-- name: GetCollection :one
SELECT * FROM collections
WHERE id = $1 LIMIT 1;

-- name: ListCollections :many
SELECT * FROM collections
ORDER BY name;

-- name: CreateCollection :one
INSERT INTO collections (
  name, description
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateCollection :exec
UPDATE collections
  SET name = $2,
  description = $3,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteCollection :exec
DELETE FROM collections
WHERE id = $1;

-- Collection Images
-- name: AddCollectionImage :exec
INSERT INTO collection_images (
  collection_id, image_url
) VALUES (
  $1, $2
);

-- name: GetCollectionImages :many
SELECT image_url FROM collection_images
WHERE collection_id = $1;

-- name: DeleteCollectionImage :exec
DELETE FROM collection_images
WHERE collection_id = $1 AND image_url = $2;

-- Collection Products
-- name: AddProductToCollection :exec
INSERT INTO collection_products (
  collection_id, product_id
) VALUES (
  $1, $2
)
ON CONFLICT (collection_id, product_id) DO NOTHING;

-- name: RemoveProductFromCollection :exec
DELETE FROM collection_products
WHERE collection_id = $1 AND product_id = $2;

-- name: GetCollectionProducts :many
SELECT p.* FROM products p
JOIN collection_products cp ON p.id = cp.product_id
WHERE cp.collection_id = $1;

-- name: GetUserByEmailAndPassword :one
SELECT * FROM users
WHERE email = $1 and password_hash = $2 LIMIT 1;

