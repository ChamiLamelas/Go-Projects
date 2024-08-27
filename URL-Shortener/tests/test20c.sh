curl -isSX GET http://localhost:8000/urlshortener/expand/000000000000 >> test20.out 2>&1
curl -isSX GET http://localhost:8000/urlshortener/expand/000000000001 >> test20.out 2>&1
curl -isSX POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.web3.com"}' >> test20.out 2>&1
curl -isSX GET http://localhost:8000/urlshortener/analytics/000000000000 >> test20.out 2>&1
curl -isSX GET http://localhost:8000/urlshortener/analytics/000000000001 >> test20.out 2>&1
curl -isSX GET http://localhost:8000/urlshortener/analytics/000000000002 >> test20.out 2>&1
diff test20.out test20.ref