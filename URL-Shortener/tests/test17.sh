curl -s -w "\nResponse code: %{http_code}" -X POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com"}' > test17.out 2>&1
curl -s -w "\nResponse code: %{http_code}" -X GET http://localhost:8000/urlshortener/analytics/0 >> test17.out 2>&1
diff test17.out test17.ref