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


INSERT INTO produk (id_produk, nama_bahan, spesifikasi, harga) VALUES
('9e88d752-6b94-4d2b-b5d1-13ef2b1922c1', 'Crown PFM', 'Porcelain Fused to Metal standard', 500000),
('a3c91d8e-9087-43c2-a4e9-158a2d1847f2', 'Crown Zirconia', 'Full Zirconia Monolithic', 1200000),
('2d4a6f9c-76e3-4d82-8c10-09257e1ab321', 'E-Max Press', 'Lithium Disilicate Crown', 1500000),
('f1b8c2d5-e94a-48d6-a2b1-67e4f1a2c3d4', 'Gigi Tiruan Akrilik', 'Full Denture Akrilik (RA/RB)', 2000000),
('8c5d3a9b-1e24-4f5c-b8e7-2c6f1a3b4e5d', 'Valplast', 'Gigi Tiruan Fleksibel per elemen', 400000),
('b4a7f2e1-8d3c-49b5-a6e1-5f2d3c4b1e9a', 'Retainer Hawley', 'Orthodontic Retainer', 600000),
('1e9d2c4b-3a8f-46c1-b7d5-9e2a1f4c3b8d', 'Clear Aligner', 'Plastik PETG per set', 1000000),
('7f3a1b5c-4d2e-4b8a-9c1f-3e5d2a4b6c8f', 'Night Guard', 'Soft/Hard acrylic untuk bruxism', 750000),
('c2e4b6d8-1a3f-4e5c-a8d9-7b5c3e1a2f4d', 'Inlay/Onlay Composite', 'Komposit indirect', 450000),
('5d1a3c7b-8e2f-4a9b-c6d4-1f2e3a5b4c9d', 'Bridge PFM', 'Jembatan 3 unit PFM', 1500000);