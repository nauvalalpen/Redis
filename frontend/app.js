// ============================================================
// MiniKatalog - Frontend JavaScript
// Mendemonstrasikan: Auth Session, Rate Limiting, Cache-Aside
// ============================================================

const API = 'http://localhost:8080/api';

// ============================================================
// STATE MANAGEMENT
// ============================================================
let state = {
  token: localStorage.getItem('mk_token') || null,
  username: localStorage.getItem('mk_username') || null,
  products: [],
  cacheHistory: []
};

// ============================================================
// DOM HELPERS
// ============================================================
const $ = (id) => document.getElementById(id);

function showToast(message, type = 'info', duration = 3000) {
  const toast = $('toast');
  toast.textContent = message;
  toast.className = `toast ${type}`;
  toast.classList.remove('hidden');
  setTimeout(() => toast.classList.add('hidden'), duration);
}

function showEl(id) { $(id).classList.remove('hidden'); }
function hideEl(id) { $(id).classList.add('hidden'); }

function formatPrice(price) {
  return 'Rp ' + Number(price).toLocaleString('id-ID');
}

// ============================================================
// AUTH STATE UPDATE
// ============================================================
function updateAuthUI() {
  if (state.token) {
    showEl('userInfo');
    hideEl('guestInfo');
    showEl('addProductBtn');
    $('usernameDisplay').textContent = state.username || 'User';
    $('sessionIndicator').textContent = 'Aktif';
    $('sessionIndicator').className = 'rf-status active';
  } else {
    hideEl('userInfo');
    showEl('guestInfo');
    hideEl('addProductBtn');
    hideEl('addProductForm');
    $('sessionIndicator').textContent = 'Tidak Aktif';
    $('sessionIndicator').className = 'rf-status inactive';
  }
}

// ============================================================
// API CALLS
// ============================================================
async function apiFetch(path, method = 'GET', body = null) {
  const headers = { 'Content-Type': 'application/json' };
  if (state.token) {
    headers['Authorization'] = `Bearer ${state.token}`;
  }
  const opts = { method, headers };
  if (body) opts.body = JSON.stringify(body);

  const res = await fetch(`${API}${path}`, opts);
  const data = await res.json().catch(() => ({}));
  return { ok: res.ok, status: res.status, data, headers: res.headers };
}

// ============================================================
// AUTH OPERATIONS
// ============================================================
async function register(username, password) {
  const { ok, data } = await apiFetch('/auth/register', 'POST', { username, password });
  return { ok, message: ok ? data.message : data.error };
}

async function login(username, password) {
  const { ok, status, data } = await apiFetch('/auth/login', 'POST', { username, password });
  if (ok) {
    state.token = data.token;
    state.username = data.username;
    localStorage.setItem('mk_token', data.token);
    localStorage.setItem('mk_username', data.username);
    return { ok: true };
  }
  return { ok: false, message: data.error, isRateLimit: status === 429 };
}

async function logout() {
  await apiFetch('/auth/logout', 'POST');
  state.token = null;
  state.username = null;
  localStorage.removeItem('mk_token');
  localStorage.removeItem('mk_username');
  updateAuthUI();
  showToast('Logout berhasil', 'success');
}

// ============================================================
// PRODUCT OPERATIONS
// ============================================================
async function loadProducts() {
  showEl('productLoading');
  hideEl('productGrid');
  hideEl('productEmpty');

  try {
    const start = performance.now();
    const { ok, data, headers } = await apiFetch('/products');
    const elapsed = (performance.now() - start).toFixed(1);

    if (!ok) { throw new Error(data.error || 'Gagal memuat produk'); }

    const cacheStatus = data.cache_status || headers.get('X-Cache') || 'UNKNOWN';
    const serverTime = data.response_time || headers.get('X-Response-Time') || '—';

    // Update Redis status panel
    updateCacheStatus(cacheStatus, serverTime, elapsed);

    state.products = data.data || [];
    renderProducts(state.products);
  } catch (err) {
    hideEl('productLoading');
    showToast('Gagal memuat produk: ' + err.message, 'error');
  }
}

