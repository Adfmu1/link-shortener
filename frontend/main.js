(function () {
  'use strict';

  const API_ENDPOINT = '/api/shorten';

  const form = document.getElementById('shorten-form');
  const urlInput = document.getElementById('url-input');
  const shortenBtn = document.getElementById('shorten-btn');
  const errorMessage = document.getElementById('error-message');
  const resultBox = document.getElementById('result-box');
  const shortLinkEl = document.getElementById('short-link');
  const resetBtn = document.getElementById('reset-btn');

  function showError(message) {
    errorMessage.textContent = message;
    errorMessage.hidden = false;
  }

  function clearError() {
    errorMessage.textContent = '';
    errorMessage.hidden = true;
  }

  function setLoading(isLoading) {
    shortenBtn.disabled = isLoading;
    shortenBtn.textContent = isLoading ? 'Shortening…' : 'Shorten';
  }

  form.addEventListener('submit', async function (event) {
    event.preventDefault();
    clearError();

    const longUrl = urlInput.value.trim();
    if (!longUrl) {
      showError('Please enter a link first.');
      return;
    }

    setLoading(true);

    try {
      const response = await fetch(API_ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url: longUrl }),
      });

      if (!response.ok) {
        throw new Error('Server returned ' + response.status);
      }

      const data = await response.json();

      if (!data.Code) {
        throw new Error('Response did not include a shortUrl.');
      }

      shortLinkEl.textContent = data.Code;
      shortLinkEl.href = document.URL.replace("/main", "") + data.Code;

      form.hidden = true;
      resultBox.hidden = false;
    } catch (err) {
      showError('Could not shorten link: ' + err.message);
    } finally {
      setLoading(false);
    }
  });

  resetBtn.addEventListener('click', function () {
    urlInput.value = '';
    clearError();
    resultBox.hidden = true;
    form.hidden = false;
    urlInput.focus();
  });
})();