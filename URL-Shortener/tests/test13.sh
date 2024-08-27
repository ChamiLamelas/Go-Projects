curl -s -w "\nResponse code: %{http_code}\n" -X DELETE http://localhost:8000/urlshortener/analytics/0 > test13.out 2>&1
diff test13.out test13.ref