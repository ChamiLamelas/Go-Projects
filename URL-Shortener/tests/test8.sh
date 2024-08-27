curl -s -w "\nResponse code: %{http_code}\n" -X POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com"}' > test8.out 2>&1
curl -s -w "\nResponse code: %{http_code}\n" -X POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com"}' >> test8.out 2>&1
diff test8.out test8.ref