curl -s -w "\nResponse code: %{http_code}\n" -X DELETE http://localhost:8000/urlshortener/expand/0 > test12.out 2>&1
diff test12.out test12.ref