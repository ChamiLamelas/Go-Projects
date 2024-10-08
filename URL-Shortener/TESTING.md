# Testing 

## Overview 

This file outlines the test cases in the [tests folder](./tests/). The test cases were used in testing the implementation. Most of these test cases were derived from [DESIGN.md](./DESIGN.md).

The tests generally work by sending HTTP requests with `curl` to the server and collecting the output (HTTP response code, response body) into `.out` files. The output is then diffed with a saved reference output file. 

> Note: throughout this file, we assume that the current working directory is `tests`.

## Files 

- `boot.sh` is used to start the server without touching the database file if one exists.
- `fresh_boot.sh` is used to wipe the database and then start the server with a fresh database.
- `testXx.ref` are the reference output files.
- `testXx.sh` are the test scripts to be run representing the client. For tests that passed, these should be empty.

## Test Cases

### Test 1

**Description:** Simple behavior test where we check if an alias can be automatically assigned.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test1.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 2

**Description:** check if an alias can be custom assigned.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test2.sh` in a second terminal.
3. `Ctrl + C` the server. 

### Test 3

**Description:** check if an alias can be automatically assigned and then expanded.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test3.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 4

**Description:** check if an alias can be custom assigned and then expanded.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test4.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 5

**Description:** check if an alias can be automatically assigned, expanded, and then analytics can be ran.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test5.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 6

**Description:** check if two aliases are properly created.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test6.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 7

**Description:** check if two aliases are properly created between server boots.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test7a.sh` in a second terminal.
3. `Ctrl + C` the server.
4. Run `bash boot.sh` in the first terminal.
5. Run `bash test7b.sh` in the second terminal.
6. `Ctrl + C` the server.

### Test 8

**Description:** check if server properly handles a user trying to make a second automatic alias for the same URL.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test8.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 9

**Description:** check if server properly handles a user trying to make a second custom alias for the same URL.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test9.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 10

**Description:** check if server properly handles a user trying to use the same custom alias for two URLs.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test10.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 11

**Description:** check if server properly rejects requests with wrong method for `shorten/` route.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test11.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 12 

**Description:** check if server properly rejects requests with wrong method for `expand/` route. 

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test12.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 13

**Description:** check if server properly rejects requests with wrong method for `analytics/` route.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test13.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 14

**Description:** check if server properly rejects unknown routes.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test14.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 15

**Description:** check if server properly handles a user trying to expand an alias that does not exist.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test15.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 16

**Description:** check if server properly handles a user trying to get analytics on an alias that does not exist.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test16.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 17

**Description:** check if server properly handles a user trying to get analytics on an alias that was never expanded.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test17.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 18

**Description:** check if server properly handles when a user has made a custom alias that happens to match the next automatic aliases.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test18.sh` in a second terminal.
3. `Ctrl + C` the server.

### Test 19

**Description:** check if custom alias is created in one boot and then in the second boot, it can be expanded.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test19a.sh` in a second terminal.
3. `Ctrl + C` the server.
4. Run `bash boot.sh` in the first terminal.
5. Run `bash test19b.sh` in the second terminal.
6. `Ctrl + C` the server.

### Test 20

**Description:** check if automatic and custom aliases can be created over multiple boots, can be expanded, and analytics are maintained over multiple boots.

1. Run `bash fresh_boot.sh` in one terminal.
2. Run `bash test20a.sh` in a second terminal.
3. `Ctrl + C` the server.
4. Run `bash boot.sh` in the first terminal.
5. Run `bash test20b.sh` in the second terminal.
6. `Ctrl + C` the server.
7. Run `bash boot.sh` in the first terminal.
8. Run `bash test20c.sh` in the second terminal.
9. `Ctrl + C` the server.

### Test 21

**Description:** check if back to back `shorten/` requests are properly handled. This is meant to test for race conditions in how aliases are automatically assigned. This test is handled a bit differently, and it involves checking things manually. 

1. Uncomment the two logging statements at the top of `ShortenAutomatic` in `server.go`.
2. Run `bash fresh_boot.sh` in one terminal.
3. Run `bash test21.sh` in a second terminal.
4. `Ctrl + C` the server.
5. First, check the log statements in the server terminal. You will see that shorten requests are being handled concurrently. For example, it may look something like this: 

    ```text
    2024/08/27 12:34:50.319581 Beginning to service shorten request for https://www.web2.com
    2024/08/27 12:34:50.319581 Beginning to service shorten request for https://www.web4.com
    2024/08/27 12:34:50.319581 Beginning to service shorten request for https://www.web1.com
    2024/08/27 12:34:50.326429 Finished servicing shorten request for https://www.web2.com
    2024/08/27 12:34:50.332235 Finished servicing shorten request for https://www.web4.com
    2024/08/27 12:34:50.334306 Beginning to service shorten request for https://www.web3.com
    2024/08/27 12:34:50.334306 Beginning to service shorten request for https://www.web5.com
    2024/08/27 12:34:50.338010 Finished servicing shorten request for https://www.web1.com
    2024/08/27 12:34:50.342716 Finished servicing shorten request for https://www.web3.com
    2024/08/27 12:34:50.348275 Finished servicing shorten request for https://www.web5.com
    ```

    Then, check each of the `.out` files produced for this test. The alias should be unique for each going from 0 to 4. The order of assignment should match the "finished" prints in the server print log. 

    For example, with the output above, the web2 gets alias 0, web4 gets alias 1, web1 gets alias 2, web3 gets alias 3, and web5 gets alias 4.
