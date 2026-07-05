import React, { useState } from 'react';

const syncAria = (el: HTMLInputElement) => {
  el.setAttribute('aria-invalid', el.matches(':user-invalid') ? 'true' : 'false');
};

export default function App() {
  const [url, setUrl] = useState('');
  const [shortenedUrl, setShortenedUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [copied, setCopied] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setShortenedUrl('');
    setIsLoading(true);

    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

    try {
      const response = await fetch(`${apiUrl}/shorten`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url }),
      });

      if (!response.ok) {
        const errData = await response.json().catch(() => ({}));
        throw new Error(errData.error || 'Failed to shorten URL');
      }

      const data = await response.json();
      setShortenedUrl(data.short_url);
    } catch (err: any) {
      setError(err.message || 'An error occurred. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(shortenedUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy to clipboard', err);
    }
  };

  return (
    <div className="container">
      <header className="header">
        <div className="logo-container">
          <span className="logo-icon">ST</span>
        </div>
        <h1 className="title">Short Track</h1>
        <p className="subtitle">High-performance link management, simplified.</p>
      </header>

      <main className="card">
        <form onSubmit={handleSubmit} noValidate={false}>
          <div className="form-group">
            <label htmlFor="url" className="label">
              Destination URL
            </label>
            <span id="url-hint" className="hint">
              Enter a valid web link (must start with http:// or https://)
            </span>
            <div className="input-wrapper">
              <input
                type="url"
                id="url"
                className="input"
                placeholder="https://example.com/some/very/long/path"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                onBlur={(e) => syncAria(e.currentTarget)}
                onInput={(e) => syncAria(e.currentTarget)}
                required
                disabled={isLoading}
                aria-describedby="url-hint"
              />
            </div>
            <span className="error-msg">
              ❌ Please enter a valid URL (e.g. https://example.com).
            </span>
          </div>

          <button type="submit" className="btn" disabled={isLoading || !url}>
            {isLoading && <div className="spinner" />}
            {isLoading ? 'Shortening...' : 'Shorten Link'}
          </button>
        </form>

        {error && (
          <div className="error-msg" style={{ display: 'block', marginTop: '1rem' }}>
            ⚠️ {error}
          </div>
        )}

        {shortenedUrl && (
          <section className="result-section">
            <h2 className="result-label">Your Shortened Link</h2>
            <div className="result-box">
              <a
                href={shortenedUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="result-url"
              >
                {shortenedUrl}
              </a>
              <button
                onClick={handleCopy}
                className={`btn-copy ${copied ? 'copied' : ''}`}
                aria-label="Copy link to clipboard"
              >
                {copied ? '✓ Copied' : 'Copy'}
              </button>
            </div>
          </section>
        )}
      </main>
    </div>
  );
}