function updateCacheStatus(cacheStatus, serverTime, clientElapsed) {
  const indicator = $('cacheIndicator');
  const rtDisplay = $('responseTimeDisplay');

  if (cacheStatus === 'CACHE_HIT') {
    indicator.textContent = '⚡ CACHE HIT';
    indicator.className = 'rf-status hit';
  } else if (cacheStatus === 'CACHE_MISS') {
    indicator.textContent = '🔄 CACHE MISS';
    indicator.className = 'rf-status miss';
  }

  rtDisplay.textContent = `Server: ${serverTime} | Client: ${clientElapsed}ms`;
  rtDisplay.className = 'rf-status';

  // Tambah ke history
  const isHit = cacheStatus === 'CACHE_HIT';
  state.cacheHistory.push({ status: cacheStatus, serverTime, clientElapsed });

  const historyEl = $('cacheHistory');
  const event = document.createElement('span');
  event.className = `cache-event ${isHit ? 'hit' : 'miss'}`;
  event.textContent = `#${state.cacheHistory.length} ${isHit ? '⚡HIT' : '🔄MISS'} (${clientElapsed}ms)`;
  event.title = `Server: ${serverTime}`;
  historyEl.appendChild(event);
}

function renderProducts(products) {
  hideEl('productLoading');

  if (!products || products.length === 0) {
    showEl('productEmpty');
    return;
  }

  showEl('productGrid');
  const grid = $('productGrid');
  grid.innerHTML = '';

  products.forEach(p => {
    const card = document.createElement('div');
    card.className = 'product-card';
    card.innerHTML = `
      <div class="product-card-header">
        <div class="product-name">${escHtml(p.name)}</div>
        ${p.category ? `<span class="product-category">${escHtml(p.category)}</span>` : ''}
      </div>
      <div class="product-price">${formatPrice(p.price)}</div>
      ${p.description ? `<div class="product-desc">${escHtml(p.description)}</div>` : ''}
      <div class="product-meta">
        <span class="product-stock">📦 Stok: ${p.stock || 0}</span>
        <span class="product-by">oleh ${escHtml(p.created_by || 'unknown')}</span>
      </div>
      ${state.token ? `
        <div class="product-actions">
          <button class="btn btn-sm btn-danger" onclick="deleteProduct('${p.id}')">🗑️ Hapus</button>
        </div>
      ` : ''}
    `;
    grid.appendChild(card);
  });
}

async function createProduct(formData) {
  const { ok, data } = await apiFetch('/products', 'POST', formData);
  return { ok, message: ok ? 'Produk berhasil ditambahkan!' : (data.error || 'Gagal menambahkan produk') };
}

async function deleteProduct(id) {
  if (!confirm('Hapus produk ini? Cache akan otomatis di-refresh.')) return;
  const { ok, data } = await apiFetch(`/products/${id}`, 'DELETE');
  if (ok) {
    showToast('Produk dihapus. Cache di-invalidate → refresh akan CACHE MISS!', 'success', 4000);
    loadProducts();
  } else {
    showToast(data.error || 'Gagal menghapus', 'error');
  }
}

