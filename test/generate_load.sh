#!/bin/bash

URL="http://localhost:8080/google"
COUNT=50

echo "🚀 Generating $COUNT clicks for $URL..."

for ((i=1;i<=COUNT;i++)); do
   UA="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/$((RANDOM % 10 + 90)).0.0.0 Safari/537.36"
   FAKE_IP="192.168.$((RANDOM % 255)).$((RANDOM % 255))"

   curl -s -o /dev/null -L \
     -H "User-Agent: $UA" \
     -H "X-Forwarded-For: $FAKE_IP" \
     "$URL"

   echo -n "."
done

echo -e "\n✅ Done! Check http://localhost:8080/api/v1/analytics/summary"
