#!/usr/bin/env bash
set -e
API_URL="http://localhost:8888"

# Test /greeting/{name}
echo "Testing /greeting/{name}..."
curl -s "$API_URL/greeting/world" | grep 'Hello'

# Test /examples POST
echo "Testing POST /examples..."
EXAMPLE_JSON=$(curl -s -X POST "$API_URL/examples" -H 'Content-Type: application/json' -d '{"id":"test1","name":"Test Example"}')
EXAMPLE_ID=$(echo "$EXAMPLE_JSON" | jq -r .id)
if [[ "$EXAMPLE_ID" == "null" || -z "$EXAMPLE_ID" ]]; then
	echo "POST /examples failed: no id returned"
	exit 1
fi

# Test /examples GET
echo "Testing GET /examples..."
curl -s "$API_URL/examples" | grep 'Test Example'

# Test /examples/{id} GET
echo "Testing GET /examples/{id}..."
curl -s "$API_URL/examples/$EXAMPLE_ID" | grep 'Test Example'

# Test /examples/{id} PUT
echo "Testing PUT /examples/{id}..."
curl -s -X PUT "$API_URL/examples/$EXAMPLE_ID" -H 'Content-Type: application/json' -d '{"id":"$EXAMPLE_ID","name":"Updated Example"}' | grep 'Updated Example'

# Test /examples/{id} PATCH
echo "Testing PATCH /examples/{id}..."
curl -s -X PATCH "$API_URL/examples/$EXAMPLE_ID" -H 'Content-Type: application/json' -d '{"name":"Patched Example"}' | grep 'Patched Example'

# Test /examples/{id} DELETE
echo "Testing DELETE /examples/{id}..."
curl -s -X DELETE "$API_URL/examples/$EXAMPLE_ID"

# Confirm deletion
echo "Confirming deletion..."
if curl -s "$API_URL/examples/$EXAMPLE_ID" | grep 'Patched Example'; then
	echo "DELETE /examples/{id} failed"
	exit 1
fi

echo "All endpoint tests passed!"
