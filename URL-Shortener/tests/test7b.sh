curl -isSX POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.nytimes.com"}' >> test7.out 2>&1
diff test7.out test7.ref