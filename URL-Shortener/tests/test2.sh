curl -s -w "\nResponse code: %{http_code}\n" -X POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com", "alias":"google"}' > test2.out 2>&1
diff test2.out test2.ref