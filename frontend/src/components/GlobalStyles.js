import { createGlobalStyle } from 'styled-components';

const GlobalStyles = createGlobalStyle`
  :root {
    --bg-primary: #0a0a0a;
    --bg-secondary: #1a1a1a;
    --bg-tertiary: #2a2a2a;
    --bg-glass: rgba(26, 26, 26, 0.8);
    --bg-glass-dark: rgba(10, 10, 10, 0.9);
    --purple-primary: #8b5cf6;
    --purple-secondary: #a855f7;
    --purple-accent: #7c3aed;
    --purple-glow: rgba(139, 92, 246, 0.3);
    --text-primary: #ffffff;
    --text-secondary: #e5e7eb;
    --text-muted: #9ca3af;
    --success: #10b981;
    --success-bg: rgba(16, 185, 129, 0.1);
    --danger: #ef4444;
    --danger-bg: rgba(239, 68, 68, 0.1);
    --neutral: #9ca3af;
    --neutral-bg: rgba(156, 163, 175, 0.1);
    --border: rgba(139, 92, 246, 0.2);
    --shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.8);
    --glow: 0 0 20px rgba(139, 92, 246, 0.4);
  }

  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  body {
    font-family: 'Inter', system-ui, -apple-system, sans-serif;
    background: radial-gradient(ellipse at top, #1e1b4b 0%, var(--bg-primary) 50%, #000000 100%);
    min-height: 100vh;
    color: var(--text-primary);
    overflow-x: hidden;
  }

  /* Animation keyframes */
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @keyframes slideUp {
    from { transform: translateY(20px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
  }

  @keyframes pulse {
    0% { box-shadow: 0 0 0 0 rgba(139, 92, 246, 0.4); }
    70% { box-shadow: 0 0 0 10px rgba(139, 92, 246, 0); }
    100% { box-shadow: 0 0 0 0 rgba(139, 92, 246, 0); }
  }
`;

export default GlobalStyles;
