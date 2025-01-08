# Personal Project SIAKAD24-Backend
> Backend untuk project lawas UKOM jaman SMK SIAKAD 24, menerapkan microservice, front-end menggunakan laravel + React.js nantinya

## How to Run
1. Install dependensi pada local
2. Jalankan program dengan ``go run cmd/api/main.go``

### Auto Live-Reload (Opsional)
> Digunakan untuk live-reload project tanpa perlu restart manual command runningnya
> Untuk lebih jelasnya bisa kunjungi link ini untuk reflex https://elpahlevi.medium.com/menerapkan-live-reloading-di-golang-b9260f795b00, atau bisa menggunakan Air jika OS yg anda gunakan adalah windows (bisa juga pada linux, macos)

#### INSTALL PACKAGE REFLEX
1. Install package ``reflex`` jika belum ada > ``go install github.com/cespare/reflex@latest``
2. Check package sudah terinstall secara proper atau belum ``reflex -h``
3. Jalankan reflex untuk live-reload ``reflex -s -r '\.go$$' go run cmd/api/main.go``

#### INSTALL PACKAGE AIR
1. ``go install github.com/air-verse/air@latest``
2. buat file konfigurasi ``.air.toml`` pada root workdir anda atau bisa dengan menjalankan command berikut ``air init``
3. ubah konfigurasi file ``main.go`` yg akan anda jalankan sesuai dengan project anda pada file ``.air.toml``, berikut konfigurasi nya
```
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

# Membunuh proses lama sebelum menjalankan ulang aplikasi
kill_before = true

[build]
  # Perintah untuk build aplikasi
  cmd = "go build -o tmp/main.exe cmd/api/main.go"
  bin = "tmp/main.exe"  # Menyimpan binary hasil build di direktori tmp/
  
  # Opsi delay build (dalam milidetik)
  delay = 1000
  
  # Direktori yang dikecualikan dari pemantauan
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  
  # Ekstensi file yang dipantau untuk perubahan
  include_ext = ["go", "tpl", "tmpl", "html"]
  
  # Jika ada perubahan pada file-file ini, akan membangun ulang aplikasi
  include_file = []
  
  # Tidak mengubah pengaturan berikut jika tidak diperlukan
  follow_symlink = false
  full_bin = ""
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  post_cmd = []
  pre_cmd = []
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  silent = false
  time = false

[misc]
  clean_on_exit = false

[proxy]
  app_port = 0
  enabled = false
  proxy_port = 0

[screen]
  clear_on_rebuild = false
  keep_scroll = true

[watch_dir]
  paths = ["cmd", "internal", "internal/handlers", "api", "tools"]

```  
4. jalankan live-reload menggunakan command ``air``