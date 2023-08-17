# About
Aplikasi ini adalah microservice Microsite Refund API yang memiliki 3 endpoint:
1. /m/select-store
2. /m/store-state
3. /m/v4/nearest-store

# How to run
1. lakukan migrasi database dari repository chat-commerce-backend (karena microservice ini tidak memiliki migrasi yang terpisah)
2. siapkan application.yaml yang dibutuhkan
3. jalankan unit test untuk memastikan flow aplikasi sesuai desain. Setelah passed, dapat jalankan aplikasi.
```bash
cd dev && go test ./...
cd ../bin && ./mi_storesapi serve
```