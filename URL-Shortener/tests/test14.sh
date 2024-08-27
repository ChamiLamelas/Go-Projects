curl -s -w "\nResponse code: %{http_code}\n" -X GET http://localhost:8000/urlshortener/blah/0 > test14.out 2>&1
diff test14.out test14.ref