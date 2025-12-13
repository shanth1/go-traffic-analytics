#!/bin/bash

source .env

BASE_URL="http://localhost:8080"
ADMIN_EMAIL="admin@gotrace.com"
CLIENT_EMAIL="client@gotrace.com"
PASSWORD="password"

GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting E2E Test Suite for GoTrace${NC}\n"

check_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✔ Success${NC}"
    else
        echo -e "${RED}✘ Failed${NC}"
        exit 1
    fi
}

# ==============================================================================
# 1. PUBLIC ROUTES
# ==============================================================================
echo -e "${BLUE}--- [1] Testing Public Redirects ---${NC}"

echo -n "Redirecting /google (Expect 307)... "
CODE=$(curl -o /dev/null -s -w "%{http_code}" "$BASE_URL/google")
if [ "$CODE" == "307" ]; then echo -e "${GREEN}✔ OK ($CODE)${NC}"; else echo -e "${RED}✘ Fail ($CODE)${NC}"; fi

echo -n "Redirecting /unknown (Expect 404)... "
CODE=$(curl -o /dev/null -s -w "%{http_code}" "$BASE_URL/unknown-slug-123")
if [ "$CODE" == "404" ]; then echo -e "${GREEN}✔ OK ($CODE)${NC}"; else echo -e "${RED}✘ Fail ($CODE)${NC}"; fi


# ==============================================================================
# 2. AUTHENTICATION
# ==============================================================================
echo -e "\n${BLUE}--- [2] Testing Auth ---${NC}"

# 2.1 Login Client
echo -n "Logging in as Client... "
CLIENT_RES=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$CLIENT_EMAIL\", \"password\":\"$PASSWORD\"}")
CLIENT_TOKEN=$(echo $CLIENT_RES | jq -r .token)

if [ "$CLIENT_TOKEN" != "null" ]; then echo -e "${GREEN}✔ Token received${NC}"; else echo -e "${RED}✘ Failed${NC}"; exit 1; fi

# 2.2 Login Admin
echo -n "Logging in as Admin... "
ADMIN_RES=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$ADMIN_EMAIL\", \"password\":\"$PASSWORD\"}")
ADMIN_TOKEN=$(echo $ADMIN_RES | jq -r .token)
if [ "$ADMIN_TOKEN" != "null" ]; then echo -e "${GREEN}✔ Token received${NC}"; else echo -e "${RED}✘ Failed${NC}"; exit 1; fi

# 2.3 Register New User
NEW_EMAIL="newuser_$(date +%s)@test.com"
echo -n "Registering new user ($NEW_EMAIL)... "
REG_RES=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -H "X-API-Key: ${API_KEY}" \
    -d "{\"email\":\"$NEW_EMAIL\", \"password\":\"password\"}")
NEW_USER_ID=$(echo $REG_RES | jq -r .id)

if [ "$NEW_USER_ID" != "null" ]; then
    echo -e "${GREEN}✔ Created (ID: $NEW_USER_ID)${NC}"
else
    echo -e "${RED}✘ Failed ($REG_RES)${NC}"
fi


# ==============================================================================
# 3. CLIENT FLOWS (Campaigns & Links)
# ==============================================================================
echo -e "\n${BLUE}--- [3] Testing Client Features ---${NC}"

# 3.1 Create Campaign
echo -n "Creating Campaign 'Marketing 2025'... "
CAMP_RES=$(curl -s -X POST "$BASE_URL/api/v1/campaigns" \
    -H "Authorization: Bearer $CLIENT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name": "Marketing 2025"}')
CAMP_ID=$(echo $CAMP_RES | jq -r .data.id)
if [ "$CAMP_ID" != "null" ]; then echo -e "${GREEN}✔ OK (ID: $CAMP_ID)${NC}"; else echo -e "${RED}✘ Fail${NC}"; exit 1; fi

# 3.2 List Campaigns
echo -n "Listing Campaigns... "
curl -s -f -H "Authorization: Bearer $CLIENT_TOKEN" "$BASE_URL/api/v1/campaigns" > /dev/null
check_status $?

# 3.3 Create Link IN Campaign
SLUG="promo-$(date +%s)"
echo -n "Creating Link inside Campaign (slug: $SLUG)... "
LINK_RES=$(curl -s -X POST "$BASE_URL/api/v1/links" \
    -H "Authorization: Bearer $CLIENT_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"campaign_id\": \"$CAMP_ID\", \"target_url\": \"https://example.com\", \"custom_slug\": \"$SLUG\"}")
