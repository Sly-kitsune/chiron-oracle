export default {
  async fetch(request) {
    const html = `
      <!DOCTYPE html>
      <html lang="en">
      <head>
        <meta charset="UTF-8">
        <title>Chiron Oracle ✨</title>
        <style>
          body {
            font-family: 'Segoe UI', sans-serif;
            background: linear-gradient(135deg, #fdf6e3, #ffe4e1);
            color: #333;
            padding: 2rem;
            max-width: 800px;
            margin: auto;
          }
          h1 {
            color: #b22222;
            text-align: center;
            font-size: 2.5rem;
          }
          p {
            line-height: 1.6;
            font-size: 1.2rem;
          }
          .oracle-box {
            background: #fff;
            border-radius: 12px;
            padding: 1.5rem;
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
            margin-top: 2rem;
          }
          footer {
            margin-top: 2rem;
            font-size: 0.9em;
            text-align: center;
            color: #777;
          }
        </style>
      </head>
      <body>
        <h1>✨ Chiron Oracle ✨</h1>
        <div class="oracle-box">
          <p>Welcome, seeker. Your wounds and strengths will be revealed soon...</p>
          <p>This is your oracle’s first heartbeat on the Cloudflare edge.</p>
        </div>
        <footer>Powered by Cloudflare Workers · Built by Rayana</footer>
      </body>
      </html>
    `;
    return new Response(html, {
      headers: { "content-type": "text/html" },
    });
  },
};
