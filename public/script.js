document.getElementById("oracle-form").addEventListener("submit", function(e) {
  e.preventDefault();
  const city = document.getElementById("city").value;
  const date = document.getElementById("date").value;
  const time = document.getElementById("time").value;

  const resultBox = document.getElementById("result");
  resultBox.innerHTML = `
    <p>🌌 Oracle received:</p>
    <p><strong>City:</strong> ${city}</p>
    <p><strong>Date:</strong> ${date}</p>
    <p><strong>Time:</strong> ${time}</p>
    <p>✨ Your interpretation will appear here soon...</p>
  `;
});