LINK_ID=$(echo $LINK_RES | jq -r .data.id)
if [ "$LINK_ID" != "null" ]; then echo -e "${GREEN}✔ OK (ID: $LINK_ID)${NC}"; else echo -e "${RED}✘ Fail${NC}"; exit 1; fi

# 3.4 Check created link works
echo -n "Checking redirect for new link /$SLUG... "
CODE=$(curl -o /dev/null -s -w "%{http_code}" "$BASE_URL/$SLUG")
if [ "$CODE" == "307" ]; then echo -e "${GREEN}✔ OK${NC}"; else echo -e "${RED}✘ Fail ($CODE)${NC}"; fi

# 3.5 List Links by Campaign
echo -n "Get Links by Campaign... "
LINKS_LIST=$(curl -s -H "Authorization: Bearer $CLIENT_TOKEN" "$BASE_URL/api/v1/campaigns/$CAMP_ID/links")
COUNT=$(echo $LINKS_LIST | jq '.data | length')
if [ "$COUNT" -gt 0 ]; then echo -e "${GREEN}✔ OK (Count: $COUNT)${NC}"; else echo -e "${RED}✘ Fail (Empty)${NC}"; fi


# ==============================================================================
# 4. ANALYTICS
# ==============================================================================
echo -e "\n${BLUE}--- [4] Testing Analytics Endpoints ---${NC}"

endpoints=("summary" "stream?group_by=os" "flow" "geo" "quality")

for ep in "${endpoints[@]}"; do
    echo -n "GET /analytics/$ep ... "
    curl -s -f -H "Authorization: Bearer $CLIENT_TOKEN" "$BASE_URL/api/v1/analytics/$ep" > /dev/null
    if [ $? -eq 0 ]; then echo -e "${GREEN}✔ OK${NC}"; else echo -e "${RED}✘ Fail${NC}"; fi
done


# ==============================================================================
# 5. ADMIN FLOWS
# ==============================================================================
echo -e "\n${BLUE}--- [5] Testing Admin Features ---${NC}"

# 5.1 Get Users
echo -n "Admin: Get Users... "
USERS_LIST=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/api/v1/admin/users")
# Находим ID нашего нового пользователя из шага 2.3
FOUND_ID=$(echo $USERS_LIST | jq -r ".data[] | select(.email==\"$NEW_EMAIL\") | .id")

if [ "$FOUND_ID" == "$NEW_USER_ID" ]; then
    echo -e "${GREEN}✔ OK (Found new user in list)${NC}"
else
    echo -e "${RED}✘ Fail (User not found in list)${NC}"
    exit 1
fi

# 5.2 Ban User
echo -n "Admin: Banning User $NEW_EMAIL... "
curl -s -X PATCH "$BASE_URL/api/v1/admin/users/$NEW_USER_ID/status" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"is_active": false}' > /dev/null
check_status $?

# 5.3 Verify Ban (Login attempt)
echo -n "Verify Ban (Login should fail)... "
BAN_LOGIN=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$NEW_EMAIL\", \"password\":\"password\"}")

if [ "$BAN_LOGIN" == "401" ] || [ "$BAN_LOGIN" == "403" ]; then
    echo -e "${GREEN}✔ OK (Access Denied)${NC}"
else
    echo -e "${RED}✘ Fail (Got code $BAN_LOGIN)${NC}"
fi

# 5.4 Change Plan
echo -n "Admin: Upgrade User Plan to 'enterprise'... "
curl -s -X PATCH "$BASE_URL/api/v1/admin/users/$NEW_USER_ID/plan" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"plan_id": "enterprise"}' > /dev/null
check_status $?


# ==============================================================================
# 6. CLEANUP
# ==============================================================================
echo -e "\n${BLUE}--- [6] Cleanup ---${NC}"

echo -n "Deleting Link $LINK_ID... "
curl -s -X DELETE "$BASE_URL/api/v1/links/$LINK_ID" -H "Authorization: Bearer $CLIENT_TOKEN"
check_status $?

echo -e "\n${BLUE}🎉 All Tests Completed!${NC}"
