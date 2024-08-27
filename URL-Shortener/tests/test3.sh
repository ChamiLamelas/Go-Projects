curl -s -w "\nResponse code: %{http_code}\n" -X POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com"}' > test3.out 2>&1
curl -s -w "\nResponse code: %{http_code}\n" -X GET http://localhost:8000/urlshortener/expand/0 >> test3.out 2>&1
diff test3.out test3.ref