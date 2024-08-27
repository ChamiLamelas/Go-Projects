curl -s -w "\nResponse code: %{http_code}" -X POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com"}' > test6.out 2>&1
curl -s -w "\nResponse code: %{http_code}" -X POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.nytimes.com"}' >> test6.out 2>&1
diff test6.out test6.ref