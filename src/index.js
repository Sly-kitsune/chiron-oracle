export default {
  async fetch(request) {
    return new Response("Chiron Oracle Worker is alive 🌍", {
      headers: { "content-type": "text/plain" },
    });
  },
};
