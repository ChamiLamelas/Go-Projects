curl -s -w "\nResponse code: %{http_code}\n" -X GET http://localhost:8000/urlshortener/expand/0 > test15.out 2>&1
diff test15.out test15.ref