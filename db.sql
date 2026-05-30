CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE dokter (
	id_dokter UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	nama_dokter VARCHAR(100) NOT NULL,
	no_hp VARCHAR(20) UNIQUE NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	alamat TEXT,
	klinik VARCHAR(100),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE produk (
	id_produk UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	nama_bahan VARCHAR(100) NOT NULL,
	spesifikasi TEXT,
	harga NUMERIC(12,2) NOT NULL CHECK (harga>0),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE karyawan (
	id_karyawan  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	nama VARCHAR(255) NOT NULL,
	no_hp VARCHAR(20) UNIQUE NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	role VARCHAR(20) NOT NULL CHECK(role IN('cs', 'teknisi', 'bos')),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE pesanan(
	id_pesanan UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id_dokter UUID NOT NULL REFERENCES dokter(id_dokter),
	status_pesanan VARCHAR(30) NOT NULL DEFAULT 'menunggu'
		CHECK(status_pesanan IN(
		'menunggu', 'disetujui', 'ditolak', 'produksi', 'selesai', 'revisi'
		)),
	tgl_pesanan TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE detail_pesanan(
	id_detail UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id_pesanan UUID NOT NULL REFERENCES pesanan(id_pesanan) ON DELETE CASCADE,
	id_produk UUID NOT NULL REFERENCES produk(id_produk),
	warna VARCHAR(20) NOT NULL,
	ukuran VARCHAR(10) NOT NULL,
	jumlah INT NOT NULL DEFAULT 1 CHECK (jumlah>0),
	harga_satuan NUMERIC(12,2) NOT NULL,
	subtotal NUMERIC(12,2) GENERATED ALWAYS AS (jumlah * harga_satuan) STORED,
	catatan_tambahan TEXT NULL,
	kode_gigi VARCHAR(10) NOT NULL
);

CREATE TABLE persetujuan (
	id_persetujuan UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id_pesanan UUID NOT NULL REFERENCES pesanan(id_pesanan) ON DELETE CASCADE,
	id_karyawan UUID NOT NULL REFERENCES karyawan(id_karyawan),
	status VARCHAR(20) NOT NULL CHECK (status IN('disetujui', 'ditolak', 'revisi')),
	catatan_bos TEXT,
	tgl_persetujuan TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE revisi(
	id_revisi  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id_pesanan UUID NOT NULL REFERENCES pesanan(id_pesanan) ON DELETE CASCADE,
	status_revisi VARCHAR(20) NOT NULL CHECK(status_revisi IN('pending', 'dikerjakan', 'selesai')),
	deskripsi_revisi TEXT NOT NULL,
	tgl_revisi TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pengerjaan(
	id_pengerjaan  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id_pesanan UUID NOT NULL REFERENCES pesanan(id_pesanan) ON DELETE CASCADE,
	id_karyawan UUID NOT NULL REFERENCES karyawan(id_karyawan),
	id_revisi UUID NULL REFERENCES revisi(id_revisi) ON DELETE CASCADE,
	status_produksi VARCHAR(20) NOT NULL DEFAULT 'antrian' CHECK (status_produksi IN('antrian', 'dikerjakan', 'revisi', 'selesai')),
	catatan_karyawan TEXT,
	tgl_mulai TIMESTAMP,
	tgl_selesai TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pengiriman(
	id_pengiriman  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id_pesanan UUID NOT NULL REFERENCES pesanan(id_pesanan) ON DELETE CASCADE,
	nama_jasa VARCHAR(50) NOT NULL,
	no_resi VARCHAR(100) NOT NULL,
	tgl_kirim TIMESTAMP,
	tgl_diterima TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE detail_pengiriman(
	id_detail_pengiriman  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id_pengiriman UUID NOT NULL REFERENCES pengiriman(id_pengiriman) ON DELETE CASCADE,
	status VARCHAR(50) NOT NULL DEFAULT 'Menunggu' CHECK(status IN('Menunggu', 'Dijemput Kurir', 'Dalam Pengiriman', 'Diterima')),
	waktu TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	keterangan TEXT
);

CREATE TABLE transaksi(
	id_transaksi  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id_pesanan UUID NOT NULL REFERENCES pesanan(id_pesanan) ON DELETE CASCADE,
	id_karyawan UUID NOT NULL REFERENCES karyawan(id_karyawan),
	total_harga NUMERIC(12,2) NOT NULL,
	metode_pembayaran VARCHAR(30) CHECK(metode_pembayaran IN('transfer', 'tunai', 'qris', 'gopay')),
	status_pembayaran VARCHAR(30) DEFAULT 'belum bayar' CHECK (status_pembayaran IN('lunas', 'belum bayar')),
	tgl_transaksi TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pesanan_dokter ON pesanan(id_dokter);
CREATE INDEX idx_pesanan_status ON pesanan(status_pesanan);
CREATE INDEX idx_pengerjaan_status ON pengerjaan(status_produksi);
CREATE INDEX idx_transaksi_status ON transaksi(status_pembayaran);

select * from karyawan;
