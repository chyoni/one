# ONE COIN

### Index

- #01 Init

- #02 Block, Chain

- #03 bbolt DB and persist

- #04 Restore blockchain

  > gob encode, decode

- #05 Restore blocks

- #06 REST API Server Part 1

- #07 REST API Server Part 2 (all blocks)

- #08 REST API Server Part 3 (a block)

- #09 REST API Server Part 4 (middleware)

- #10 REST API Server Part 5 (add block)

- #11 REST API Server Part 6 (blockchain status)

- #12 Mempool and transaction

- #13 REST API Server Part 7 (transaction)

- #14 Mempool database

- #15 Mining Part 1

- #16 Mining Part 2

- #17 Mining DONE

- #18 Unspend Transaction Output

- #19 isOnMempool Part 1

- #20 isOnMempool DONE

- #21 CLI

- #22 Wallet Part 1 (generate privateKey)

- #23 Wallet Part 2 (Sign, Verify func)

- #24 Wallet Part DONE

- #25 P2P Part 1 (websocket upgrade and connection between node to node)

- #26 P2P Part 2 (who connected with me)

- #27 P2P Part 3 (read and write)

- #28 P2P Part 4 (send message and seperate database by node)

- #29 P2P Part 5 (when disconnection occur, delete peer and fix data races)

- #30 P2P Part 6 (request all blocks message and send all blocks message)

- #31 P2P Part 7 (synchronized blockchain between peer to peer)

- #32 P2P Part 8 (Clearing data races)

- #33 P2P Part 9 (New block broadcast to peers)

- #34 P2P Part 10 (New transaction broadcast to peers)

- #35 P2P Part 11 (synchronized mempool between nodes)

- #36 P2P Part 12 (broadcast new peer)

- #37 P2P Part DONE! (clearing data races)

- #38 Unit Test Part 1 (blockchain package)

  #### TDD with interface

  - 테스트를 할 때 실제 Database에서 가져오는 데이터를 사용하고 싶지 않고 그저 코드의 로직을 테스트하고 싶기 때문에 실제로 돌릴 때 사용될 function과 테스트할 때 사용될 function을
    구분해야 한다. 이를 위해 interface로 해당 func을 구현하고 해당 interface를 implement하는 두 struct를 사용하여 하나는 실제 환경 나머지 하나는 테스트 환경을 위해 사용하게끔 코드를 작성한다.

  #### Go Test Command

  - go test -v ./.../ -> 현재 위치에서 모든 폴더의 \_test.go 형식의 test파일을 실행
  - go test -v ./blockchain -> 현재 위치에서 blockchain폴더의 \_test.go 형식의 test파일을 실행

- #39 Unit Test Part 2 (blockchain package)

- #40 Unit Test Part 3 (blockchain package)

  - go test -v -coverprofile cover.out ./... && go tool cover -html=cover.out

- #41 Unit Test Part 4 (blockchain package)

  - test for infinite loop

- #42 Unit Test Part 5 (blockchain package)

- #43 Unit Test Part 6 (blockchain package)

- #44 Unit Test Part 7 (blockchain package)

- #45 Unit Test Part 8 (blockchain package)

- #46 Unit Test Part 9 (blockchain package)
