import React, { useEffect, useMemo, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter, Navigate, NavLink, Outlet, Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import { 
  CalendarDays, Check, ChevronLeft, ChevronRight, CircleUserRound, Clock3, 
  FileText, LogOut, MapPin, Menu, RefreshCw, ShieldCheck, Ticket, 
  UsersRound, UserRound, X, XCircle 
} from 'lucide-react';
import './styles.css';

// ==========================================
// CONFIGS & UTILS
// ==========================================
const API = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const today = () => { 
  const d = new Date(); 
  const off = d.getTimezoneOffset(); 
  return new Date(d.getTime() - off * 60000).toISOString().slice(0, 10); 
};

const dateOnly = d => new Date(d.getFullYear(), d.getMonth(), d.getDate());

const addDays = (value, amount) => { 
  const d = new Date(`${value}T00:00:00`); 
  d.setDate(d.getDate() + amount); 
  return d.toISOString().slice(0, 10); 
};

const thaiDate = value => value 
  ? new Intl.DateTimeFormat('th-TH', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' }).format(new Date(`${String(value).slice(0, 10)}T00:00:00`)) 
  : '-';

const thaiShort = value => value 
  ? new Intl.DateTimeFormat('th-TH', { day: 'numeric', month: 'short', year: 'numeric' }).format(new Date(`${String(value).slice(0, 10)}T00:00:00`)) 
  : '-';

const unwrap = json => json?.data ?? json;

async function request(path, options = {}) { 
  const token = localStorage.getItem('special_booking_token'); 
  const headers = { 
    'Content-Type': 'application/json', 
    ...(options.headers || {}) 
  }; 
  
  if (token) headers.Authorization = `Bearer ${token}`; 
  
  const res = await fetch(`${API}${path}`, { ...options, headers }); 
  const json = await res.json().catch(() => ({})); 
  
  if (!res.ok) { 
    const error = new Error(json.error || json.message || `Request failed (${res.status})`); 
    error.payload = json; 
    throw error; 
  } 
  return unwrap(json); 
}

// ==========================================
// SHARED / SMALL COMPONENTS
// ==========================================
function Logo() { 
  return (
    <div className="brand">
      <div className="brand-mark"><ShieldCheck size={19}/></div>
      <div>
        <strong>Fitness Management System</strong>
        <span>ระบบบริหารฟิตเนส</span>
        <small>By T18</small>
      </div>
    </div>
  ); 
}

function Icon({ name, size = 17 }) { 
  const map = { 
    home: FileText, 
    booking: CalendarDays, 
    users: UsersRound, 
    member: CircleUserRound, 
    area: MapPin, 
    clock: Clock3, 
    ticket: Ticket, 
    trainer: UserRound 
  }; 
  const C = map[name] || FileText; 
  return <C size={size}/>; 
}

const nav = [
  { path: '/', label: 'หน้าแรก', icon: 'home' }, 
  { path: '/reports', label: 'ระบบรายงานและสรุปผล', icon: 'home' }, 
  { path: '/staff', label: 'ระบบจัดการพนักงาน', icon: 'users' }, 
  { path: '/members', label: 'ระบบบริการสมาชิก', icon: 'member' }, 
  { path: '/classes', label: 'ระบบจัดการคลาสกลุ่ม', icon: 'booking' }, 
  { path: '/special-areas', label: 'ระบบจัดการพื้นที่พิเศษ', icon: 'area' }, 
  { path: '/trainers', label: 'ระบบจัดการเทรนเนอร์', icon: 'trainer' }, 
  { path: '/equipment', label: 'ระบบเครื่องออกกำลังกาย', icon: 'area' }, 
  { path: '/maintenance', label: 'ระบบซ่อมบำรุง', icon: 'clock' }, 
  { path: '/lost-and-found', label: 'ระบบจัดการของสูญหาย', icon: 'ticket' }, 
  { path: '/member-data', label: 'ระบบจัดการข้อมูลสมาชิก', icon: 'member' }
];

function Stat({ icon, label, value, hint }) {
  return (
    <div className="stat">
      <div className="stat-icon"><Icon name={icon} size={16}/></div>
      <div>
        <span>{label}</span>
        <strong>{value}</strong>
        {hint && <small>{hint}</small>}
      </div>
    </div>
  );
}

function PageError({ error, onRetry }) {
  return error ? (
    <div className="page-error">
      <XCircle size={17}/>
      <span>{error}</span>
      <button onClick={onRetry}>ลองใหม่</button>
    </div>
  ) : null;
}

// ==========================================
// AUTH & LAYOUT
// ==========================================
function Login({ onAuth }) { 
  const [mode, setMode] = useState('login'); 
  const [form, setForm] = useState({ full_name: '', email: '', phone: '', password: '' }); 
  const [busy, setBusy] = useState(false); 
  const [error, setError] = useState(''); 

  const submit = async e => {
    e.preventDefault();
    setBusy(true);
    setError('');
    try {
      const data = await request(`/auth/${mode === 'login' ? 'login' : 'register'}`, {
        method: 'POST',
        body: JSON.stringify(form)
      });
      localStorage.setItem('special_booking_token', data.token);
      onAuth(data.user);
    } catch (err) {
      setError(err.message);
    } finally {
      setBusy(false);
    }
  }; 

  const set = (k, v) => setForm({ ...form, [k]: v }); 

  return (
    <main className="auth-shell">
      <div className="auth-art">
        <Logo/>
        <div className="art-orbit"></div>
        <div className="art-copy">
          <p className="eyebrow">FITNESS OPERATIONS</p>
          <h1>พื้นที่ที่ดี<br/><em>เริ่มจากการจัดการที่ชัดเจน</em></h1>
          <p>จัดการพื้นที่ นัดหมายเทรนเนอร์ และดูตารางทั้งหมดจากระบบเดียว</p>
        </div>
        <div className="art-stamp">3A82DE<br/><span>FITNESS / 2026</span></div>
      </div>
      
      <section className="auth-panel">
        <div className="auth-panel-inner">
          <p className="eyebrow blue">WELCOME BACK</p>
          <h2>{mode === 'login' ? 'เข้าสู่ระบบ' : 'สร้างบัญชีสมาชิก'}</h2>
          <p className="muted">{mode === 'login' ? 'จัดการงานฟิตเนสของคุณ' : 'กรอกข้อมูลเพื่อเริ่มต้นใช้งาน'}</p>
          
          <form onSubmit={submit}>
            {mode === 'register' && (
              <label>ชื่อ-นามสกุล
                <input required value={form.full_name} onChange={e => set('full_name', e.target.value)} placeholder="เช่น ริว พีรพัฒน์"/>
              </label>
            )}
            <label>อีเมล
              <input required type="email" value={form.email} onChange={e => set('email', e.target.value)} placeholder="you@example.com"/>
            </label>
            {mode === 'register' && (
              <label>เบอร์โทรศัพท์
                <input required value={form.phone} onChange={e => set('phone', e.target.value)} placeholder="0812345678"/>
              </label>
            )}
            <label>รหัสผ่าน
              <input required minLength="8" type="password" value={form.password} onChange={e => set('password', e.target.value)} placeholder="อย่างน้อย 8 ตัวอักษร"/>
            </label>
            
            {error && <div className="form-error"><XCircle size={16}/>{error}</div>}
            
            <button className="primary wide" disabled={busy}>
              {busy ? 'กำลังตรวจสอบ...' : mode === 'login' ? 'เข้าสู่ระบบ' : 'สมัครสมาชิก'}
              <span>→</span>
            </button>
          </form>
          
          <button className="text-button" onClick={() => { setMode(mode === 'login' ? 'register' : 'login'); setError(''); }}>
            {mode === 'login' ? 'ยังไม่มีบัญชี? สมัครสมาชิก' : 'มีบัญชีอยู่แล้ว? เข้าสู่ระบบ'}
          </button>
        </div>
      </section>
    </main>
  ); 
}

function Layout({ user, onLogout }) {
  const loc = useLocation();
  const active = nav.find(n => n.path === loc.pathname) || 
    nav.find(n => loc.pathname.startsWith('/special-areas') ? n.path === '/special-areas' : loc.pathname.startsWith('/trainers') ? n.path === '/trainers' : false) || 
    nav[0];
  const [mobile, setMobile] = useState(false);

  return (
    <div className="app-shell">
      <aside className={`sidebar ${mobile ? 'open' : ''}`}>
        <Logo/>
        <div className="nav-title">เมนูหลัก</div>
        <nav>
          {nav.map(n => (
            <NavLink 
              key={n.path} 
              end={n.path === '/'} 
              to={n.path} 
              className={({ isActive }) => isActive ? 'active' : ''} 
              onClick={() => setMobile(false)}
            >
              <Icon name={n.icon}/>
              <span>{n.label}</span>
            </NavLink>
          ))}
        </nav>
        <div className="profile">
          <div className="avatar">{(user?.full_name || 'A').slice(0, 1)}</div>
          <div>
            <span>ผู้ใช้งาน</span>
            <strong>{user?.full_name || user?.email}</strong>
          </div>
          <button onClick={onLogout} title="ออกจากระบบ"><LogOut size={15}/></button>
        </div>
      </aside>

      <main className="main">
        <header className="topbar">
          <button className="mobile-menu" onClick={() => setMobile(!mobile)}><Menu/></button>
          <div>
            <p className="breadcrumb">Fitness Management System / <strong>{active?.label}</strong></p>
            <h1>{active?.label}</h1>
            <p className="subhead">จัดการตาราง นัดหมาย และข้อมูลการใช้งานให้เป็นระบบ</p>
          </div>
          <div className="top-actions">
            <button className="refresh" onClick={() => window.location.reload()}><RefreshCw size={17}/></button>
            <div className="user-chip">
              <span className="status-dot"></span>
              {user?.role === 'admin' || user?.role === 'staff' ? 'ผู้ดูแลระบบ' : 'สมาชิก'}
            </div>
          </div>
        </header>
        <Outlet/>
      </main>
    </div>
  );
}

// ==========================================
// PAGES
// ==========================================
function Home() {
  const navigate = useNavigate();
  return (
    <div className="home-page">
      <section className="home-hero">
        <div>
          <p className="eyebrow blue">FITNESS OPERATIONS / HOME</p>
          <h2>จัดการทุกพื้นที่<br/><em>จากจุดเดียว</em></h2>
          <p>หน้าหลักของระบบรวมทางลัดไปยังพื้นที่พิเศษและตารางนัดหมายเทรนเนอร์</p>
          <button className="primary" onClick={() => navigate('/special-areas')}>
            เปิดระบบพื้นที่พิเศษ <span>→</span>
          </button>
        </div>
        <div className="home-orbit">
          <span>FITNESS<br/><b>2026</b></span>
        </div>
      </section>

      <section className="home-modules">
        <div className="module-intro">
          <p className="eyebrow">WORKSPACE MAP</p>
          <h3>โมดูลหลัก</h3>
          <p>เลือกพื้นที่ทำงานที่ต้องการ ระบบจะแยก URL ให้ชัดเจน</p>
        </div>
        <button className="module-tile featured" onClick={() => navigate('/special-areas')}>
          <span className="module-number">01</span>
          <MapPin size={20}/>
          <strong>ระบบจัดการพื้นที่พิเศษ</strong>
          <small>/special-areas</small>
          <span className="module-arrow">↗</span>
        </button>
        <button className="module-tile" onClick={() => navigate('/trainers')}>
          <span className="module-number">02</span>
          <UserRound size={20}/>
          <strong>ระบบจัดการเทรนเนอร์</strong>
          <small>/trainers</small>
          <span className="module-arrow">↗</span>
        </button>
        <div className="home-note">
          <Check size={17}/>
          <span><strong>ระบบพร้อมใช้งาน</strong><br/>ดูตารางนัดหมายรวมได้จากหน้าเทรนเนอร์</span>
        </div>
      </section>
    </div>
  );
}

function BookingPage() {
  const navigate = useNavigate();
  const [date, setDate] = useState(today());
  const [areas, setAreas] = useState([]);
  const [selected, setSelected] = useState(null);
  const [availability, setAvailability] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [confirm, setConfirm] = useState(null);
  const [success, setSuccess] = useState(null);

  const load = async () => {
    setLoading(true);
    setError('');
    try {
      const [a, av] = await Promise.all([
        request('/areas'),
        request(`/availability?date=${date}`)
      ]);
      setAreas(a || []);
      setAvailability(av?.items || []);
      if (!selected && a?.[0]) setSelected(a[0].id);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, [date]);

  const current = availability.find(x => x.area?.id === selected) || availability[0];

  useEffect(() => {
    if (current?.area?.id && current.area.id !== selected) {
      setSelected(current.area.id);
    }
  }, [current?.area?.id]);

  const slots = current?.slots || [];

  const confirmBooking = async () => {
    if (!confirm || !current?.area) return;
    try {
      const b = await request('/bookings', {
        method: 'POST',
        body: JSON.stringify({
          area_id: current.area.id,
          time_slot_id: confirm.time_slot.id,
          booking_date: date,
          note: ''
        })
      });
      setConfirm(null);
      setSuccess(b);
      load();
    } catch (e) {
      setConfirm({ ...confirm, error: e.message });
    }
  };

  return (
    <div className="content">
      <div className="toolbar">
        <label className="date-input">
          <CalendarDays size={17}/>
          <span>วันที่</span>
          <input type="date" value={date} onChange={e => setDate(e.target.value)}/>
        </label>
        <button className="outline-action" onClick={() => navigate('/special-areas/history')}>
          <Ticket size={15}/>ดูประวัติการจอง
        </button>
      </div>

      <PageError error={error} onRetry={load}/>

      <div className="stats">
        <Stat icon="clock" label="วันที่" value={thaiDate(date)} hint="วันที่เลือก"/>
        <Stat icon="area" label="พื้นที่ที่เลือก" value={current?.area?.name || 'กำลังโหลด'} hint={current?.area ? `รองรับ ${current.area.capacity} คน` : ''}/>
        <Stat icon="ticket" label="สล็อตที่ว่าง" value={`${slots.filter(x => x.available).length} / ${slots.length || 0}`} hint="พร้อมให้จองวันนี้"/>
      </div>

      <section className="section-block area-block">
        <div className="section-head">
          <div>
            <span className="section-index">01</span>
            <div>
              <p className="eyebrow">SPACE DIRECTORY</p>
              <h2>เลือกประเภทพื้นที่</h2>
            </div>
          </div>
          <span className="count-label">{areas.length} พื้นที่</span>
        </div>
        <div className="area-tabs">
          {loading 
            ? [1, 2, 3, 4].map(i => <div className="skeleton area-skeleton" key={i}/>)
            : areas.map(a => (
                <button 
                  key={a.id} 
                  className={`area-tab ${a.id === selected ? 'selected' : ''}`} 
                  onClick={() => setSelected(a.id)}
                >
                  <span>{a.name}</span>
                  <small>รองรับ {a.capacity} คน</small>
                </button>
              ))
          }
        </div>
      </section>

      <section className="section-block slots-block">
        <div className="section-head">
          <div>
            <span className="section-index">02</span>
            <div>
              <p className="eyebrow">LIVE AVAILABILITY / {date}</p>
              <h2>ตารางสล็อตเวลา <span>• {current?.area?.name || ''}</span></h2>
            </div>
          </div>
          <div className="legend">
            <span><i className="legend-dot available-dot"></i>ว่าง</span>
            <span><i className="legend-dot booked-dot"></i>จองแล้ว</span>
          </div>
        </div>
        <div className="slots-grid">
          {loading 
            ? Array.from({ length: 12 }).map((_, i) => <div className="skeleton slot-skeleton" key={i}/>)
            : slots.length 
              ? slots.map(item => (
                  <button 
                    disabled={!item.available} 
                    key={item.time_slot.id} 
                    className={`slot ${item.available ? 'available' : 'booked'}`} 
                    onClick={() => item.available && setConfirm(item)}
                  >
                    <div>
                      <strong>{item.time_slot.start_time} - {item.time_slot.end_time}</strong>
                      <span>{item.available ? 'ว่าง · คลิกเพื่อจอง' : 'จองแล้ว'}</span>
                    </div>
                    <i></i>
                  </button>
                ))
              : <div className="empty-slot">ไม่พบช่วงเวลาสำหรับวันนี้</div>
          }
        </div>
      </section>

      <BookingConfirm data={confirm} area={current?.area} date={date} onClose={() => setConfirm(null)} onConfirm={confirmBooking}/>
      <BookingSuccess booking={success} onClose={() => setSuccess(null)} onHistory={() => { setSuccess(null); navigate('/special-areas/history'); }}/>
    </div>
  );
}

function BookingConfirm({ data, area, date, onClose, onConfirm }) {
  if (!data) return null;
  return (
    <div className="dialog-backdrop">
      <div className="dialog confirm-dialog">
        <button className="dialog-close" onClick={onClose}><X size={17}/></button>
        <div className="dialog-symbol confirm-symbol"><CalendarDays size={26}/></div>
        <p className="dialog-kicker blue">ยืนยันการจอง</p>
        <h3>ตรวจสอบรายละเอียด</h3>
        <p className="dialog-sub">กดยืนยันเพื่อบันทึกการจองลงฐานข้อมูล</p>
        <div className="receipt">
          <span>พื้นที่<strong>{area?.name}</strong></span>
          <span>วันที่<strong>{thaiDate(date)}</strong></span>
          <span>เวลา<strong>{data.time_slot.start_time} - {data.time_slot.end_time}</strong></span>
        </div>
        {data.error && <div className="dialog-error"><XCircle size={15}/>{data.error}</div>}
        <div className="dialog-actions">
          <button className="secondary" onClick={onClose}>กลับไปเลือก</button>
          <button className="primary" onClick={onConfirm}>ยืนยันการจอง</button>
        </div>
      </div>
    </div>
  );
}

function BookingSuccess({ booking, onClose, onHistory }) {
  if (!booking) return null;
  return (
    <div className="dialog-backdrop">
      <div className="dialog success">
        <button className="dialog-close" onClick={onClose}><X size={17}/></button>
        <div className="dialog-symbol"><Check size={31}/></div>
        <p className="dialog-kicker">จองสำเร็จ</p>
        <h3>บันทึกการจองแล้ว</h3>
        <p className="dialog-sub">ระบบบันทึกการจองพื้นที่พิเศษเรียบร้อยแล้ว</p>
        <div className="receipt">
          <span>พื้นที่<strong>{booking.area?.name}</strong></span>
          <span>วันที่<strong>{thaiDate(booking.booking_date)}</strong></span>
          <span>เวลา<strong>{booking.time_slot?.start_time} - {booking.time_slot?.end_time}</strong></span>
          <span>รหัสการจอง<strong>#{booking.booking_code}</strong></span>
        </div>
        <div className="dialog-actions">
          <button className="secondary" onClick={onClose}>ปิด</button>
          <button className="primary" onClick={onHistory}>ดูประวัติการจอง</button>
        </div>
      </div>
    </div>
  );
}

function HistoryPage() {
  const navigate = useNavigate();
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [detail, setDetail] = useState(null);

  const load = async () => {
    setLoading(true);
    try {
      setRows(await request('/bookings') || []);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const cancel = async row => {
    if (!confirm(`ยืนยันยกเลิกการจอง ${row.booking_code}?`)) return;
    try {
      await request(`/bookings/${row.id}/cancel`, { method: 'PATCH' });
      load();
    } catch (e) {
      setError(e.message);
    }
  };

  return (
    <div className="history-page">
      <div className="history-heading">
        <div>
          <button className="back-link" onClick={() => navigate('/special-areas')}>
            <ChevronLeft size={16}/>กลับไปหน้าจอง
          </button>
          <p className="eyebrow">SPECIAL AREAS / HISTORY</p>
          <h2>ประวัติการจอง</h2>
          <p className="muted">รายการจองทั้งหมดของคุณ</p>
        </div>
        <div className="history-heading-actions">
          <button className="outline-action" onClick={() => navigate('/special-areas')}>
            <CalendarDays size={15}/>จองพื้นที่
          </button>
          <button className="refresh" onClick={load}><RefreshCw size={17}/></button>
        </div>
      </div>

      <PageError error={error} onRetry={load}/>

      {loading ? (
        <div className="history-list">
          {[1, 2, 3].map(i => <div className="history-row skeleton" key={i}/>)}
        </div>
      ) : rows.length ? (
        <div className="history-list">
          {rows.map(row => (
            <div className="history-row" key={row.id}>
              <button className="history-hit" onClick={() => setDetail(row)}>
                <div className="history-date">
                  <strong>{row.booking_date?.slice(8, 10)}</strong>
                  <span>{thaiShort(row.booking_date).split(' ')[1]}</span>
                </div>
                <div className="history-main">
                  <strong>{row.area?.name}</strong>
                  <span>{row.time_slot?.start_time} - {row.time_slot?.end_time} · #{row.booking_code}</span>
                  {row.cooldown_until && (
                    <small className="cooldown-note">
                      cooldown ถึง {new Date(row.cooldown_until).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' })}
                    </small>
                  )}
                </div>
                <em className={`tag ${row.status === 'confirmed' ? 'tag-green' : row.status === 'cancelled' ? 'tag-red' : 'tag-blue'}`}>
                  {row.status === 'confirmed' ? 'ยืนยันแล้ว' : row.status === 'cancelled' ? 'ยกเลิกแล้ว' : row.status}
                </em>
                <ChevronRight size={18}/>
              </button>
              {(row.status === 'confirmed' || row.status === 'pending') && (
                <button className="cancel-link" onClick={() => cancel(row)}>ยกเลิก</button>
              )}
            </div>
          ))}
        </div>
      ) : (
        <div className="empty-history">
          <Ticket size={30}/>
          <h3>ยังไม่มีประวัติการจอง</h3>
          <button className="primary" onClick={() => navigate('/special-areas')}>ไปหน้าจอง</button>
        </div>
      )}

      <DetailDialog booking={detail} onClose={() => setDetail(null)}/>
    </div>
  );
}

function DetailDialog({ booking, onClose }) {
  if (!booking) return null;
  return (
    <div className="dialog-backdrop">
      <div className="dialog detail-dialog">
        <button className="dialog-close" onClick={onClose}><X size={17}/></button>
        <div className="dialog-symbol detail-symbol"><Ticket size={25}/></div>
        <p className="dialog-kicker blue">BOOKING DETAIL</p>
        <h3>รายละเอียดการจอง</h3>
        <div className="detail-code">#{booking.booking_code}</div>
        <div className="receipt detail-receipt">
          <span>พื้นที่<strong>{booking.area?.name}</strong></span>
          <span>วันที่<strong>{thaiDate(booking.booking_date)}</strong></span>
          <span>เวลา<strong>{booking.time_slot?.start_time} - {booking.time_slot?.end_time}</strong></span>
          <span>สถานะ<strong>{booking.status}</strong></span>
          <span>ผู้จอง<strong>{booking.user?.full_name || 'ฉัน'}</strong></span>
          <span>หมายเหตุ<strong>{booking.note || '-'}</strong></span>
        </div>
        <button className="primary wide" onClick={onClose}>ปิดรายละเอียด</button>
      </div>
    </div>
  );
}

// ==========================================
// TRAINERS & APPOINTMENTS
// ==========================================
function TrainerCard({ trainer, onSchedule, disabled }) {
  return (
    <article className={`trainer-card ${disabled ? 'needs-member' : ''}`}>
      <div className="trainer-avatar">{trainer.full_name?.slice(0, 2)}</div>
      <div className="trainer-card-main">
        <div className="trainer-card-top">
          <div>
            <strong>{trainer.full_name}</strong>
            <span>{trainer.code} · {trainer.specialty || 'เทรนเนอร์'}</span>
          </div>
          <em className="tag tag-green">พร้อมรับนัดหมาย</em>
        </div>
        <div className="trainer-meta">
          <span><UserRound size={14}/> {trainer.phone || 'ไม่ระบุเบอร์'}</span>
          <span><CalendarDays size={14}/> ตารางสัปดาห์นี้</span>
        </div>
        <button className="primary compact" disabled={disabled} onClick={() => onSchedule(trainer)}>
          เลือกเวลานัดหมาย <ChevronRight size={15}/>
        </button>
      </div>
    </article>
  );
}

function MemberResult({ member, selected, onSelect }) { 
  return (
    <button type="button" className={`member-result ${selected ? 'selected' : ''}`} onClick={() => onSelect(member)}>
      <div className="member-avatar">{member.full_name?.slice(0, 1)}</div>
      <div>
        <strong>{member.full_name}</strong>
        <span>{member.email} · {member.phone}</span>
      </div>
      <Check size={16} className="member-check"/>
    </button>
  ); 
}

function MemberTarget({ member, onClear }) { 
  return (
    <div className="member-target">
      <div className="member-avatar target">{member.full_name?.slice(0, 1)}</div>
      <div>
        <span>สมาชิกที่เลือก</span>
        <strong>{member.full_name} · {member.email}</strong>
        <small>{member.phone || 'ไม่ระบุเบอร์'}</small>
      </div>
      <button type="button" onClick={onClear} title="เปลี่ยนสมาชิก"><X size={16}/></button>
    </div>
  ); 
}

function AppointmentModal({ trainer, member, onClose, onSaved }) { 
  const [form, setForm] = useState({ date: today(), start: '18:00', end: '19:00', note: '' }); 
  const [busy, setBusy] = useState(false); 
  const [error, setError] = useState(''); 

  if (!trainer || !member) return null; 

  const submit = async e => { 
    e.preventDefault(); 
    setBusy(true); 
    setError(''); 
    try { 
      await request('/admin/trainer-appointments', { 
        method: 'POST', 
        body: JSON.stringify({ 
          trainer_id: trainer.id, 
          member_id: member.id, 
          appointment_date: form.date, 
          start_time: form.start, 
          end_time: form.end, 
          note: form.note 
        }) 
      }); 
      onSaved(); 
    } catch (err) { 
      setError(err.message); 
    } finally { 
      setBusy(false); 
    } 
  }; 

  return (
    <div className="dialog-backdrop">
      <div className="dialog appointment-dialog">
        <button className="dialog-close" onClick={onClose}><X size={17}/></button>
        <div className="dialog-symbol confirm-symbol"><UserRound size={26}/></div>
        <p className="dialog-kicker blue">จับคู่สมาชิกกับเทรนเนอร์</p>
        <h3>ยืนยันนัดหมาย</h3>
        <p className="dialog-sub">ระบบจะบันทึกความสัมพันธ์ระหว่างสมาชิก เทรนเนอร์ และช่วงเวลาลงฐานข้อมูล</p>
        
        <div className="pair-preview">
          <span><small>สมาชิก</small><strong>{member.full_name}</strong></span>
          <span><small>เทรนเนอร์</small><strong>{trainer.full_name}</strong></span>
        </div>

        <form className="appointment-form" onSubmit={submit}>
          <label>วันที่
            <input type="date" required value={form.date} onChange={e => setForm({ ...form, date: e.target.value })}/>
          </label>
          <div className="time-fields">
            <label>เริ่ม
              <input type="time" required value={form.start} onChange={e => setForm({ ...form, start: e.target.value })}/>
            </label>
            <label>สิ้นสุด
              <input type="time" required value={form.end} onChange={e => setForm({ ...form, end: e.target.value })}/>
            </label>
          </div>
          <label>หมายเหตุ
            <textarea rows="3" value={form.note} onChange={e => setForm({ ...form, note: e.target.value })} placeholder="เป้าหมายหรือรายละเอียดเพิ่มเติม"/>
          </label>
          
          {error && <div className="dialog-error"><XCircle size={15}/>{error}</div>}
          
          <div className="dialog-actions">
            <button type="button" className="secondary" onClick={onClose}>ยกเลิก</button>
            <button className="primary" disabled={busy}>{busy ? 'กำลังบันทึก...' : 'ยืนยันการจับคู่'}</button>
          </div>
        </form>
      </div>
    </div>
  ); 
}

function TrainersPage() { 
  const navigate = useNavigate(); 
  const [trainers, setTrainers] = useState([]); 
  const [members, setMembers] = useState([]); 
  const [memberSearch, setMemberSearch] = useState(''); 
  const [selectedMember, setSelectedMember] = useState(null); 
  const [search, setSearch] = useState(''); 
  const [loading, setLoading] = useState(true); 
  const [memberLoading, setMemberLoading] = useState(false); 
  const [error, setError] = useState(''); 
  const [selectedTrainer, setSelectedTrainer] = useState(null); 

  const loadTrainers = async () => { 
    setLoading(true); 
    setError(''); 
    try { 
      setTrainers(await request(`/trainers?search=${encodeURIComponent(search)}`) || []); 
    } catch(e) { 
      setError(e.message); 
    } finally { 
      setLoading(false); 
    } 
  }; 

  const loadMembers = async (value = memberSearch) => { 
    setMemberLoading(true); 
    setError(''); 
    try { 
      setMembers(await request(`/admin/members?search=${encodeURIComponent(value)}`) || []); 
    } catch(e) { 
      setError(e.message); 
    } finally { 
      setMemberLoading(false); 
    } 
  }; 

  useEffect(() => { 
    loadTrainers(); 
    loadMembers(''); 
  }, []); 

  const doTrainerSearch = e => { e.preventDefault(); loadTrainers(); }; 
  const doMemberSearch = e => { e.preventDefault(); loadMembers(); }; 

  return (
    <div className="trainer-page">
      <div className="trainer-page-head">
        <div>
          <p className="eyebrow blue">FITNESS OPERATIONS / TRAINERS</p>
          <h2>ระบบจัดการเทรนเนอร์</h2>
          <p>ผู้ดูแลระบบจับคู่สมาชิกกับเทรนเนอร์ และจัดการตารางนัดหมายจากจุดเดียว</p>
        </div>
        <button className="primary" onClick={() => navigate('/trainers/schedule')}>
          <CalendarDays size={16}/>ดูตารางนัดหมายรวม
        </button>
      </div>

      <PageError error={error} onRetry={() => { loadTrainers(); loadMembers(); }}/>

      <section className="member-picker">
        <div className="picker-title">
          <div>
            <p className="eyebrow">01 / MEMBER MATCHING</p>
            <h3>เลือกสมาชิกเพื่อจับคู่</h3>
            <span>ค้นหาสมาชิกก่อนเลือกเทรนเนอร์</span>
          </div>
          {selectedMember && <span className="selected-label"><Check size={14}/>เลือกแล้ว</span>}
        </div>

        {selectedMember ? (
          <MemberTarget member={selectedMember} onClear={() => setSelectedMember(null)}/>
        ) : (
          <>
            <form className="member-search" onSubmit={doMemberSearch}>
              <input value={memberSearch} onChange={e => setMemberSearch(e.target.value)} placeholder="ชื่อ, รหัสสมาชิก, อีเมล หรือเบอร์โทร"/>
              <button className="primary" disabled={memberLoading}>{memberLoading ? 'กำลังค้นหา...' : 'ค้นหาสมาชิก'}</button>
            </form>
            {members.length > 0 && (
              <div className="member-results">
                {members.slice(0, 5).map(m => (
                  <MemberResult key={m.id} member={m} selected={false} onSelect={setSelectedMember}/>
                ))}
              </div>
            )}
            {!memberLoading && members.length === 0 && <p className="picker-empty">ยังไม่พบสมาชิกที่ค้นหา</p>}
          </>
        )}
      </section>

      <section className="trainer-pairing-step">
        <div className="section-head trainer-section-head">
          <div>
            <span className="section-index">02</span>
            <div>
              <p className="eyebrow">TRAINER DIRECTORY</p>
              <h3>เลือกเทรนเนอร์เพื่อจับคู่</h3>
            </div>
          </div>
          <span className="legend"><i className="legend-dot available-dot"></i>พร้อมรับนัดหมาย</span>
        </div>

        <form className="trainer-search" onSubmit={doTrainerSearch}>
          <label>ค้นหาเทรนเนอร์</label>
          <div>
            <input value={search} onChange={e => setSearch(e.target.value)} placeholder="ชื่อเทรนเนอร์, รหัส หรือความเชี่ยวชาญ"/>
            <button className="primary">ค้นหา</button>
          </div>
        </form>

        <div className="trainer-grid">
          {loading 
            ? [1, 2, 3].map(i => <div className="trainer-card skeleton" key={i}/>)
            : trainers.map(t => (
                <TrainerCard key={t.id} trainer={t} disabled={!selectedMember} onSchedule={setSelectedTrainer}/>
              ))
          }
        </div>

        {!selectedMember && (
          <div className="pair-hint">
            <UsersRound size={16}/> เลือกสมาชิกด้านบนก่อน จึงจะเปิดปุ่มเลือกเวลานัดหมายได้
          </div>
        )}
      </section>

      {selectedTrainer && (
        <AppointmentModal 
          trainer={selectedTrainer} 
          member={selectedMember} 
          onClose={() => setSelectedTrainer(null)} 
          onSaved={() => { setSelectedTrainer(null); navigate('/trainers/schedule'); }}
        />
      )}
    </div>
  ); 
}

function SchedulePage() {
  const navigate = useNavigate();
  const [start, setStart] = useState(addDays(today(), -new Date().getDay() + 1));
  const [rows, setRows] = useState([]);
  const [trainers, setTrainers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // ฟังก์ชันคำนวณวันในสัปดาห์
  const weekDays = (startDate) => Array.from({ length: 7 }, (_, i) => addDays(startDate, i));
  const days = weekDays(start);

  const load = async () => {
    setLoading(true);
    try {
      const [a, b] = await Promise.all([
        request('/trainers'),
        request(`/trainer-appointments?from=${days[0]}&to=${days[6]}`)
      ]);
      setTrainers(a || []);
      setRows(b || []);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, [start]);

  const by = (trainerId, date) => rows.filter(r => r.trainer_id === trainerId && String(r.appointment_date).slice(0, 10) === date);

  return (
    <div className="schedule-page">
      <div className="schedule-heading">
        <div>
          <button className="back-link" onClick={() => navigate('/trainers')}>
            <ChevronLeft size={16}/>กลับไดเรกทอรีเทรนเนอร์
          </button>
          <p className="eyebrow">TRAINERS / SCHEDULE</p>
          <h2>ตารางนัดหมายรวม</h2>
          <p className="muted">ภาพรวมตารางของเทรนเนอร์ทั้งหมดในสัปดาห์นี้</p>
        </div>
        <button className="primary" onClick={() => navigate('/trainers')}>
          <UserRound size={15}/>จัดการเทรนเนอร์
        </button>
      </div>

      <PageError error={error} onRetry={load}/>

      <div className="schedule-toolbar">
        <button className="outline-action" onClick={() => setStart(addDays(start, -7))}>
          <ChevronLeft size={15}/>สัปดาห์ก่อน
        </button>
        <strong>{thaiShort(days[0])} - {thaiShort(days[6])}</strong>
        <button className="outline-action" onClick={() => setStart(addDays(start, 7))}>
          สัปดาห์ถัดไป<ChevronRight size={15}/>
        </button>
      </div>

      <div className="schedule-stats">
        <Stat icon="ticket" label="สล็อตทั้งหมด" value={rows.length + 21} hint="สัปดาห์นี้"/>
        <Stat icon="clock" label="สล็อตว่าง" value={Math.max(0, 21 - rows.length)} hint="พร้อมจอง"/>
        <Stat icon="ticket" label="สล็อตจองแล้ว" value={rows.length} hint="อัปเดตล่าสุด"/>
        <Stat icon="trainer" label="เทรนเนอร์" value={trainers.length} hint="ที่เปิดใช้งาน"/>
      </div>

      <section className="schedule-table-wrap">
        <div className="schedule-table-head">
          <span>เทรนเนอร์</span>
          {days.map(d => (
            <span key={d}>
              {new Intl.DateTimeFormat('th-TH', { weekday: 'short' }).format(new Date(`${d}T00:00:00`))}
              <small>{d.slice(8, 10)} / {d.slice(5, 7)}</small>
            </span>
          ))}
        </div>

        {loading 
          ? [1, 2, 3].map(i => <div className="schedule-row skeleton" key={i}/>)
          : trainers.map(t => (
              <div className="schedule-row" key={t.id}>
                <div className="schedule-trainer">
                  <div className="trainer-avatar small">{t.full_name?.slice(0, 2)}</div>
                  <span>
                    <strong>{t.full_name}</strong>
                    <small>{t.specialty || t.code}</small>
                  </span>
                </div>
                {days.map(d => (
                  <div className="schedule-cell" key={d}>
                    {by(t.id, d).length 
                      ? by(t.id, d).map(a => (
                          <button 
                            className={`schedule-pill ${a.status === 'cancelled' ? 'cancelled' : ''}`} 
                            key={a.id} 
                            onClick={() => alert(`${a.trainer?.full_name || t.full_name}\n${a.start_time} - ${a.end_time}\nสมาชิก: ${a.member?.full_name || '-'}`)}
                          >
                            <strong>{a.start_time}-{a.end_time}</strong>
                            <small>{a.member?.full_name || 'สมาชิก'}</small>
                          </button>
                        ))
                      : <span className="dash">ว่าง</span>
                    }
                  </div>
                ))}
              </div>
            ))
        }
      </section>
    </div>
  );
}

function Placeholder({ title }) {
  return (
    <div className="empty-page">
      <div className="empty-icon"><FileText size={28}/></div>
      <p className="eyebrow">MODULE READY</p>
      <h2>{title}</h2>
      <p>หน้านี้แยก URL แล้ว ส่วนระบบเทรนเนอร์ใช้งานได้ที่ <strong>/trainers</strong></p>
      <NavLink className="primary empty-link" to="/trainers">ไปหน้าจัดการเทรนเนอร์</NavLink>
    </div>
  );
}

// ==========================================
// APP ROOT
// ==========================================
function App() {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const raw = localStorage.getItem('special_booking_user');
    if (raw) {
      try {
        setUser(JSON.parse(raw));
      } catch {}
    }
  }, []);

  const auth = u => {
    localStorage.setItem('special_booking_user', JSON.stringify(u));
    setUser(u);
  };

  const logout = () => {
    localStorage.removeItem('special_booking_token');
    localStorage.removeItem('special_booking_user');
    setUser(null);
  };

  if (!user) return <Login onAuth={auth}/>;

  return (
    <Routes>
      <Route element={<Layout user={user} onLogout={logout}/>}>
        <Route path="/" element={<Home/>}/>
        <Route path="/special-areas" element={<BookingPage/>}/>
        <Route path="/special-areas/history" element={<HistoryPage/>}/>
        <Route path="/trainers" element={<TrainersPage/>}/>
        <Route path="/trainers/schedule" element={<SchedulePage/>}/>
        <Route path="/home" element={<Navigate replace to="/"/>}/>
        {nav
          .filter(n => !['/', '/special-areas', '/trainers'].includes(n.path))
          .map(n => (
            <Route key={n.path} path={n.path} element={<Placeholder title={n.label}/>}/>
          ))
        }
      </Route>
    </Routes>
  );
}

// Render Application
createRoot(document.getElementById('root')).render(
  <BrowserRouter>
    <App/>
  </BrowserRouter>
);