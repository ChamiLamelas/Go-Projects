curl -isSX POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com"}' > test3.out 2>&1
curl -isSX GET http://localhost:8000/urlshortener/expand/000000000000 >> test3.out 2>&1
diff test3.out test3.ref