// ============================================================
// ESCAPE HTML (Security)
// ============================================================
function escHtml(str) {
  if (!str) return '';
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

// ============================================================
// MODAL MANAGEMENT
// ============================================================
function showModal(id) { showEl(id); }
function hideModal(id) { hideEl(id); }

// ============================================================
// EVENT LISTENERS
// ============================================================
document.addEventListener('DOMContentLoaded', () => {
  updateAuthUI();
  loadProducts();

  // --- Navbar buttons ---
  $('showLoginBtn').addEventListener('click', () => showModal('loginModal'));
  $('showRegisterBtn').addEventListener('click', () => showModal('registerModal'));
  $('logoutBtn').addEventListener('click', logout);

  // --- Modal close ---
  $('loginOverlay').addEventListener('click', () => hideModal('loginModal'));
  $('registerOverlay').addEventListener('click', () => hideModal('registerModal'));
  $('closeLoginModal').addEventListener('click', () => hideModal('loginModal'));
  $('closeRegisterModal').addEventListener('click', () => hideModal('registerModal'));

  // --- Switch between modals ---
  $('switchToRegister').addEventListener('click', (e) => {
    e.preventDefault();
    hideModal('loginModal');
    showModal('registerModal');
  });
  $('switchToLogin').addEventListener('click', (e) => {
    e.preventDefault();
    hideModal('registerModal');
    showModal('loginModal');
  });

  // --- Refresh button ---
  $('refreshBtn').addEventListener('click', loadProducts);

  // --- Cache info ---
  $('clearCacheNote').addEventListener('click', () => {
    showToast('Cache otomatis di-invalidate saat ada produk baru/hapus. TTL cache: 5 menit.', 'info', 5000);
  });

  // --- Add product toggle ---
  $('addProductBtn').addEventListener('click', () => {
    showEl('addProductForm');
    hideEl('addProductBtn');
  });
  $('cancelProductBtn').addEventListener('click', () => {
    hideEl('addProductForm');
    showEl('addProductBtn');
    hideEl('productFormError');
  });

  // ============================================================
  // REGISTER FORM
  // ============================================================
  $('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    hideEl('registerError');
    hideEl('registerSuccess');

    const username = $('regUsername').value.trim();
    const password = $('regPassword').value;

    const { ok, message } = await register(username, password);
    if (ok) {
      showEl('registerSuccess');
      $('registerSuccess').textContent = message;
      $('registerForm').reset();
      showToast('Registrasi berhasil! Silakan login.', 'success');
      setTimeout(() => {
        hideModal('registerModal');
        showModal('loginModal');
      }, 1500);
    } else {
      showEl('registerError');
      $('registerError').textContent = message;
    }
  });

  // ============================================================
  // LOGIN FORM
  // ============================================================
  $('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    hideEl('loginError');
    hideEl('rateLimitInfo');

    const username = $('loginUsername').value.trim();
    const password = $('loginPassword').value;
    const btn = $('loginSubmitBtn');

    btn.disabled = true;
    btn.textContent = 'Memverifikasi...';

    const result = await login(username, password);

    btn.disabled = false;
    btn.textContent = 'Login';

    if (result.ok) {
      hideModal('loginModal');
      updateAuthUI();
      showToast(`Selamat datang, ${state.username}! Session tersimpan di Redis.`, 'success', 4000);
      $('loginForm').reset();
    } else {
      showEl('loginError');
      $('loginError').textContent = result.message;
      if (result.isRateLimit) {
        showEl('rateLimitInfo');
      }
    }
  });

  // ============================================================
  // ADD PRODUCT FORM
  // ============================================================
  $('submitProductBtn').addEventListener('click', async () => {
    hideEl('productFormError');

    const name = $('pName').value.trim();
    const price = parseInt($('pPrice').value) || 0;
    const category = $('pCategory').value.trim();
    const stock = parseInt($('pStock').value) || 0;
    const description = $('pDesc').value.trim();

    if (!name) {
      showEl('productFormError');
      $('productFormError').textContent = 'Nama produk wajib diisi';
      return;
    }

    const btn = $('submitProductBtn');
    btn.disabled = true;
    btn.textContent = 'Menyimpan...';

    const { ok, message } = await createProduct({ name, price, category, stock, description });

    btn.disabled = false;
    btn.textContent = 'Simpan Produk';

    if (ok) {
      showToast(message + ' Cache di-invalidate → request berikutnya CACHE MISS!', 'success', 4000);
      hideEl('addProductForm');
      showEl('addProductBtn');
      $('pName').value = '';
      $('pPrice').value = '';
      $('pCategory').value = '';
      $('pStock').value = '';
      $('pDesc').value = '';
      loadProducts();
    } else {
      showEl('productFormError');
      $('productFormError').textContent = message;
    }
  });
});
