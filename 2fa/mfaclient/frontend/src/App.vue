<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from "vue";
import {
  GetAccounts,
  AddAccount,
  AddAccountFromURI,
  DeleteAccount,
} from "../wailsjs/go/main/App";
import { main } from "../wailsjs/go/models";

const accounts = ref<main.AccountWithCode[]>([]);
const error = ref("");
const showAddForm = ref(false);
const addMode = ref<"manual" | "uri">("manual");
const issuer = ref("");
const label = ref("");
const secret = ref("");
const uri = ref("");
const addError = ref("");
const copiedId = ref("");

let refreshInterval: number | null = null;

async function loadAccounts() {
  try {
    accounts.value = (await GetAccounts()) || [];
    error.value = "";
  } catch (e: any) {
    error.value = e;
  }
}

async function handleAddManual() {
  addError.value = "";
  try {
    await AddAccount(
      issuer.value,
      label.value,
      secret.value.toUpperCase().replace(/\s/g, ""),
    );
    issuer.value = "";
    label.value = "";
    secret.value = "";
    showAddForm.value = false;
    await loadAccounts();
  } catch (e: any) {
    addError.value = e;
  }
}

async function handleAddFromURI() {
  addError.value = "";
  try {
    await AddAccountFromURI(uri.value);
    uri.value = "";
    showAddForm.value = false;
    await loadAccounts();
  } catch (e: any) {
    addError.value = e;
  }
}

const confirmDeleteAccount = ref<main.AccountWithCode | null>(null);

async function handleDelete(account: main.AccountWithCode) {
  confirmDeleteAccount.value = account;
}

async function confirmDelete() {
  if (!confirmDeleteAccount.value) return;
  try {
    await DeleteAccount(confirmDeleteAccount.value.id);
    await loadAccounts();
  } catch (e: any) {
    error.value = e;
  } finally {
    confirmDeleteAccount.value = null;
  }
}

async function copyCode(account: main.AccountWithCode) {
  await navigator.clipboard.writeText(account.code);
  copiedId.value = account.id;
  setTimeout(() => {
    copiedId.value = "";
  }, 1500);
}

function resetForm() {
  showAddForm.value = !showAddForm.value;
  addError.value = "";
  issuer.value = "";
  label.value = "";
  secret.value = "";
  uri.value = "";
}

onMounted(() => {
  loadAccounts();
  refreshInterval = window.setInterval(loadAccounts, 1000);
});

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval);
});
</script>

<template>
  <main>
    <header>
      <h1>🔐 2FA Authenticator</h1>
      <button class="header-btn" @click="resetForm">
        {{ showAddForm ? "✕ Cancel" : "+ Add Account" }}
      </button>
    </header>

    <div v-if="showAddForm" class="add-form">
      <div class="tabs">
        <button
          :class="{ active: addMode === 'manual' }"
          @click="addMode = 'manual'"
        >
          Manual Entry
        </button>
        <button :class="{ active: addMode === 'uri' }" @click="addMode = 'uri'">
          Paste URI
        </button>
      </div>

      <div v-if="addMode === 'manual'" class="form-fields">
        <input v-model="issuer" type="text" placeholder="Issuer" />
        <input v-model="label" type="text" placeholder="Account" />
        <input v-model="secret" type="text" placeholder="Secret Key" />
        <button class="submit-btn" @click="handleAddManual">Add Account</button>
      </div>

      <div v-else class="form-fields">
        <input v-model="uri" type="text" placeholder="otpauth://totp/..." />
        <button class="submit-btn" @click="handleAddFromURI">
          Add from URI
        </button>
      </div>

      <div v-if="addError" class="error-msg">{{ addError }}</div>
    </div>

    <div v-else class="accounts">
      <div v-for="account in accounts" :key="account.id" class="account">
        <div class="account-row">
          <div>
            <strong>{{ account.issuer || "Unknown" }}</strong>
            <span class="label">{{ account.label }}</span>
          </div>
          <button @click="handleDelete(account)" title="Delete">🗑️</button>
        </div>
        <div class="account-row">
          <button class="code" @click="copyCode(account)" title="Click to copy">
            {{ account.code }} {{ copiedId === account.id ? "✅" : "📋" }}
          </button>
          <span :class="{ warning: account.time_remaining <= 10 }"
            >⏱️ {{ account.time_remaining }}s</span
          >
        </div>
      </div>
    </div>

    <div
      v-if="confirmDeleteAccount"
      class="confirm-overlay"
      @click.self="confirmDeleteAccount = null"
    >
      <div class="confirm-dialog">
        <p>
          Delete {{ confirmDeleteAccount.issuer }} ({{
            confirmDeleteAccount.label
          }})?
        </p>
        <div class="confirm-actions">
          <button @click="confirmDeleteAccount = null">Cancel</button>
          <button class="danger-btn" @click="confirmDelete">Delete</button>
        </div>
      </div>
    </div>
  </main>
</template>
