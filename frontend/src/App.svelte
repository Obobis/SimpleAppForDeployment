<script>
  let username = "";
  let password = "";
  let isRegistering = false;
  let loggedIn = false;

  let ws;
  let wsStatus = "disconnected";

  function connectWebSocket() {
    ws = new WebSocket(`${location.origin.replace(/^http/, 'ws')}/ws`);

    ws.onopen = () => {
      wsStatus = "connected";
      ws.send(JSON.stringify({ 
        action: "auth", 
        data: { username, token: "session-token" } 
      }));
    };

    ws.onerror = (err) => {
      wsStatus = "error";
      console.error("WebSocket error:", err);
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("Received:", data);
      if (data.action === "auth_success") {
        loggedIn = true;
      }
    };
  }

  async function register() {
    const res = await fetch("/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    if (res.ok) {
      alert("Registration successful! Please login.");
      isRegistering = false;
    } else {
      alert("Registration failed!");
    }
  }

  async function login() {
    const res = await fetch("/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    if (res.ok) {
      connectWebSocket();
    } else {
      alert("Login failed!");
    }
  }

  function sendTestMessage() {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ action: "ping" }));
    }
  }
</script>

<main>
  {#if !loggedIn}
    <div class="auth-form">
      <h1>{isRegistering ? 'Register' : 'Login'}</h1>
      <input bind:value={username} placeholder="Username" />
      <input bind:value={password} type="password" placeholder="Password" />
      {#if isRegistering}
        <button on:click={register}>Register</button>
        <button type="button" on:click={() => isRegistering = false}>Login instead</button>
      {:else}
        <button on:click={login}>Login</button>
        <button type="button" on:click={() => isRegistering = true}>Register</button>
      {/if}
    </div>
  {:else}
    <div class="chat">
      <h1>Welcome, {username}!</h1>
      <p>WebSocket status: {wsStatus}</p>
      
      <button on:click={sendTestMessage} disabled={wsStatus !== "connected"}>
        Send Ping
      </button>
    </div>
  {/if}
</main>

<style>
  .auth-form, .chat {
    max-width: 400px;
    margin: 0 auto;
    padding: 20px;
  }
  input, button {
    display: block;
    width: 100%;
    margin: 10px 0;
    padding: 8px;
  }
</style>
gt