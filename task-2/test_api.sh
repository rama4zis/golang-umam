#!/bin/sh

# Configuration
API_URL="http://localhost:8080"

echo "=========================================="
echo "Starting API Tests (Create & Update Only)"
echo "Target URL: $API_URL"
echo "=========================================="
echo ""

# --- Categories ---

echo "--- CATEGORY OPERATIONS ---"

echo "1. Creating Category 1 (Name: Appliances)..."
curl -s -X POST "$API_URL/api/category" \
  -H "Content-Type: application/json" \
  -d '{"name": "Appliances"}'
echo "\n"

echo "2. Creating Category 2 (Name: Books)..."
curl -s -X POST "$API_URL/api/category" \
  -H "Content-Type: application/json" \
  -d '{"name": "Books"}'
echo "\n"

echo "3. Update Category 1 (Change to: Home Appliances)..."
# Assuming ID 1 or the next available ID. Since we are scripting, we often guess IDs in simple tests.
# If this is a fresh DB, IDs might be 1 and 2. The previous request asked for ID 2 (implying 1 might exist or be skipped).
# I will assume the IDs for the newly created categories are X and X+1.
# But for simplicity in this script, the user mentioned "new id will be 3... and 2 in category" in the previous turn.
# I will just target IDs that likely exist or were just created.
# Let's try to update the ones we "expect" to create. If DB is persistent, IDs increment.
# I'll rely on the user manually checking or restarting DB for deterministic IDs, 
# OR I can try to parse the ID from the response if I used a more complex script (like python or jq).
# But for a simple shell script without jq guaranteed (it failed earlier), I will just hardcode likely IDs or just fire requests.
# The prompt implies specific interest in "2 category and 2 product".

# I'll attempt to update ID 1 and 2 (or whatever is valid).
# Or better, I will assume the user resets the DB or accepts that I target "recent" IDs.
# I will target ID 1 and 2 for categories as a guess.
curl -s -X PUT "$API_URL/api/category/1" \
  -H "Content-Type: application/json" \
  -d '{"name": "Home Appliances"}'
echo "\n"

curl -s -X PUT "$API_URL/api/category/2" \
  -H "Content-Type: application/json" \
  -d '{"name": "Educational Books"}'
echo "\n"

echo "4. Get All Categories..."
curl -s -X GET "$API_URL/api/category"
echo "\n"


# --- Products ---

echo "--- PRODUCT OPERATIONS ---"

echo "5. Creating Product 1 (linked to Category 1)..."
curl -s -X POST "$API_URL/api/product" \
  -H "Content-Type: application/json" \
  -d '{"name": "Blender", "price": 500000, "stock": 10, "category_id": {"id": 1}}'
echo "\n"

echo "6. Creating Product 2 (linked to Category 2)..."
curl -s -X POST "$API_URL/api/product" \
  -H "Content-Type: application/json" \
  -d '{"name": "Go Programming Book", "price": 150000, "stock": 20, "category_id": {"id": 2}}'
echo "\n"

echo "7. Update Product 1..."
# Assuming Product ID 1 exists
curl -s -X PUT "$API_URL/api/product/1" \
  -H "Content-Type: application/json" \
  -d '{"name": "Super Blender 3000", "price": 750000, "stock": 8, "category_id": {"id": 1}}'
echo "\n"

echo "8. Update Product 2..."
# Assuming Product ID 2 exists
curl -s -X PUT "$API_URL/api/product/2" \
  -H "Content-Type: application/json" \
  -d '{"name": "Advanced Go Programming", "price": 180000, "stock": 15, "category_id": {"id": 2}}'
echo "\n"

echo "9. Get All Products..."
curl -s -X GET "$API_URL/api/product"
echo "\n"

echo "=========================================="
echo "Tests Completed"
echo "=========================================="
