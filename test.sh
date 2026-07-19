#!/bin/bash

# Somr Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_URL="http://localhost:8080"

echo -e "${BLUE}=== Launching Automated AeroForge Tests ===${NC}\n"

# 1. Test GET /welcome
echo -e "${BLUE}[TEST 1] GET /welcome?user=Charles${NC}"
RESPONSE=$(curl -s "$BASE_URL/welcome?user=Charles")
echo "Response: $RESPONSE"
if [[ $RESPONSE == *"Welcome to AeroForge Control, Charles!"* ]]; then
    echo -e "${GREEN}-> SUCCESS${NC}\n"
else
    echo -e "${RED}-> FAILURE${NC}\n"
fi

# 2. Test GET /parts
echo -e "${BLUE}[TEST 2] GET /parts${NC}"
RESPONSE=$(curl -s "$BASE_URL/parts")
echo "Response: $RESPONSE"
if [[ $RESPONSE == *"Ariane6_Booster"* ]]; then
    echo -e "${GREEN}-> SUCCESS${NC}\n"
else
    echo -e "${RED}-> FAILURE${NC}\n"
fi

# 3. Test GET /telemetry (ISS)
echo -e "${BLUE}[TEST 3] GET /telemetry (Live ISS Tracking)${NC}"
RESPONSE=$(curl -s "$BASE_URL/telemetry")
echo "Response: $RESPONSE"
if [[ $RESPONSE == *"NOMINAL"* ]]; then
    echo -e "${GREEN}-> SUCCESS${NC}\n"
else
    echo -e "${RED}-> FAILURE (External telemetry API unreachable)${NC}\n"
fi

# 4. Test POST /validate (Valid Payload)
echo -e "${BLUE}[TEST 4] POST /validate (Valid Payload)${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/validate" \
  -H "Content-Type: application/json" \
  -d '{"name":"Ariane6_Booster","geometry":"half-cone-rectangular"}')
echo "Response: $RESPONSE"
if [[ $RESPONSE == *"APPROVED"* ]]; then
    echo -e "${GREEN}-> SUCCESS${NC}\n"
else
    echo -e "${RED}-> FAILURE${NC}\n"
fi

# 5. Test POST /validate (Invalid Payload - Wrong Geometry)
echo -e "${BLUE}[TEST 5] POST /validate (Invalid Payload - Wrong Geometry)${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/validate" \
  -H "Content-Type: application/json" \
  -d '{"name":"Ariane6_Booster","geometry":"wrong-geometry"}')
echo "Response: $RESPONSE"
if [[ $RESPONSE == *"CRITICAL ERROR"* ]]; then
    echo -e "${GREEN}-> SUCCESS (Error successfully caught)${NC}\n"
else
    echo -e "${RED}-> FAILURE (API should have rejected the part)${NC}\n"
fi

echo -e "${BLUE}=== End of Test Suite ===${NC}"