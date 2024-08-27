curl -s -w "\nResponse code: %{http_code}" -X POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com"}' > test5.out 2>&1
curl -s -w "\nResponse code: %{http_code}" -X GET http://localhost:8000/urlshortener/expand/0 >> test5.out 2>&1
curl -s -w "\nResponse code: %{http_code}" -X GET http://localhost:8000/urlshortener/analytics/0 >> test5.out 2>&1
diff test5.out test5.ref