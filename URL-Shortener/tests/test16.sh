curl -s -w "\nResponse code: %{http_code}\n" -X GET http://localhost:8000/urlshortener/analytics/0 > test16.out 2>&1
diff test16.out test16.ref