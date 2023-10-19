$BODY = '{ "id": 133, "weight": 2.9, "name": "Bosch", "length": 3.4, "worktime": 12, "diameter": 15 }'

$response = curl -d "$BODY" localhost:4000/v1/drills

$response