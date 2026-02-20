/**
 * WebSocket real-time push — end-to-end CLI test
 * Jalankan: node scripts/manual-tests/test_websocket.mjs
 */

import { createRequire } from "module";
import { execSync } from "child_process";
import path from "path";
import { fileURLToPath } from "url";

// Resolve 'ws' dari node_modules wscat agar tidak perlu npm install terpisah
const wscatBin = execSync("which wscat").toString().trim();
const wscatRoot = path.resolve(path.dirname(wscatBin), "../lib/node_modules/wscat");
const require = createRequire(import.meta.url);
const WebSocket = require(path.join(wscatRoot, "node_modules/ws"));

const BASE = "http://localhost:8080/api/v1";
const WS_URL = "ws://localhost:8080/api/v1/ws";

// ── helpers ──────────────────────────────────────────────────────────────────
async function post(url, token, body) {
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify(body),
  });
  return res.json();
}

function sleep(ms) {
  return new Promise((r) => setTimeout(r, ms));
}

// ── main ─────────────────────────────────────────────────────────────────────
async function main() {
  console.log("═".repeat(55));
  console.log("  WebSocket Real-Time Chat — End-to-End Test");
  console.log("═".repeat(55));

  // 1. Login
  const csResp = await post(`${BASE}/cs/login`, null, {
    email: "cs@dev.local",
    password: "Password123!",
  });
  if (!csResp.success) throw new Error(`CS login failed: ${JSON.stringify(csResp)}`);
  const csToken = csResp.data.token;
  const csId = csResp.data.user_id;

  const custResp = await post(`${BASE}/users/login`, null, {
    email: "customer@dev.local",
    password: "Password123!",
  });
  if (!custResp.success) throw new Error(`Customer login failed: ${JSON.stringify(custResp)}`);
  const custToken = custResp.data.token;

  console.log(`[LOGIN] CS id      : ${csId}`);
  console.log(`[LOGIN] Customer id: ${custResp.data.user.id ?? custResp.data.user_id}`);

  // 2. CS membuka koneksi WebSocket
  const received = [];
  const ws = new WebSocket(`${WS_URL}?token=${csToken}`);

  await new Promise((resolve, reject) => {
    ws.on("open", () => {
      console.log(`[WS]    CS connected ✓`);
      resolve();
    });
    ws.on("error", reject);
    ws.on("message", (data) => {
      const raw = data.toString();
      received.push(raw);
      console.log(`[WS]    CS received frame: ${raw}`);
    });
  });

  // 3. Customer membuka / mendapatkan percakapan dengan CS
  const convResp = await post(`${BASE}/chat`, custToken, {
    participant_id: csId,
  });
  if (!convResp.success) throw new Error(`StartConversation failed: ${JSON.stringify(convResp)}`);
  const convId = convResp.data.id;
  console.log(`[REST]  conversation : ${convId}`);

  await sleep(200);

  // 4. Customer mengirim pesan
  const msgResp = await post(
    `${BASE}/chat/${convId}/messages`,
    custToken,
    { message: "Halo CS! Ada yang bisa dibantu?" }
  );
  if (!msgResp.success) throw new Error(`SendMessage failed: ${JSON.stringify(msgResp)}`);
  const msgId = msgResp.data.id;
  console.log(`[REST]  message sent : ${msgId}`);

  // 5. Tunggu push WS tiba
  await sleep(800);
  ws.close();
  await sleep(200);

  // 6. Validasi
  console.log("\n" + "─".repeat(55));
  console.log(`[WS]    Total frames diterima CS : ${received.length}`);

  if (received.length === 0) {
    console.error("❌  FAIL — CS tidak menerima frame WebSocket");
    process.exit(1);
  }

  const frame = JSON.parse(received[0]);
  console.log("Frame pertama:");
  console.log(JSON.stringify(frame, null, 2));

  if (frame.type !== "chat_message") throw new Error(`type salah: ${frame.type}`);
  if (frame.conversation_id !== convId) throw new Error(`conversation_id mismatch`);

  console.log("\n✅  PASS — CS berhasil menerima push real-time via WebSocket");
}

main().catch((err) => {
  console.error("\n❌  ERROR:", err.message);
  process.exit(1);
});
