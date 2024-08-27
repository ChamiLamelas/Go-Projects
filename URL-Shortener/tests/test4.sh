curl -isSX POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com", "alias":"google"}' > test4.out 2>&1
curl -isSX GET http://localhost:8000/urlshortener/expand/google >> test4.out 2>&1
diff test4.out test4.ref