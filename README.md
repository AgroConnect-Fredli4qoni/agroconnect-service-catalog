# Repositori: agroconnect-service-catalog
## Microservice Manajemen Komoditas Hasil Panen (Golang + MongoDB)
**Skema Sertifikasi**: BNSP Full-Stack Developer | **Kandidat**: Fredli Fourqoni

---

### Deskripsi
Layanan Service Catalog bertanggung jawab atas pengelolaan data hasil panen pertanian, komoditas holtikultura, varietas bibit, dan pupuk menggunakan basis data dokumen **MongoDB**.

### Fitur Utama:
1. **CRUD Komoditas**: Tambah, ubah, hapus, dan lihat detail komoditas tani.
2. **Filter & Pencarian**: Filter kategori (Pangan Pokok, Sayur, Buah, Bumbu) dan pencarian nama hasil tani.
3. **Manajemen Stok Real-Time**: Sinkronisasi penyesuaian stok saat transaksi berhasil dibuat oleh Service Order.
4. **Health Check**: Endpoint `/health` untuk monitoring ketersediaan kontainer.

### Port
- Port Eksternal/Internal: `8081`
