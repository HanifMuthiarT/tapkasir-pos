"use client";

import { useEffect, useMemo, useRef, useState, type FormEvent } from "react";
import { useSearchParams } from "next/navigation";
import {
  ArrowLeft,
  ArrowRight,
  Banknote,
  Check,
  Coffee,
  CreditCard,
  Heart,
  MapPin,
  Minus,
  Plus,
  QrCode,
  Search,
  ShoppingBag,
  Trash2,
  X,
} from "lucide-react";
import {
  categories,
  menuItems,
  rupiah,
  type CategoryFilter,
  type MenuItem,
} from "../data/menu";

type Page = "beranda" | "favorit" | "pesanan";
type Drawer = "keranjang" | "checkout" | "preview" | null;
type Payment = "qris" | "transfer" | "cash";

export default function CustomerPage() {
  const params = useSearchParams();
  const meja = (params.get("meja") || "").trim().slice(0, 16);
  const [page, setPage] = useState<Page>("beranda");
  const [category, setCategory] = useState<CategoryFilter>("Semua Menu");
  const [searchOpen, setSearchOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [favorites, setFavorites] = useState<string[]>([]);
  const [cart, setCart] = useState<Record<string, number>>({});
  const [drawer, setDrawer] = useState<Drawer>(null);
  const [detail, setDetail] = useState<MenuItem | null>(null);
  const [detailQuantity, setDetailQuantity] = useState(1);
  const [toast, setToast] = useState<{ name: string; quantity: number } | null>(null);
  const [customerName, setCustomerName] = useState("");
  const [customerPhone, setCustomerPhone] = useState("");
  const [payment, setPayment] = useState<Payment>("qris");

  const categoryListRef = useRef<HTMLDivElement>(null);
  const [categoryScroll, setCategoryScroll] = useState({ left: false, right: false });

  // Tampilkan petunjuk geser hanya saat daftar kategori lebih lebar dari layar.
  useEffect(() => {
    const list = categoryListRef.current;
    if (!list) return;

    const updateScroll = () => {
      const maxScroll = list.scrollWidth - list.clientWidth;
      setCategoryScroll({
        left: list.scrollLeft > 6,
        right: maxScroll > 6 && list.scrollLeft < maxScroll - 6,
      });
    };

    updateScroll();
    list.addEventListener("scroll", updateScroll, { passive: true });
    window.addEventListener("resize", updateScroll);
    const observer = new ResizeObserver(updateScroll);
    observer.observe(list);

    return () => {
      list.removeEventListener("scroll", updateScroll);
      window.removeEventListener("resize", updateScroll);
      observer.disconnect();
    };
  }, [page]);

  // Pengingat ditampilkan sementara setelah produk berhasil masuk keranjang.
  useEffect(() => {
    if (!toast) return;
    const timer = window.setTimeout(() => setToast(null), 7000);
    return () => window.clearTimeout(timer);
  }, [toast]);

  const cartItems = useMemo(
    () => menuItems.filter((item) => (cart[item.id] || 0) > 0),
    [cart]
  );
  const totalQuantity = cartItems.reduce((sum, item) => sum + cart[item.id], 0);
  const subtotal = cartItems.reduce((sum, item) => sum + item.price * cart[item.id], 0);

  const visibleItems = menuItems.filter((item) => {
    const correctPage = page !== "favorit" || favorites.includes(item.id);
    const correctCategory = category === "Semua Menu" || item.category === category;
    const correctSearch = `${item.name} ${item.description}`
      .toLowerCase()
      .includes(search.toLowerCase().trim());
    return correctPage && correctCategory && correctSearch;
  });

  function changeQuantity(id: string, difference: number) {
    setCart((current) => {
      const next = { ...current };
      const quantity = Math.max(0, (next[id] || 0) + difference);
      if (quantity === 0) delete next[id];
      else next[id] = quantity;
      return next;
    });
  }

  function openDetail(item: MenuItem) {
    setDetailQuantity(1);
    setToast(null);
    setDrawer(null);
    setDetail(item);
  }

  function addDetailToCart() {
    if (!detail || !detail.available) return;
    const quantity = detailQuantity;
    const item = detail;

    setCart((current) => ({
      ...current,
      [item.id]: (current[item.id] || 0) + quantity,
    }));
    setDetail(null);
    setToast({ name: item.name, quantity });
  }

  function toggleFavorite(id: string) {
    setFavorites((current) =>
      current.includes(id) ? current.filter((value) => value !== id) : [...current, id]
    );
  }

  function openPage(nextPage: Page) {
    setPage(nextPage);
    setCategory("Semua Menu");
    setSearch("");
    setSearchOpen(false);
    setDetail(null);
    setDrawer(null);
    setToast(null);
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  function goToMenu() {
    setPage("beranda");
    document.getElementById("menu")?.scrollIntoView({ behavior: "smooth" });
  }

  function submitCheckout(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!meja || !customerName.trim() || !customerPhone.trim() || totalQuantity === 0) return;
    setDrawer("preview");
  }

  function closeModal() {
    setDrawer(null);
    setDetail(null);
  }

  return (
    <div className="tk-app">
      <header className="tk-header">
        <div className="tk-header-inner">
          <nav className="tk-header-nav" aria-label="Navigasi pelanggan">
            <button className={page === "beranda" ? "tk-nav-active" : ""} onClick={() => openPage("beranda")}>Beranda</button>
            <button className={page === "pesanan" ? "tk-nav-active" : ""} onClick={() => openPage("pesanan")}>Pesanan</button>
            <button className={page === "favorit" ? "tk-nav-active" : ""} onClick={() => openPage("favorit")}>Favorit</button>
          </nav>
          <button className="tk-logo" onClick={() => openPage("beranda")} aria-label="Kembali ke beranda">
            <span className="tk-logo-icon"><Coffee size={22} /></span>
            <span>tap<span>kasir</span><i>.</i></span>
          </button>
          <div className="tk-header-right">
            {meja && <span className="tk-table-badge"><MapPin size={14} /> Meja {meja}</span>}
            <button className="tk-icon-button" aria-label="Cari menu" onClick={() => { setPage("beranda"); setSearchOpen((v) => !v); }}><Search size={20} /></button>
            <button className="tk-icon-button tk-cart-icon" aria-label="Buka keranjang" onClick={() => setDrawer("keranjang")}>
              <ShoppingBag size={20} />
              {totalQuantity > 0 && <span className="tk-counter">{totalQuantity}</span>}
            </button>
          </div>
        </div>
      </header>

      {page === "beranda" && (
        <section className="tk-hero" aria-label="Selamat datang">
          <div className="tk-hero-copy">
            <span className="tk-eyebrow"><span /> YOUR COFFEE MOMENT</span>
            <h1>Lorem Ipsum<br /><em>Lorem Ipsum</em></h1>
            <p>lorem ipsum</p>
            <button className="tk-primary" onClick={goToMenu}>Jelajahi Menu <ArrowRight size={17} /></button>
            <div
                className="tk-hero-note"
                style={{ visibility: "hidden" }}
                >
                <MapPin size={16} />
                <span>Scan QR di meja untuk mulai memesan</span>
                </div>
          </div>
          <div className="tk-hero-photo" role="img" aria-label="Interior kafe bernuansa hangat">
          </div>
        </section>
      )}

      {page === "pesanan" ? (
        <section className="tk-empty-page">
          <span className="tk-kicker">PESANANMU</span>
          <h2>Belum ada pesanan</h2>
          <p>Riwayat dan status pesanan akan muncul di sini setelah sistem terhubung dengan backend.</p>
          <button className="tk-primary" onClick={() => openPage("beranda")}>Lihat Menu <ArrowRight size={16} /></button>
        </section>
      ) : (
        <main id="menu" className="tk-menu-section">
          <div className="tk-section-head">
            <div>
              <span className="tk-kicker">{page === "favorit" ? "PILIHANMU" : "PILIH SESUAI SELERAMU"}</span>
              <h2>{page === "favorit" ? "Menu favoritmu" : "Temukan menu favoritmu"}</h2>
              <p>{page === "favorit" ? "Semua yang kamu sukai, di satu tempat." : ""}</p>
            </div>
            <button className="tk-search-toggle" onClick={() => setSearchOpen((open) => !open)}><Search size={17} /> Cari menu</button>
          </div>

          {searchOpen && (
            <div className="tk-search-wrap">
              <Search size={18} />
              <input aria-label="Cari produk" placeholder="Cari espresso, matcha, croissant..." value={search} onChange={(event) => setSearch(event.target.value)} autoFocus />
              {search && <button onClick={() => setSearch("")} aria-label="Hapus pencarian"><X size={18} /></button>}
            </div>
          )}

          <div className={`tk-category-bar ${categoryScroll.right ? "tk-category-bar-more" : ""}`}>
            <div className="tk-category-list" role="group" aria-label="Kategori menu" ref={categoryListRef}>
              {categories.map((item) => (
                <button
                  key={item}
                  className={`tk-category ${category === item ? "tk-category-selected" : ""}`}
                  onClick={() => setCategory(item)}
                  aria-pressed={category === item}
                >
                  {item}
                </button>
              ))}
            </div>
            
          </div>

          {visibleItems.length === 0 ? (
            <div className="tk-empty-result">
              <Heart size={36} />
              <h3>{page === "favorit" ? "Belum ada menu favorit" : "Menu tidak ditemukan"}</h3>
              <p>{page === "favorit" ? "Klik ikon hati pada menu yang kamu sukai." : "Coba kata kunci atau kategori lain."}</p>
            </div>
          ) : (
            <div className="tk-product-grid">
              {visibleItems.map((item) => (
                <article key={item.id} className="tk-product-card">
                  <div className="tk-product-photo" style={{ backgroundImage: `url("${item.image}")` }} role="img" aria-label={item.name}>
                    <button type="button" className="tk-product-photo-open" aria-label={`Lihat detail ${item.name}`} onClick={() => openDetail(item)} />
                    {item.tag && <span className="tk-product-tag">{item.tag}</span>}
                    <button className={`tk-heart ${favorites.includes(item.id) ? "tk-heart-active" : ""}`} aria-label={favorites.includes(item.id) ? `Hapus ${item.name} dari favorit` : `Favoritkan ${item.name}`} onClick={() => toggleFavorite(item.id)}>
                      <Heart size={18} fill={favorites.includes(item.id) ? "currentColor" : "none"} />
                    </button>
                  </div>
                  <div className="tk-product-info">
                    <span className="tk-product-category">{item.category}</span>
                    <button className="tk-product-title" onClick={() => openDetail(item)}>{item.name}</button>
                    <div className="tk-product-bottom">
                      <strong>{rupiah(item.price)}</strong>
                      <button className="tk-add" onClick={() => openDetail(item)} disabled={!item.available}><Plus size={16} /> Tambah</button>
                    </div>
                  </div>
                </article>
              ))}
            </div>
          )}
        </main>
      )}

      <footer className="tk-footer">
        <div><strong>tapkasir.</strong><p>Pesan dengan mudah, nikmati waktumu.</p></div>
        <span>Tampilan pelanggan </span>
      </footer>

      <nav className="tk-bottom-nav" aria-label="Navigasi mobile">
        <button className={page === "beranda" ? "tk-mobile-active" : ""} onClick={() => openPage("beranda")}><Coffee size={20} /> Menu</button>
        <button className={page === "pesanan" ? "tk-mobile-active" : ""} onClick={() => openPage("pesanan")}><ShoppingBag size={20} /> Pesanan</button>
        <button className={page === "favorit" ? "tk-mobile-active" : ""} onClick={() => openPage("favorit")}><Heart size={20} /> Favorit</button>
      </nav>

      {toast && (
        <div className="tk-toast" role="status" aria-live="polite" aria-atomic="true">
          <span className="tk-toast-icon"><Check size={19} /></span>
          <div className="tk-toast-copy">
            <strong>{toast.quantity} × {toast.name} ditambahkan</strong>
            <p>Produk sudah di keranjang. Buka keranjang untuk melanjutkan pesanan.</p>
          </div>
          <button className="tk-toast-action" onClick={() => { setToast(null); setDrawer("keranjang"); }}>
            Lihat keranjang <ArrowRight size={15} />
          </button>
          <button className="tk-toast-dismiss" aria-label="Tutup pengingat" onClick={() => setToast(null)}><X size={17} /></button>
        </div>
      )}

      {(drawer || detail) && (
        <div className={`tk-overlay ${detail ? "tk-overlay-detail" : ""}`}>
          <button className="tk-overlay-backdrop" aria-label="Tutup dialog" onClick={closeModal} />
          <section className="tk-modal" role="dialog" aria-modal="true" aria-label={detail ? `Detail ${detail.name}` : "Panel pelanggan"}>
            <div className="tk-modal-head">
              <button className="tk-modal-back" aria-label="Kembali" onClick={() => {
                if (detail) setDetail(null);
                else if (drawer === "checkout") setDrawer("keranjang");
                else if (drawer === "preview") setDrawer("checkout");
                else closeModal();
              }}><ArrowLeft size={19} /></button>
              <h2>{detail ? "Detail Menu" : drawer === "keranjang" ? "Keranjang" : drawer === "checkout" ? "Checkout" : "Pratinjau Pesanan"}</h2>
              <button className="tk-modal-close" aria-label="Tutup" onClick={closeModal}><X size={19} /></button>
            </div>
            <div className="tk-modal-body">
              {detail && (
                <div className="tk-detail">
                  <div className="tk-detail-photo" style={{ backgroundImage: `url("${detail.image}")` }} role="img" aria-label={detail.name} />
                  <span className="tk-product-category">{detail.category}</span>
                  <h3>{detail.name}</h3>
                  <p>{detail.description}</p>
                  <div className="tk-detail-price">
                    <span>Harga satuan</span>
                    <strong>{rupiah(detail.price)}</strong>
                  </div>
                  <div className="tk-detail-quantity">
                    <div>
                      <strong>Jumlah pesanan</strong>
                      <small>Atur jumlah sebelum memasukkan ke keranjang</small>
                    </div>
                    <div className="tk-detail-stepper">
                      <button type="button" aria-label="Kurangi jumlah" onClick={() => setDetailQuantity((n) => Math.max(1, n - 1))} disabled={detailQuantity <= 1}><Minus size={17} /></button>
                      <output aria-live="polite">{detailQuantity}</output>
                      <button type="button" aria-label="Tambah jumlah" onClick={() => setDetailQuantity((n) => Math.min(99, n + 1))} disabled={detailQuantity >= 99}><Plus size={17} /></button>
                    </div>
                  </div>
                  <button className="tk-primary tk-full" onClick={addDetailToCart} disabled={!detail.available}>
                    <ShoppingBag size={17} /> Masukkan ke Keranjang · {rupiah(detail.price * detailQuantity)}
                  </button>
                </div>
              )}

              {!detail && drawer === "keranjang" && (
                <div>
                  {cartItems.length === 0 ? (
                    <div className="tk-cart-empty"><ShoppingBag size={42} /><h3>Keranjang masih kosong</h3><p>Pilih menu yang kamu suka terlebih dahulu.</p><button className="tk-primary" onClick={closeModal}>Jelajahi Menu</button></div>
                  ) : (
                    <>
                      <div className="tk-cart-list">
                        {cartItems.map((item) => (
                          <div className="tk-cart-row" key={item.id}>
                            <div className="tk-cart-photo" style={{ backgroundImage: `url("${item.image}")` }} />
                            <div className="tk-cart-description"><strong>{item.name}</strong><span>{rupiah(item.price)}</span><div className="tk-quantity"><button aria-label={`Kurangi ${item.name}`} onClick={() => changeQuantity(item.id, -1)}><Minus size={14} /></button><b>{cart[item.id]}</b><button aria-label={`Tambah ${item.name}`} onClick={() => changeQuantity(item.id, 1)}><Plus size={14} /></button></div></div>
                            <button className="tk-trash" aria-label={`Hapus ${item.name}`} onClick={() => changeQuantity(item.id, -cart[item.id])}><Trash2 size={17} /></button>
                          </div>
                        ))}
                      </div>
                      <div className="tk-total"><span>Subtotal ({totalQuantity} item)</span><strong>{rupiah(subtotal)}</strong></div>
                      <p className="tk-helper">Harga final dan ketersediaan stok nantinya dikonfirmasi oleh backend.</p>
                      <button className="tk-primary tk-full" onClick={() => setDrawer("checkout")}>Lanjut Checkout <ArrowRight size={16} /></button>
                    </>
                  )}
                </div>
              )}

              {!detail && drawer === "checkout" && (
                <form className="tk-checkout" onSubmit={submitCheckout}>
                  <div className="tk-checkout-table"><MapPin size={19} /><span><small>LOKASI PEMESANAN</small><strong>{meja ? `Meja ${meja}` : "Nomor meja belum terdeteksi"}</strong></span></div>
                  {!meja && <p className="tk-warning">Untuk melanjutkan, buka halaman dari QR Code meja. Saat pengembangan, coba URL <code>/?meja=07</code>.</p>}
                  <label>Nama pelanggan<input required maxLength={60} placeholder="Masukkan nama" value={customerName} onChange={(event) => setCustomerName(event.target.value)} /></label>
                  <label>Nomor telepon<input required type="tel" inputMode="tel" pattern="[0-9+ ]{8,16}" title="Masukkan 8–16 karakter berupa angka, spasi, atau tanda +" placeholder="Contoh: 081234567890" value={customerPhone} onChange={(event) => setCustomerPhone(event.target.value)} /></label>
                  <h3>Metode pembayaran</h3>
                  <div className="tk-payment-list">
                    <button type="button" className={payment === "qris" ? "tk-payment-selected" : ""} onClick={() => setPayment("qris")}><QrCode size={23} /><span><strong>QRIS</strong><small>Bayar dengan e-wallet atau mobile banking</small></span>{payment === "qris" && <Check size={18} />}</button>
                    <button type="button" className={payment === "transfer" ? "tk-payment-selected" : ""} onClick={() => setPayment("transfer")}><CreditCard size={23} /><span><strong>Transfer bank</strong><small>Instruksi rekening setelah order dibuat</small></span>{payment === "transfer" && <Check size={18} />}</button>
                    <button type="button" className={payment === "cash" ? "tk-payment-selected" : ""} onClick={() => setPayment("cash")}><Banknote size={23} /><span><strong>Tunai di kasir</strong><small>Bayar dengan menyebutkan nama dan nomor meja</small></span>{payment === "cash" && <Check size={18} />}</button>
                  </div>
                  <div className="tk-total"><span>Subtotal</span><strong>{rupiah(subtotal)}</strong></div>
                  <button className="tk-primary tk-full" type="submit" disabled={!meja || totalQuantity === 0}>Lihat Pratinjau Pesanan <ArrowRight size={16} /></button>
                </form>
              )}

              {!detail && drawer === "preview" && (
                <div className="tk-preview">
                  <span className="tk-preview-icon"><Coffee size={33} /></span>
                  <span className="tk-kicker">PRATINJAU • BUKAN TRANSAKSI AKTIF</span>
                  <h3>Pesanan siap ditinjau</h3>
                  <p>Ini hanya simulasi frontend. Belum ada pesanan yang dikirim ke kasir atau pembayaran yang diproses.</p>
                  <div className="tk-receipt">
                    <div><span>Nama</span><strong>{customerName.trim()}</strong></div>
                    <div><span>Nomor meja</span><strong>{meja}</strong></div>
                    <div><span>Nomor telepon</span><strong>{customerPhone.trim()}</strong></div>
                    <div><span>Jumlah item</span><strong>{totalQuantity}</strong></div>
                    <div><span>Pembayaran</span><strong>{payment === "qris" ? "QRIS" : payment === "transfer" ? "Transfer bank" : "Tunai di kasir"}</strong></div>
                    <div className="tk-receipt-total"><span>Total sementara</span><strong>{rupiah(subtotal)}</strong></div>
                  </div>
                  <div className="tk-warning">QRIS asli, nomor rekening, dan konfirmasi pembayaran baru tersedia setelah integrasi backend.</div>
                  <button className="tk-primary tk-full" onClick={() => setDrawer("checkout")}>Kembali ke Checkout</button>
                </div>
              )}
            </div>
          </section>
        </div>
      )}
    </div>
  );
}
